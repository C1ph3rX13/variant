package halosgate2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Binject/debug/pe"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"unsafe"
)

type (
	DWORD     uint32
	ULONGLONG uint64
	WORD      uint16
	BYTE      uint8
	LONG      uint32
)

type unNtd struct {
	pModule uintptr
	size    uintptr
}

// Library - describes a loaded library
type Library struct {
	Name        string
	BaseAddress uintptr
	Exports     map[string]uint64
}

// sstring is the stupid internal windows definition of a unicode string.
type sstring struct {
	Length    uint16
	MaxLength uint16
	PWstr     *uint16
}

type Export struct {
	Name           string
	VirtualAddress uintptr
}

type imageExportDir struct {
	_, _                  uint32
	_, _                  uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}

const (
	MEM_COMMIT                       = 0x001000
	MEM_RESERVE                      = 0x002000
	IDX                              = 32
	IMAGE_NUMBEROF_DIRECTORY_ENTRIES = 16
)

type _IMAGE_FILE_HEADER struct {
	Machine              WORD
	NumberOfSections     WORD
	TimeDateStamp        DWORD
	PointerToSymbolTable DWORD
	NumberOfSymbols      DWORD
	SizeOfOptionalHeader WORD
	Characteristics      WORD
}

type _IMAGE_DATA_DIRECTORY struct {
	VirtualAddress DWORD
	Size           DWORD
}

type IMAGE_DATA_DIRECTORY _IMAGE_DATA_DIRECTORY

type _IMAGE_OPTIONAL_HEADER64 struct {
	Magic                       WORD
	MajorLinkerVersion          BYTE
	MinorLinkerVersion          BYTE
	SizeOfCode                  DWORD
	SizeOfInitializedData       DWORD
	SizeOfUninitializedData     DWORD
	AddressOfEntryPoint         DWORD
	BaseOfCode                  DWORD
	ImageBase                   ULONGLONG
	SectionAlignment            DWORD
	FileAlignment               DWORD
	MajorOperatingSystemVersion WORD
	MinorOperatingSystemVersion WORD
	MajorImageVersion           WORD
	MinorImageVersion           WORD
	MajorSubsystemVersion       WORD
	MinorSubsystemVersion       WORD
	Win32VersionValue           DWORD
	SizeOfImage                 DWORD
	SizeOfHeaders               DWORD
	CheckSum                    DWORD
	Subsystem                   WORD
	DllCharacteristics          WORD
	SizeOfStackReserve          ULONGLONG
	SizeOfStackCommit           ULONGLONG
	SizeOfHeapReserve           ULONGLONG
	SizeOfHeapCommit            ULONGLONG
	LoaderFlags                 DWORD
	NumberOfRvaAndSizes         DWORD
	DataDirectory               [IMAGE_NUMBEROF_DIRECTORY_ENTRIES]IMAGE_DATA_DIRECTORY
}

type _IMAGE_DOS_HEADER struct { // DOS .EXE header
	E_magic    WORD     // Magic number
	E_cblp     WORD     // Bytes on last page of file
	E_cp       WORD     // Pages in file
	E_crlc     WORD     // Relocations
	E_cparhdr  WORD     // Size of header in paragraphs
	E_minalloc WORD     // Minimum extra paragraphs needed
	E_maxalloc WORD     // Maximum extra paragraphs needed
	E_ss       WORD     // Initial (relative) SS value
	E_sp       WORD     // Initial SP value
	E_csum     WORD     // Checksum
	E_ip       WORD     // Initial IP value
	E_cs       WORD     // Initial (relative) CS value
	E_lfarlc   WORD     // File address of relocation table
	E_ovno     WORD     // Overlay number
	E_res      [4]WORD  // Reserved words
	E_oemid    WORD     // OEM identifier (for E_oeminfo)
	E_oeminfo  WORD     // OEM information; E_oemid specific
	E_res2     [10]WORD // Reserved words
	E_lfanew   LONG     // File address of new exe header
}

type IMAGE_DOS_HEADER _IMAGE_DOS_HEADER

type IMAGE_FILE_HEADER _IMAGE_FILE_HEADER
type IMAGE_OPTIONAL_HEADER64 _IMAGE_OPTIONAL_HEADER64
type IMAGE_OPTIONAL_HEADER IMAGE_OPTIONAL_HEADER64

