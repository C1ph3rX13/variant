package wdll

import "syscall"

// NewLazyDLLAndProc 返回指定 dll 的指定函数
func NewLazyDLLAndProc(dll, name string) *syscall.LazyProc {
	return syscall.NewLazyDLL(dll).NewProc(name)
}

func UuidFromStringA() *syscall.LazyProc {
	return NewLazyDLLAndProc(rpcrt4DLL, uuidFromStringA)
}

func AllocADsMem() *syscall.LazyProc {
	return NewLazyDLLAndProc(activedsDLL, allocADsMem)
}

func I_QueryTagInformation() *syscall.LazyProc {
	return NewLazyDLLAndProc(advapi32DLL, i_QueryTagInformation)
}

func EnumPageFilesW() *syscall.LazyProc {
	return NewLazyDLLAndProc(psapiDLL, enumPageFilesW)
}

func EnumerateLoadedModules() *syscall.LazyProc {
	return NewLazyDLLAndProc(dbghelpDLL, enumerateLoadedModules)
}
