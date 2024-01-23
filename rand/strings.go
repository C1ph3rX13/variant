package rand

import (
	"math/rand"
	"time"
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
	b := make([]rune, rand.Intn(10)+1)

	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}

	return string(b)
}
