package wdll

import "syscall"

func GetTickCount() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, getTickCount)
}

func GetPhysicallyInstalledSystemMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(kernel32DLL, getPhysicallyInstalledSystemMemory)
}
