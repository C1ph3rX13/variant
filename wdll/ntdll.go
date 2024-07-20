package wdll

import "syscall"

func RtlCopyMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlCopyMemory)
}

func RtlCopyBytes() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlCopyBytes)
}

func NtQueueApcThreadEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, ntQueueApcThreadEx)
}

func EtwpCreateEtwThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwpCreateEtwThread)
}

func RtlEthernetAddressToStringA() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlEthernetAddressToStringA)
}

func RtlEthernetStringToAddressA() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlEthernetStringToAddressA)
}

func RtlIpv4StringToAddressA() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlIpv4StringToAddressA)
}

func RtlIpv4AddressToStringA() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlIpv4AddressToStringA)
}

func NtAllocateVirtualMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, ntAllocateVirtualMemory)
}

func NtWriteVirtualMemory() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, ntWriteVirtualMemory)
}

func EtwEventWrite() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwEventWrite)
}

func EtwEventWriteEx() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwEventWriteEx)
}

func EtwEventWriteFull() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwEventWriteFull)
}

func EtwEventWriteString() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwEventWriteString)
}

func EtwEventWriteTransfer() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, etwEventWriteTransfer)
}

func NtQueryInformationThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, ntQueryInformationThread)
}

func RtlCreateUserThread() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, rtlCreateUserThread)
}

func NtQueryInformationProcess() *syscall.LazyProc {
	return NewLazyDLLAndProc(ntdllDLL, ntQueryInformationProcess)
}
