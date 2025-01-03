package gostrip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"variant/log"
	"variant/rand"

	"github.com/Binject/debug/pe"
)

func PEBoom(buff []byte, size int) []byte {
	boomBuff := make([]byte, size*1024*1024) // 设置膨胀体积

	peBuff, err := pe.NewFile(bytes.NewReader(buff)) // 读取PE文件
	if err != nil {
		log.Fatal("Error Reading Inputted File")
	}

	var sectionAlignment, fileAlignment, scAddr uint32
	var imageBase uint64
	var shellcode []byte
	lastSection := peBuff.Sections[peBuff.NumberOfSections-1]
	switch file := (peBuff.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		imageBase = uint64(file.ImageBase)
		sectionAlignment = file.SectionAlignment
		fileAlignment = file.FileAlignment
		scAddr = align(lastSection.Size, fileAlignment, lastSection.Offset)
	case *pe.OptionalHeader64:
		imageBase = file.ImageBase
		sectionAlignment = file.SectionAlignment
		fileAlignment = file.FileAlignment
		scAddr = align(lastSection.Size, fileAlignment, lastSection.Offset)
	}

	buf := bytes.NewBuffer(boomBuff)
	w := bufio.NewWriter(buf)
	WriteErr := binary.Write(w, binary.LittleEndian, imageBase)
	if err != nil {
		log.Fatalf("Error Writing: %v", WriteErr)
	}
	flushErr := w.Flush()
	if err != nil {
		log.Fatalf("Error Flushing: %v", flushErr)
	}
	shellcode = buf.Bytes()

	shellcodeLen := len(shellcode)
	newSection := new(pe.Section)
	newSection.Name = "." + rand.RandomLetters(5)
	o := []byte(newSection.Name)

	newSection.OriginalName = [8]byte{o[0], o[1], o[2], o[3], o[4], o[5], 0, 0}
	newSection.VirtualSize = uint32(shellcodeLen)
	newSection.VirtualAddress = align(lastSection.VirtualSize, sectionAlignment, lastSection.VirtualAddress)
	newSection.Size = align(uint32(shellcodeLen), fileAlignment, 0)
	newSection.Offset = align(lastSection.Size, fileAlignment, lastSection.Offset)
	newSection.Characteristics = pe.IMAGE_SCN_CNT_CODE | pe.IMAGE_SCN_MEM_EXECUTE | pe.IMAGE_SCN_MEM_READ
	peBuff.InsertionAddr = scAddr
	peBuff.InsertionBytes = shellcode

	switch file := (peBuff.OptionalHeader).(type) {
	case *pe.OptionalHeader32:
		v := newSection.VirtualSize
		if v == 0 {
			v = newSection.Size
		}
		file.SizeOfImage = align(v, sectionAlignment, newSection.VirtualAddress)
		file.CheckSum = 0
	case *pe.OptionalHeader64:
		v := newSection.VirtualSize
		if v == 0 {
			v = newSection.Size
		}
		file.SizeOfImage = align(v, sectionAlignment, newSection.VirtualAddress)
		file.CheckSum = 0
	}
	peBuff.FileHeader.NumberOfSections++
	peBuff.Sections = append(peBuff.Sections, newSection)
	Bytes, _ := peBuff.Bytes()

	log.Info("Boom Successfully")
	return Bytes
}

func align(size, align, addr uint32) uint32 {
	if size%align == 0 {
		return addr + size
	}
	return addr + (size/align+1)*align
}
