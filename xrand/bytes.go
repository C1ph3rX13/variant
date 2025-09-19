package xrand

import (
	"crypto/rand"
	"fmt"
	"variant/log"
)

// Bytes 生成 n 个加密安全的随机字节
// 如果 n <= 0 或系统随机数生成器失败，返回错误
func Bytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, fmt.Errorf("invalid length: must be > 0")
	}

	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return buf, nil
}

// MustBytes 生成 n 个加密安全的随机字节
// 如果 n <= 0 或系统随机数生成器失败，触发 panic
func MustBytes(n int) []byte {
	buf, err := Bytes(n)
	if err != nil {
		log.Fatalf("Critical security failure: %v", err)
	}
	return buf
}

// MustStrToBytes 生成长度为 n 的随机字符串并转换为字节切片
// 如果 n <= 0 或系统随机数生成器失败，触发 panic
// 适用于关键操作（如密钥生成）失败应立即终止程序的场景
func MustStrToBytes(n int) []byte {
	if n <= 0 {
		log.Fatalf("invalid length: must be > 0")
	}

	buf := N(n)
	return []byte(buf)
}
