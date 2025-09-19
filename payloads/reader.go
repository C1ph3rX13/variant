package payloads

import (
	"encoding/hex"
	"os"
	"regexp"
)

// hexEscapeRegex 匹配C语言风格的十六进制转义符（如 \x41）
// 匹配格式: \x 后跟两个十六进制字符（支持大小写）
var hexEscapeRegex = regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)

// CReader 读取包含C语言风格十六进制编码的文件，提取并转换所有 \xNN 转义符。
// 输入示例: "\x48\x65\x6C\x6C\x6F" 会被解析为 "Hello"
func CReader(cPath string) ([]byte, error) {
	cData, err := os.ReadFile(cPath)
	if err != nil {
		return nil, err
	}

	matches := hexEscapeRegex.FindAllString(string(cData), -1)
	var result []byte

	for _, match := range matches {
		// 解码 \x 后的十六进制部分
		decoded, err := hex.DecodeString(match[2:])
		if err != nil {
			return nil, err
		}
		result = append(result, decoded...)
	}

	return result, nil
}

// CStringsReader 读取C风格十六进制字符串并将其转换为字节序列。
// 输入示例: "\x41\x42" 会被转换为字节序列 [92 120 52 49 92 120 52 50]
// （对应字符串的ASCII码: '\' 'x' '4' '1' '\' 'x' '4' '2'）
func CStringsReader(cPath string) ([]byte, error) {
	cData, err := os.ReadFile(cPath)
	if err != nil {
		return nil, err
	}

	matches := hexEscapeRegex.FindAllString(string(cData), -1)
	var result []byte

	for _, s := range matches {
		// 将匹配字符串的每个ASCII字符转换为字节
		result = append(result, []byte(s)...)
	}

	return result, nil
}
