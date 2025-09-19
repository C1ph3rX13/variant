package crypto

import (
	"encoding/hex"
	"fmt"
)

func Sm4CbcEncryptHex(plainText string, key []byte, iv []byte) ([]byte, error) {
	p, err := hex.DecodeString(plainText)
	if err != nil {
		return nil, fmt.Errorf("HEX: decodeString error: %v", err)
	}

	// 设置 Key 进行加密
	cipherText, err := Sm4CbcEncrypt(key, p, iv)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func Sm4CbcDecryptHex(cipherText string, key []byte, iv []byte) ([]byte, error) {
	c, err := hex.DecodeString(cipherText)
	if err != nil {
		return nil, fmt.Errorf("hex decode error: %v", err)
	}

	plainText, err := Sm4CbcDecrypt(key, c, iv)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
