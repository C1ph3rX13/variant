package enc

import (
	"encoding/hex"
	"os"
	"path"
	"regexp"
)

// BinReader 读取bin类型的Payload
func BinReader(binPath string) ([]byte, error) {
	binData, err := os.ReadFile(binPath)
	if err != nil {
		return nil, err
	}

	return binData, nil
}

// CReader 读取C语言类型的Payload
func CReader(cPath string) ([]byte, error) {
	cData, err := os.ReadFile(path.Join(cPath))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\\x[0-9a-f]{2}`)
	matches := re.FindAllString(string(cData), -1)

	var buf []byte
	for _, match := range matches {
		bytes, _ := hex.DecodeString(match[2:])
		buf = append(buf, bytes...)
	}

	return buf, nil
}
