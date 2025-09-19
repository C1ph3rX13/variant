package inject

import (
	"debug/pe"
)

/* AddressOfEntryPointInject Type */

// IMAGE_NT_HEADERS64 type
type IMAGE_NT_HEADERS64 struct {
	Signature      uint32
	FileHeader     pe.FileHeader
	OptionalHeader pe.OptionalHeader64
}

type ImageDosHeader struct {
	E_magic    uint16
	E_cblp     uint16
	E_cp       uint16
	E_crlc     uint16
	E_cparhdr  uint16
	Eminalloc  uint16
	E_maxalloc uint16
	E_ss       uint16
	E_sp       uint16
	E_csum     uint16
	Eip        uint16
	E_cs       uint16
	E_lfarlc   uint16
	E_ovno     uint16
	E_res      []uint16
	E_oemid    uint16
	E_oeminfo  uint16
	E_res2     []uint16
	E_lfanew   uint32
}

type PImageDosHeader *ImageDosHeader

/* AddressOfEntryPointInject Type */

/* NtCreateSection Type */

type clientID struct {
	UniqueProcess uintptr
	UniqueThread  uintptr
}

type objectAttrs struct {
	Length                   uintptr
	RootDirectory            uintptr
	ObjectName               uintptr
	Attributes               uintptr
	SecurityDescriptor       uintptr
	SecurityQualityOfService uintptr
}

// SecCommit is the SEC_COMMIT const from winnt.h
const SecCommit = 0x08000000

// SectionWrite is the SECTION_MAP_WRITE const from winnt.h
const SectionWrite = 0x2

// SectionRead is the SECTION_MAP_READ const from winnt.h
const SectionRead = 0x4

// SectionExecute is the SECTION_MAP_EXECUTE const from winnt.h
const SectionExecute = 0x8

// SectionRWX is the combination of READ, WRITE, and EXECUTE
const SectionRWX = SectionWrite | SectionRead | SectionExecute

/* NtCreateSection Type */

type IMAGE_REL_BASED uint16

const (
	IMAGE_REL_BASED_ABSOLUTE       IMAGE_REL_BASED = 0  //The base relocation is skipped. This type can be used to pad a block.
	IMAGE_REL_BASED_HIGH           IMAGE_REL_BASED = 1  //The base relocation adds the high 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the high value of a 32-bit word.
	IMAGE_REL_BASED_LOW            IMAGE_REL_BASED = 2  //The base relocation adds the low 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the low half of a 32-bit word.
	IMAGE_REL_BASED_HIGHLOW        IMAGE_REL_BASED = 3  //The base relocation applies all 32 bits of the difference to the 32-bit field at offset.
	IMAGE_REL_BASED_HIGHADJ        IMAGE_REL_BASED = 4  //The base relocation adds the high 16 bits of the difference to the 16-bit field at offset. The 16-bit field represents the high value of a 32-bit word. The low 16 bits of the 32-bit value are stored in the 16-bit word that follows this base relocation. This means that this base relocation occupies two slots.
	IMAGE_REL_BASED_MIPS_JMPADDR   IMAGE_REL_BASED = 5  //The relocation interpretation is dependent on the machine type.When the machine type is MIPS, the base relocation applies to a MIPS jump instruction.
	IMAGE_REL_BASED_ARM_MOV32      IMAGE_REL_BASED = 5  //This relocation is meaningful only when the machine type is ARM or Thumb. The base relocation applies the 32-bit address of a symbol across a consecutive MOVW/MOVT instruction pair.
	IMAGE_REL_BASED_RISCV_HIGH20   IMAGE_REL_BASED = 5  //This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the high 20 bits of a 32-bit absolute address.
	IMAGE_REL_BASED_THUMB_MOV32    IMAGE_REL_BASED = 7  //This relocation is meaningful only when the machine type is Thumb. The base relocation applies the 32-bit address of a symbol to a consecutive MOVW/MOVT instruction pair.
	IMAGE_REL_BASED_RISCV_LOW12I   IMAGE_REL_BASED = 7  //This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the low 12 bits of a 32-bit absolute address formed in RISC-V I-type instruction format.
	IMAGE_REL_BASED_RISCV_LOW12S   IMAGE_REL_BASED = 8  //This relocation is only meaningful when the machine type is RISC-V. The base relocation applies to the low 12 bits of a 32-bit absolute address formed in RISC-V S-type instruction format.
	IMAGE_REL_BASED_MIPS_JMPADDR16 IMAGE_REL_BASED = 9  //The relocation is only meaningful when the machine type is MIPS. The base relocation applies to a MIPS16 jump instruction.
	IMAGE_REL_BASED_DIR64          IMAGE_REL_BASED = 10 //The base relocation applies the difference to the 64-bit field at offset.
)
