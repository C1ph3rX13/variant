package rand

import (
	"math/rand"
	"time"
)

func RNumbers() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(16) + 1
}

func LNumbers() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(16) + 1
}
