package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AESCbcEncrypt(plainText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	paddingText := pkcs7PaddingAes(plainText, blockSize)

	if len(paddingText)%blockSize != 0 {
		return nil, errors.New("input data is not a multiple of the block size")
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)

	return cipherText, nil
}

func AESCbcDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(cipherText))
	blockMode.CryptBlocks(result, cipherText)
	// 去除填充
	result = pkcs7UnPaddingAes(result)
	return result, nil
}

func pkcs7PaddingAes(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

func pkcs7UnPaddingAes(text []byte) []byte {
	unPadding := int(text[len(text)-1])
	return text[:(len(text) - unPadding)]
}
