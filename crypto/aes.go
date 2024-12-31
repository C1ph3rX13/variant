package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// 检查密钥长度是否为16、24或32字节
func checkKeyLength(key []byte) error {
	switch len(key) {
	case 16, 24, 32:
		return nil
	default:
		return errors.New("key must be 16, 24, or 32 bytes long")
	}
}

// 检查IV长度是否等于AES块大小（16字节）
func checkIVLength(iv []byte) error {
	if len(iv) != aes.BlockSize {
		return fmt.Errorf("IV must be %d bytes long", aes.BlockSize)
	}
	return nil
}

// AESCBCEncrypt AES CBC模式加密
func AESCBCEncrypt(plainText, key, iv []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkIVLength(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	paddingText := pkcs7PaddingAes(plainText, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)

	return cipherText, nil
}

// AESCBCDecrypt AES CBC模式解密
func AESCBCDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkIVLength(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(cipherText))
	blockMode.CryptBlocks(result, cipherText)
	result = pkcs7UnPaddingAes(result)
	return result, nil
}

// PKCS7填充
func pkcs7PaddingAes(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

// 去除PKCS7填充
func pkcs7UnPaddingAes(text []byte) []byte {
	if len(text) == 0 {
		return text
	}
	unPadding := int(text[len(text)-1])
	if unPadding > len(text) {
		return text // 或者返回错误，根据需求决定
	}
	return text[:(len(text) - unPadding)]
}

// AESCFBEncrypt AES CFB模式加密
func AESCFBEncrypt(plainText, key, iv []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkIVLength(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

// AESCFBDecrypt AES CFB模式解密
func AESCFBDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkIVLength(iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCFBDecrypter(block, iv)
	result := make([]byte, len(cipherText))
	blockMode.XORKeyStream(result, cipherText)
	return result, nil
}
