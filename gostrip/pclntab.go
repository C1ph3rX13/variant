package gostrip

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"encoding/binary"
	"variant/log"
)

var pclntabMagic116 = []byte{0xfa, 0xff, 0xff, 0xff, 0x00, 0x00}
var pclntabMagic118 = []byte{0xf0, 0xff, 0xff, 0xff, 0x00, 0x00}
var pclntabMagic120 = []byte{0xf1, 0xff, 0xff, 0xff, 0x00, 0x00}

type pclntabVersion int

const (
	go116 pclntabVersion = iota
	go118
	go120
)

func getPclntab(raw []byte) (uint64, uint64, binary.ByteOrder, pclntabVersion) {
	switch {
	case bytes.HasPrefix(raw, []byte(elf.ELFMAG)):
		return getPclntabFromELF(raw)
	case bytes.HasPrefix(raw, []byte{0xcf, 0xfa, 0xed, 0xfe}):
		return getPclntabFromMacho(raw)
	case bytes.HasPrefix(raw, []byte{0x4d, 0x5a}):
		return getPclntabFromPE(raw)
	default:
		log.Fatal("File format is not supported.")
	}

	return 0, 0, binary.LittleEndian, go116
}

func getPclntabFromELF(raw []byte) (uint64, uint64, binary.ByteOrder, pclntabVersion) {
	elfFile, err := elf.NewFile(bytes.NewReader(raw))
	if err != nil {
		log.Fatal("Input file is not ELF format.")
	}

	if pclntabSection := elfFile.Section(".gopclntab"); pclntabSection != nil {
		data, _ := pclntabSection.Data()
		switch {
		case bytes.HasPrefix(data, pclntabMagic116):
			return pclntabSection.Offset, pclntabSection.Size, elfFile.ByteOrder, go116
		case bytes.HasPrefix(data, pclntabMagic118):
			return pclntabSection.Offset, pclntabSection.Size, elfFile.ByteOrder, go118
		case bytes.HasPrefix(data, pclntabMagic120):
			return pclntabSection.Offset, pclntabSection.Size, elfFile.ByteOrder, go120
		default:
			log.Fatal("Failed to find pclntab.")
		}
	}

	// PIE or CGO
	dataSection := elfFile.Section(".data.rel.ro")
	if dataSection != nil {
		data, err := dataSection.Data()
		if err != nil {
			log.Fatal(err)
		}

		if tabOffset := bytes.Index(data, pclntabMagic116); tabOffset != -1 {
			return dataSection.Offset + uint64(tabOffset), dataSection.Size - uint64(tabOffset), elfFile.ByteOrder, go116
		} else if tabOffset = bytes.Index(data, pclntabMagic118); tabOffset != -1 {
			return dataSection.Offset + uint64(tabOffset), dataSection.Size - uint64(tabOffset), elfFile.ByteOrder, go118
		} else if tabOffset = bytes.Index(data, pclntabMagic120); tabOffset != -1 {
			return dataSection.Offset + uint64(tabOffset), dataSection.Size - uint64(tabOffset), elfFile.ByteOrder, go120
		}
	} else {
		log.Fatal("Failed to find pclntab.")
	}

	return 0, 0, binary.LittleEndian, go116
}

func getPclntabFromMacho(raw []byte) (uint64, uint64, binary.ByteOrder, pclntabVersion) {
	machoFile, err := macho.NewFile(bytes.NewReader(raw))
	if err != nil {
		log.Fatal("Input file is not Mach-O format.")
	}

	if pclntabSection := machoFile.Section("__gopclntab"); pclntabSection != nil {
		data, _ := pclntabSection.Data()
		switch {
		case bytes.HasPrefix(data, pclntabMagic116):
			return uint64(pclntabSection.Offset), pclntabSection.Size, machoFile.ByteOrder, go116
		case bytes.HasPrefix(data, pclntabMagic118):
			return uint64(pclntabSection.Offset), pclntabSection.Size, machoFile.ByteOrder, go118
		case bytes.HasPrefix(data, pclntabMagic120):
			return uint64(pclntabSection.Offset), pclntabSection.Size, machoFile.ByteOrder, go120
		default:
			log.Fatal("Failed to find pclntab.")
		}
	} else {
		log.Fatal("Failed to find pclntab.")
	}

	return 0, 0, binary.LittleEndian, go116
}

func getPclntabFromPE(raw []byte) (uint64, uint64, binary.ByteOrder, pclntabVersion) {
	peFile, err := pe.NewFile(bytes.NewReader(raw))
	if err != nil {
		log.Fatal("Input file is not PE format.")
	}

	dataSection := peFile.Section(".rdata")
	if dataSection != nil {
		data, err := dataSection.Data()
		if err != nil {
			log.Fatal(err)
		}

		if tabOffset := bytes.Index(data, pclntabMagic116); tabOffset != -1 {
			return uint64(dataSection.Offset) + uint64(tabOffset), uint64(dataSection.Size) - uint64(tabOffset), binary.LittleEndian, go116
		} else if tabOffset = bytes.Index(data, pclntabMagic118); tabOffset != -1 {
			return uint64(dataSection.Offset) + uint64(tabOffset), uint64(dataSection.Size) - uint64(tabOffset), binary.LittleEndian, go118
		} else if tabOffset = bytes.Index(data, pclntabMagic120); tabOffset != -1 {
			return uint64(dataSection.Offset) + uint64(tabOffset), uint64(dataSection.Size) - uint64(tabOffset), binary.LittleEndian, go120
		}
	} else {
		log.Fatal("Failed to find pclntab.")
	}

	return 0, 0, binary.LittleEndian, go116
}
