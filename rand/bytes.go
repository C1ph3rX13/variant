package rand

import (
	"crypto/rand"
	"fmt"
)

func LBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rand bytes: %w", err)
	}

	return bytes, nil
}

func Bytes16Bit() []byte {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return bytes
}
