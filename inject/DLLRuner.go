package inject

import (
	"syscall"
	"variant/log"

	"golang.org/x/sys/windows"
)

func SyscallRun(path string, procName string) {
	handle, err := syscall.LoadLibrary(path)
	if err != nil {
		log.Fatalf("fail to LoadLibrary: %v", err)
	}

	proc, err := syscall.GetProcAddress(handle, procName)
	if err != nil {
		log.Fatalf("fail to GetProcAddress: %v", err)
	}

	_, _, err = syscall.SyscallN(proc)
	if err != nil {
		log.Fatalf("fail to syscallN: %v", err)
	}
}

func WindowsRun(path string, name string) {
	handle, err := windows.LoadDLL(path)
	if err != nil {
		log.Fatalf("fail to LoadDLL: %v", err)
	}

	proc, err := handle.FindProc(name)
	if err != nil {
		log.Fatalf("fail to FindProc: %v", err)
	}

	_, _, err = proc.Call()
	if err != nil {
		log.Fatalf("fail to Call: %v", err)
	}
}
