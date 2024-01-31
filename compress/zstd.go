package compress

import (
	"encoding/hex"
	"fmt"
	"github.com/klauspost/compress/zstd"
)

// ZstdCompress 返回压缩后的十六进制字符串
func ZstdCompress(s string) (string, error) {
	// 创建一个新的 Zstandard 压缩器
	zw, err := zstd.NewWriter(nil)
	if err != nil {
		return "", fmt.Errorf("<zstd.NewWriter()> err: %v", err)
	}
	defer zw.Close()

	// 将字符串写入压缩器并返回压缩数据
	compressed := make([]byte, 0, len(s))
	compressed = zw.EncodeAll([]byte(s), compressed)
	// 十六进制编码
	hexEncode := hex.EncodeToString(compressed)
	return hexEncode, nil
}

func ZstdDecompress(compressed string) (string, error) {
	// 创建一个新的 Zstandard 解压缩器
	zr, err := zstd.NewReader(nil)
	if err != nil {
		return "", fmt.Errorf("<zstd.NewReader()> err: %v", err)
	}
	defer zr.Close()

	// 十六进制解码
	hexDecode, _ := hex.DecodeString(compressed)
	// 对压缩数据进行解压缩并返回原始字符串
	decompressed, err := zr.DecodeAll(hexDecode, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decompress data: %v", err)
	}
	return string(decompressed), nil
}
