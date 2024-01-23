package wdll

import "syscall"

// NewProcByName 返回指定 dll 的指定函数的懒加载过程
func NewProcByName(dll *syscall.LazyDLL, name string) *syscall.LazyProc {
	return dll.NewProc(name)
}

func VirtualAlloc() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, virtualAlloc)
}

func VirtualProtect() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, virtualProtect)
}

func RtlCopyMemory() *syscall.LazyProc {
	return NewProcByName(DLink.Ntdll, rtlCopyMemory)
}

func RtlCopyBytes() *syscall.LazyProc {
	return NewProcByName(DLink.Ntdll, rtlCopyBytes)
}

func ConvertThreadToFiber() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, convertThreadToFiber)
}

func CreateFiber() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, createFiber)
}

func SwitchToFiber() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, switchToFiber)
}

func GetCurrentThread() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, getCurrentThread)
}

func NtQueueApcThreadEx() *syscall.LazyProc {
	return NewProcByName(DLink.Ntdll, ntQueueApcThreadEx)
}

func EtwpCreateEtwThread() *syscall.LazyProc {
	return NewProcByName(DLink.Ntdll, etwpCreateEtwThread)
}

func WaitForSingleObject() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, waitForSingleObject)
}

func CreateThread() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, createThread)
}

func OpenProcess() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, openProcess)
}

func VirtualAllocEx() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, virtualAllocEx)
}

func VirtualProtectEx() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, virtualProtectEx)
}

func WriteProcessMemory() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, writeProcessMemory)
}

func CreateRemoteThreadEx() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, createRemoteThreadEx)
}

func CloseHandle() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, closeHandle)
}

func HeapCreate() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, heapCreate)
}

func HeapAlloc() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, heapAlloc)
}

func EnumSystemLocalesA() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, enumSystemLocalesA)
}

func UuidFromStringA() *syscall.LazyProc {
	return NewProcByName(DLink.Rpcrt4, uuidFromStringA)
}

func GetCurrentProcess() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, getCurrentProcess)
}
