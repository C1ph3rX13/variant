package compress

import (
	"bytes"
	"compress/lzw"
	"encoding/hex"
	"fmt"
)

// LzwCompress 压缩函数
func LzwCompress(data []byte, ratio int) (string, error) {
	defaultRatio := 8 // 默认为 8
	if ratio > 0 && ratio <= 16 {
		defaultRatio = ratio
	}

	var buf bytes.Buffer
	writer := lzw.NewWriter(&buf, lzw.LSB, defaultRatio)
	_, err := writer.Write(data)
	if err != nil {
		return "", fmt.Errorf("压缩数据时发生错误：%v", err)
	}
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("关闭压缩写入器时发生错误：%v", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// LzwDecompress 解压函数
func LzwDecompress(data []byte, ratio int) (string, error) {
	defaultRatio := 8 // 默认为 8
	if ratio > 0 && ratio <= 16 {
		defaultRatio = ratio
	}

	decoded, err := hex.DecodeString(string(data))
	if err != nil {
		return "", fmt.Errorf("解码数据时发生错误：%v", err)
	}

	reader := lzw.NewReader(bytes.NewReader(decoded), lzw.LSB, defaultRatio)
	defer reader.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return "", fmt.Errorf("解压数据时发生错误：%v", err)
	}

	return string(buf.Bytes()), nil
}
