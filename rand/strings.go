package rand

import (
	"math/rand"
	"time"
	"variant/compress"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// LStrings 生成指定长度的字符串
func LStrings(length int) string {
	b := make([]rune, length)

	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}

	return string(b)
}

// RStrings 生成随机长度(1-10)的字符串
func RStrings() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, rand.Intn(16)+2)

	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}

	return string(b)
}

// LzwStrings 生成指定长度的字符串，返回 lzw 压缩后的字符串
func LzwStrings(length int) string {
	b := LStrings(length)
	lzwStrings, _ := compress.LzwCompress([]byte(b), 8)

	return lzwStrings
}

// ZstdStrings 生成指定长度的字符串，返回 zstd 压缩后的字符串
func ZstdStrings(length int) string {
	b := LStrings(length)
	lzwStrings, _ := compress.ZstdCompress(b)

	return lzwStrings
}
