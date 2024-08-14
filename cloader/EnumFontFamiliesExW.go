package cloader

/*
#include <windows.h>

void enumFontFamiliesExW(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        LOGFONTW lf = { 0 };
    	lf.lfCharSet = DEFAULT_CHARSET;
    	HDC dc = GetDC(0);
    	EnumFontFamiliesExW(dc, &lf, (FONTENUMPROCW)addr, 0, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumFontFamiliesExW(p []byte) {
	C.enumFontFamiliesExW((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