type _IMAGE_NT_HEADERS64 struct {
	Signature      DWORD
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER
}

type IMAGE_NT_HEADERS64 _IMAGE_NT_HEADERS64
type IMAGE_NT_HEADERS IMAGE_NT_HEADERS64

func (s sstring) String() string {
	return windows.UTF16PtrToString(s.PWstr)
}

// GetModuleLoadedOrder returns the start address of module located at i in the load order.
// This might be useful if there is a function you need that isn't ntdll,
// or if some rude individual has loaded themselves before ntdll.
func gMLO(i int) (start uintptr, size uintptr, modulepath string) {
	var badstring *sstring
	start, size, badstring = getMLO(i)
	modulepath = badstring.String()
	return
}

// GetModuleLoadedOrder returns the start address of module located at i in the load order.
// This might be useful if there is a function you need that isn't in ntdll,
// or if some rude individual has loaded themselves before ntdll.
func getMLO(i int) (start uintptr, size uintptr, modulepath *sstring)

// sysIDFromRawBytes takes a byte slice and determines if there is a sysID in the expected.
func sysIDFromRawBytes(b []byte) (uint16, error) {
	return binary.LittleEndian.Uint16(b[4:8]), nil
}

// InMemLoads returns a map of loaded dll paths to current process offsets (aka images) in the current process.
// No syscalls are made.
func inMemLoads(modulename string) (uintptr, uintptr) {
	s, si, p := gMLO(0)
	start := p
	i := 1
	if strings.Contains(strings.ToLower(p), strings.ToLower(modulename)) {
		return s, si
	}
	for {
		s, si, p = gMLO(i)
		if p != "" {
			if strings.Contains(strings.ToLower(p), strings.ToLower(modulename)) {
				return s, si
			}
		}
		if p == start {
			break
		}
		i++
	}
	return 0, 0
}

func Uint16Down(b []byte, idx uint16) uint16 {
	_ = b[1] // bounds check hint to compiler
	return uint16(b[0]) - idx | uint16(b[1])<<8
}

func Uint16Up(b []byte, idx uint16) uint16 {
	_ = b[1]
	return uint16(b[0]) + idx | uint16(b[1])<<8
}

func ntH(baseAddress uintptr) *IMAGE_NT_HEADERS {
	return (*IMAGE_NT_HEADERS)(unsafe.Pointer(baseAddress + uintptr((*IMAGE_DOS_HEADER)(unsafe.Pointer(baseAddress)).E_lfanew)))
}

func Memcpy(dst, src, size uintptr) {
	for i := uintptr(0); i < size; i++ {
		*(*uint8)(unsafe.Pointer(dst + i)) = *(*uint8)(unsafe.Pointer(src + i))
	}
}

// rvaToOffset converts an RVA value from a PE file into the file offset.
func rvaToOffset(pefile *pe.File, rva uint32) uint32 {
	for _, hdr := range pefile.Sections {
		baseoffset := uint64(rva)
		if baseoffset > uint64(hdr.VirtualAddress) &&
			baseoffset < uint64(hdr.VirtualAddress+hdr.VirtualSize) {
			return rva - hdr.VirtualAddress + hdr.Offset
		}
	}
	return rva
}

func GetExport(pModuleBase uintptr) []Export {
	var exports []Export
	var pImageNtHeaders = ntH(pModuleBase)
	// IMAGE_NT_SIGNATURE
	if pImageNtHeaders.Signature != 0x00004550 {
		return nil
	}
	var pImageExportDirectory *imageExportDir

	pImageExportDirectory = ((*imageExportDir)(unsafe.Pointer(uintptr(pModuleBase + uintptr(pImageNtHeaders.OptionalHeader.DataDirectory[0].VirtualAddress)))))

	pdwAddressOfFunctions := pModuleBase + uintptr(pImageExportDirectory.AddressOfFunctions)
	pdwAddressOfNames := pModuleBase + uintptr(pImageExportDirectory.AddressOfNames)
	pwAddressOfNameOrdinales := pModuleBase + uintptr(pImageExportDirectory.AddressOfNameOrdinals)

	for cx := uintptr(0); cx < uintptr((pImageExportDirectory).NumberOfNames); cx++ {
		var export Export
		pczFunctionName := pModuleBase + uintptr(*(*uint32)(unsafe.Pointer(pdwAddressOfNames + cx*4)))
		pFunctionAddress := pModuleBase + uintptr(*(*uint32)(unsafe.Pointer(pdwAddressOfFunctions + uintptr(*(*uint16)(unsafe.Pointer(pwAddressOfNameOrdinales + cx*2)))*4)))
		export.Name = windows.BytePtrToString((*byte)(unsafe.Pointer(pczFunctionName)))
		export.VirtualAddress = uintptr(pFunctionAddress)
		exports = append(exports, export)
	}

	return exports
}

