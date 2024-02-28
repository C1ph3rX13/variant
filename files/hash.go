package files

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"math"
	"os"
)

const (
	// Max file size for entropy, etc. is 2GB
	constMaxFileSize     = 2147483648
	constMaxEntropyChunk = 256000
)

type FileInfo struct {
	Path string
	Size int64
}

// openFile 打开文件并返回文件句柄和文件信息
func openFile(path string) (*os.File, FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, FileInfo{}, fmt.Errorf("couldn't open path (%s): %v", path, err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, FileInfo{}, err
	}

	if !fileInfo.Mode().IsRegular() {
		file.Close()
		return nil, FileInfo{}, fmt.Errorf("file (%s) is not a regular file", path)
	}

	return file, FileInfo{Path: path, Size: fileInfo.Size()}, nil
}

// calculateHash 计算指定哈希算法的文件哈希值
func calculateHash(path string, hashFunc func() hash.Hash) (string, error) {
	file, fileInfo, err := openFile(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if fileInfo.Size == 0 {
		return "", nil
	}

	if fileInfo.Size > int64(constMaxFileSize) {
		return "", fmt.Errorf("file size (%d) is too large (max allowed: %d)", fileInfo.Size, int64(constMaxFileSize))
	}

	hashValue := hashFunc()
	_, err = io.Copy(hashValue, file)
	if err != nil {
		return "", fmt.Errorf("couldn't read path (%s) to calculate hash: %v", path, err)
	}

	return hex.EncodeToString(hashValue.Sum(nil)), nil
}

// Entropy 计算文件的熵值
func Entropy(path string) (float64, error) {
	file, fileInfo, err := openFile(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	dataBytes := make([]byte, constMaxEntropyChunk)
	byteCounts := make([]int, 256)
	for {
		numBytesRead, err := file.Read(dataBytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		for i := 0; i < numBytesRead; i++ {
			byteCounts[int(dataBytes[i])]++
		}
	}

	var entropy float64
	for i := 0; i < 256; i++ {
		px := float64(byteCounts[i]) / float64(fileInfo.Size)
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}

	return math.Round(entropy*100) / 100, nil
}

// HashMD5 计算文件的 MD5 哈希值
func HashMD5(path string) (string, error) {
	return calculateHash(path, md5.New)
}

// HashSHA1 计算文件的 SHA1 哈希值
func HashSHA1(path string) (string, error) {
	return calculateHash(path, sha1.New)
}

// HashSHA256 计算文件的 SHA256 哈希值
func HashSHA256(path string) (string, error) {
	return calculateHash(path, sha256.New)
}

// HashSHA512 计算文件的 SHA512 哈希值
func HashSHA512(path string) (string, error) {
	return calculateHash(path, sha512.New)
}
