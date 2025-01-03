package gostrip

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"variant/log"
)

func GoStrip(in string, out string) {
	inFile, err := os.Open(in)
	if err != nil {
		log.Fatalf("Can't open %s", in)
	}
	defer inFile.Close()

	raw, err := io.ReadAll(inFile)
	if err != nil {
		log.Fatalf("Can't read %s", in)
	}

	offset, size, byteOrder, ver := getPclntab(raw)
	strip(raw, offset, size, byteOrder, ver)

	if out == "" {
		out = in
	}

	err = os.WriteFile(out, raw, 0775)
	if err != nil {
		log.Fatalf("Can't write %s: %s", out, err)
	}

	log.Infof("%s is stripped -> %s", in, out)
}

func strip(raw []byte, offset, size uint64, byteOrder binary.ByteOrder, ver pclntabVersion) {
	data := raw[offset : offset+size]

	ptrSize := data[7]
	uintPtr := func(b []byte) uint64 {
		if ptrSize == 4 {
			return uint64(byteOrder.Uint32(b))
		}
		return byteOrder.Uint64(b)
	}

	var funcNameOffset, fileTabOffset uint64
	if ver == go116 {
		funcNameOffset = uintPtr(data[8+2*ptrSize:])
	} else {
		funcNameOffset = uintPtr(data[8+3*ptrSize:])
	}
	funcNameTab := data[funcNameOffset:]

	if ver == go116 {
		fileTabOffset = uintPtr(data[8+4*ptrSize:])
	} else {
		fileTabOffset = uintPtr(data[8+5*ptrSize:])
	}
	fileTab := data[fileTabOffset:]

	stripNames(fileTab)
	stripNames(funcNameTab)
}

func stripNames(tab []byte) {
	reader := bytes.NewReader(tab)
	lastOffset := int64(0)

	for {
		for {
			if b, err := reader.ReadByte(); err == nil {
				if b == 0 {
					break
				}
			} else {
				log.Fatal(err)
			}
		}

		if curOffset, err := reader.Seek(0, io.SeekCurrent); err == nil {
			if curOffset-1 == lastOffset {
				break
			}

			// 替换过程
			// log.Info(string(tab[lastOffset : curOffset-1]))
			if bytes.HasPrefix(tab[lastOffset:curOffset-1], []byte("runtime.")) {
				lastOffset = curOffset
				continue
			}
			for j := lastOffset; j < curOffset-1; j++ {
				tab[j] = '?'
			}
			lastOffset = curOffset
		} else {
			log.Fatal(err)
		}
	}
}
