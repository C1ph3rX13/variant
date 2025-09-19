package xrand

import (
	"math/rand"
)

// Letter 生成单个随机字母 (a-z 或 A-Z)
func Letter() byte {
	if rand.Intn(2) == 0 {
		return 'a' + byte(rand.Intn(26))
	}
	return 'A' + byte(rand.Intn(26))
}

// N 生成指定长度的随机字符串
func N(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = Letter()
	}
	return string(b)
}

// Range 生成 min 到 max 范围内的随机长度字符串
func Range(min, max int) string {
	if min >= max || min < 1 {
		min, max = 1, max+1
	}
	return N(rand.Intn(max-min) + min)
}

// From 从指定字符集中生成指定长度的随机字符串
func From(pattern string, n int) string {
	if n <= 0 || pattern == "" {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = pattern[rand.Intn(len(pattern))]
	}
	return string(b)
}
