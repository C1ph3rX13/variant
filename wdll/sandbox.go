package wdll

import "syscall"

func GetTickCount() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, getTickCount)
}

func GetPhysicallyInstalledSystemMemory() *syscall.LazyProc {
	return NewProcByName(DLink.Kernel32, getPhysicallyInstalledSystemMemory)
}