// getSysIDFromMemory takes values to resolve, and resolves from disk
func getSysIDFromMem(funcname string, hash func(string) string) (uint16, error) {
	rawstr := func(name string) string {
		return name
	}
	if hash == nil {
		hash = rawstr
	}
	// Get dll module BaseAddr
	//	// get ntdll handle
	Ntd, _ := inMemLoads(string([]byte{'n', 't', 'd', 'l', 'l'}))
	if Ntd == 0 {
		return 0, nil
	}
	ex := GetExport(Ntd)
	for _, exp := range ex {
		if strings.ToLower(hash(exp.Name)) == strings.ToLower(funcname) || strings.ToLower(hash(strings.ToLower(exp.Name))) == strings.ToLower(funcname) {
			buff := make([]byte, 10)
			if exp.VirtualAddress <= Ntd {
				return 0, nil
			}
			Memcpy(uintptr(unsafe.Pointer(&buff[0])), uintptr(exp.VirtualAddress), 10)
			// First opcodes should be :
			//	MOV R10, RCX
			// 	MOV RAX, <syscall>
			if buff[0] == 0x4c &&
				buff[1] == 0x8b &&
				buff[2] == 0xd1 &&
				buff[3] == 0xb8 &&
				buff[6] == 0x00 &&
				buff[7] == 0x00 {
				return sysIDFromRawBytes(buff)
			} else {
				for idx := uintptr(1); idx <= 500; idx++ {
					Memcpy(uintptr(unsafe.Pointer(&buff[0])), uintptr(exp.VirtualAddress+idx*IDX), 10)
					// check neighboring syscall down
					if buff[0] == 0x4c &&
						buff[1] == 0x8b &&
						buff[2] == 0xd1 &&
						buff[3] == 0xb8 &&
						buff[6] == 0x00 &&
						buff[7] == 0x00 {
						return Uint16Down(buff[4:8], uint16(idx)), nil
					}

					Memcpy(uintptr(unsafe.Pointer(&buff[0])), uintptr(exp.VirtualAddress-idx*IDX), 10)
					// check neighboring syscall up
					if buff[0] == 0x4c && //76
						buff[1] == 0x8b && //139
						buff[2] == 0xd1 && //209
						buff[3] == 0xb8 && //184
						buff[6] == 0x00 &&
						buff[7] == 0x00 {
						//buff[4] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[4])) - idx*IDX))
						//buff[5] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[5])) - idx*IDX))
						return Uint16Up(buff[4:8], uint16(idx)), nil
					}
				}
			}

			return getSysIDFromDisk(funcname, hash)
		}
	}
	return getSysIDFromDisk(funcname, hash)
}

