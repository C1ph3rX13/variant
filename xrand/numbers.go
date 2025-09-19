package xrand

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Number 生成 1 到 n 之间的随机整数（包含 1 和 n）
func Number(n int) int {
	return r.Intn(n) + 1
}

// Digits 生成长度为 n 的纯数字字符串（0-9）
func Digits(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = '0' + byte(r.Intn(10))
	}
	return string(b)
}

// HexDigit 生成单个十六进制字符 (0-9, a-f)
func HexDigit() byte {
	return "0123456789abcdef"[r.Intn(16)]
}

// HexString 生成长度为 n 的十六进制字符串
func HexString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = HexDigit()
	}
	return string(b)
}
