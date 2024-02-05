package dynamic

import (
	"crypto/sha256"
	"fmt"
	"variant/log"
	"variant/remote"
)

// GetIcoHex 获取图标的Hash指定切片
// 1字节(byte)等于8位(bits), 在16进制表示中，1个字节占用2个16进制位
func GetIcoHex(url string, start, end int) []byte {
	if end-start != 16 && end > 64 {
		log.Fatal("invalid range, end - start = 16 && end <= 64")
	}

	ico, err := remote.Resty(url)
	if err != nil {
		log.Fatal("get ico fail: %v", err)
	}

	hash := sha256.Sum256(ico)
	icoHash := hash[start:end]

	return []byte(fmt.Sprintf("%x", icoHash))
}
