package rand

import (
	"crypto/rand"
	"fmt"
)

func LBytes(length int) ([]byte, error) {
	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rand bytes: %w", err)
	}

	return b, nil
}

func LByteStrings(length int) []byte {
	bString := RandomLetters(length)
	return []byte(bString)
}
