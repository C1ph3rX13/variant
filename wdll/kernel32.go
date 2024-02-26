package wdll

import (
	"syscall"
)

func VirtualAlloc() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, virtualAlloc)
}

func VirtualProtect() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, virtualProtect)
}

func RtlMoveMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, rtlMoveMemory)
}

func ConvertThreadToFiber() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, convertThreadToFiber)
}

func CreateFiber() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, createFiber)
}

func SwitchToFiber() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, switchToFiber)
}

func GetCurrentThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, getCurrentThread)
}

func WaitForSingleObject() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, waitForSingleObject)
}

func CreateThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, createThread)
}

func OpenProcess() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, openProcess)
}

func VirtualAllocEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, virtualAllocEx)
}

func VirtualProtectEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, virtualProtectEx)
}

func WriteProcessMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, writeProcessMemory)
}

func CreateRemoteThreadEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, createRemoteThreadEx)
}

func CloseHandle() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, closeHandle)
}

func HeapCreate() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, heapCreate)
}

func HeapAlloc() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, heapAlloc)
}

func EnumSystemLocalesA() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, enumSystemLocalesA)
}

func GetCurrentProcess() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, getCurrentProcess)
}

func QueueUserAPC() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, queueUserAPC)
}

func EnumSystemLocalesW() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, enumSystemLocalesW)
}

func EnumSystemLocalesEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, enumSystemLocalesEx)
}

func OpenThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, openThread)
}

func TerminateThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, terminateThread)
}

func ReadProcessMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, readProcessMemory)
}

func CreateToolhelp32Snapshot() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, createToolhelp32Snapshot)
}

func Thread32First() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, thread32First)
}
