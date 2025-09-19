package inject

import "unsafe"

// IMAGE_DATA_DIRECTORY represents a data directory entry.
type IMAGE_DATA_DIRECTORY struct {
	VirtualAddress uint32
	Size           uint32
}

type IMAGE_SECTION_HEADER struct {
	Name                 [8]byte
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      uint32
}

type BASE_RELOCATION_BLOCK struct {
	PageAddress uint32
	BlockSize   uint32
}

// BASE_RELOCATION_ENTRY represents the base relocation entry structure
type BASE_RELOCATION_ENTRY struct {
	OffsetType uint16 // Combined field for Offset and Type
}

// Offset extracts the Offset from the combined field
func (bre BASE_RELOCATION_ENTRY) Offset() uint16 {
	return bre.OffsetType & 0xFFF
}

// Type extracts the Type from the combined field
func (bre BASE_RELOCATION_ENTRY) Type() uint16 {
	return (bre.OffsetType >> 12) & 0xF
}

type IMAGE_IMPORT_DESCRIPTOR struct {
	Characteristics uint32
	TimeDateStamp   uint32
	ForwarderChain  uint32
	Name            uint32
	FirstThunk      uint32
}

func uintptrToBytes(ptr uintptr) []byte {
	// Create a pointer to the uintptr value
	ptrPtr := unsafe.Pointer(&ptr)

	// Convert the pointer to a byte slice
	byteSlice := make([]byte, unsafe.Sizeof(ptr))
	for i := 0; i < int(unsafe.Sizeof(ptr)); i++ {
		byteSlice[i] = *(*byte)(unsafe.Pointer(uintptr(ptrPtr) + uintptr(i)))
	}

	return byteSlice
}

const (
	IMAGE_DIRECTORY_ENTRY_IMPORT    = 0x1
	IMAGE_DIRECTORY_ENTRY_BASERELOC = 0x5
	DLL_PROCESS_ATTACH              = 0x1
)
