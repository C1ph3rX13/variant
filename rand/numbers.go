package rand

import (
	"math/rand"
	"time"
)

// RNumbers 生成 1-n 之间的随机整数
func RNumbers(n int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(n) + 1
}

// RNumbersString 生成长度为 n 的数字字符串
func RNumbersString(n int) string {
	b := make([]byte, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = byte(r.Intn(10)) + '0'
	}
	return string(b)
}
