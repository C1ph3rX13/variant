package rand

import (
	"math/rand"
	"variant/compress"
)

// RandomLetters 随机生成指定长度的 a-z 或 A-Z 的字符串
func RandomLetters(len int) string {
	b := make([]byte, len)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = byte(rand.Intn(26)) + 'a'
		} else {
			b[i] = byte(rand.Intn(26)) + 'A'
		}
	}
	return string(b)
}

// LStrings 生成指定长度的字符串
func LStrings(len int) string {
	length := rand.Intn(len) + 2
	return RandomLetters(length)
}

// RStrings 随机生成长度为 2-18 的字符串
func RStrings() string {
	length := rand.Intn(16) + 2
	return RandomLetters(length)
}

// LZWStrings 生成指定长度的字符串，返回 LZW 压缩后的字符串
func LZWStrings(length int) string {
	b := LStrings(length)
	lzwStrings, _ := compress.LzwCompress([]byte(b), 8)

	return lzwStrings
}

// ZSTDStrings 生成指定长度的字符串，返回 ZSTD 压缩后的字符串
func ZSTDStrings(length int) string {
	b := LStrings(length)
	lzwStrings, _ := compress.ZSTDCompress(b)

	return lzwStrings
}