// getSysIDFromMemory takes values to resolve, and resolves from disk.
func getSysIDFromDisk(funcname string, hash func(string) string) (uint16, error) {
	rawstr := func(name string) string {
		return name
	}

	if hash == nil {
		hash = rawstr
	}

	l := string([]byte{'c', ':', '\\', 'w', 'i', 'n', 'd', 'o', 'w', 's', '\\', 's', 'y', 's', 't', 'e', 'm', '3', '2', '\\', 'n', 't', 'd', 'l', 'l', '.', 'd', 'l', 'l'})
	p, e := pe.Open(l)
	defer p.Close()
	if e != nil {
		return 0, e
	}
	ex, e := p.Exports()
	for _, exp := range ex {
		if strings.ToLower(hash(exp.Name)) == strings.ToLower(funcname) || strings.ToLower(hash(strings.ToLower(exp.Name))) == strings.ToLower(funcname) {
			offset := rvaToOffset(p, exp.VirtualAddress)
			b, e := p.Bytes()
			if e != nil {
				return 0, e
			}
			buff := b[offset : offset+10]

			// First opcodes should be :
			//	MOV R10, RCX
			//	MOV RAX, <syscall>
			if buff[0] == 0x4c &&
				buff[1] == 0x8b &&
				buff[2] == 0xd1 &&
				buff[3] == 0xb8 &&
				buff[6] == 0x00 &&
				buff[7] == 0x00 {
				return sysIDFromRawBytes(buff)
			} else {
				for idx := uintptr(1); idx <= 500; idx++ {
					// check neighboring syscall down
					if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[0])) + idx*IDX)) == 0x4c &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[1])) + idx*IDX)) == 0x8b &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[2])) + idx*IDX)) == 0xd1 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[3])) + idx*IDX)) == 0xb8 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[6])) + idx*IDX)) == 0x00 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[7])) + idx*IDX)) == 0x00 {
						buff[4] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[4])) + idx*IDX))
						buff[5] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[5])) + idx*IDX))
						return Uint16Down(buff[4:8], uint16(idx)), nil
					}

					// check neighboring syscall up
					if *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[0])) - idx*IDX)) == 0x4c &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[1])) - idx*IDX)) == 0x8b &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[2])) - idx*IDX)) == 0xd1 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[3])) - idx*IDX)) == 0xb8 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[6])) - idx*IDX)) == 0x00 &&
						*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[7])) - idx*IDX)) == 0x00 {
						buff[4] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[4])) - idx*IDX))
						buff[5] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&buff[5])) - idx*IDX))
						return Uint16Up(buff[4:8], uint16(idx)), nil
					}
				}
			}
			return 0, errors.New("Could not find sID")
		}
	}
	return 0, errors.New("Could not find sID")
}

// DiskFuncPtr returns a pointer to the function (Virtual Address)
func DiskFuncPtr(moduleName string, funcname_hash string, hash func(string) string) (uint64, string, error) {
	rawstr := func(name string) string {
		return name
	}
	if hash == nil {
		hash = rawstr
	}

	// Get dll module BaseAddr
	phModule, _ := inMemLoads(moduleName)

	if phModule == 0 {
		syscall.LoadLibrary(moduleName)
		phModule, _ = inMemLoads(moduleName)
		if phModule == 0 {
			return 0, "", fmt.Errorf("Can't Load %s" + moduleName)
		}
	}

	// get dll exports
	pef, err := dllExports(moduleName)
	defer pef.Close()
	if err != nil {
		return 0, "", err
	}

	ex, err := pef.Exports()
	if err != nil {
		return 0, "", err
	}

	for _, exp := range ex {
		if strings.ToLower(hash(exp.Name)) == strings.ToLower(funcname_hash) || strings.ToLower(hash(strings.ToLower(exp.Name))) == strings.ToLower(funcname_hash) {
			return uint64(phModule) + uint64(exp.VirtualAddress), exp.Name, nil
		}
	}
	return 0, "", fmt.Errorf("could not find function!!! ")
}

func dllExports(dllname string) (*pe.File, error) {
	l := string([]byte{'c', ':', '\\', 'w', 'i', 'n', 'd', 'o', 'w', 's', '\\', 's', 'y', 's', 't', 'e', 'm', '3', '2', '\\'}) + dllname
	p, e := pe.Open(l)
	if e != nil {
		return nil, e
	}
	return p, nil
}

func hgSyscall(callid uint16, argh ...uintptr) (errcode uint32)

// HgSyscall calls the system fucntion specified by callid with n arguments. Work much the same as syscall.Syscall
// return value is the call error code and option error text. All args are uintptrs to make it easy.
func HgSyscall(callid uint16, argh ...uintptr) (errcode uint32, err error) {
	errcode = hgSyscall(callid, argh...)

	if errcode != 0 {
		err = fmt.Errorf("non-zero return from syscall")
	}
	return errcode, err
}

// NtdllHgate takes the exported syscall name and gets the ID it refers to. This function will access the ntdll file _on disk_, and relevant events/logs will be generated for those actions.
func DiskHgate(funcname string, hash func(string) string) (uint16, error) {
	return getSysIDFromDisk(funcname, hash)
}

// NtdllHgate takes the exported syscall name and gets the ID it refers to. This function will access the ntdll file _on disk_, and relevant events/logs will be generated for those actions.
func MemHgate(funcname string, hash func(string) string) (uint16, error) {
	return getSysIDFromMem(funcname, hash)
}
