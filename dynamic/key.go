package dynamic

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"variant/log"
)

func GetSelfSHA256Last16() string {
	currentFile, err := os.Executable()
	if err != nil {
		log.Fatalf("<key os.Executable()> err: %w", err)
	}

	file, err := os.Open(currentFile)
	if err != nil {
		log.Fatalf("<key os.Open()> err: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		log.Fatalf("<key sha256.New()> err: %w", err)
	}

	hashInBytes := hash.Sum(nil)
	last16Chars := fmt.Sprintf("%x", hashInBytes[24:32])

	return last16Chars
}

func GetSelfSHA256Nth(n int) (string, error) {
	currentFile, err := os.Executable()
	if err != nil {
		return "", err
	}

	file, err := os.Open(currentFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:]

	// 获取指定位数的字节
	if n > len(hashInBytes) {
		n = len(hashInBytes)
	}
	nthBytes := hashInBytes[n-1 : n]

	return fmt.Sprintf("%x", nthBytes), nil
}
