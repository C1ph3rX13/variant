package crypto

import (
	"encoding/base32"
)

func AesBase32Encrypt(plainText, key, iv []byte) (string, error) {
	// AES 加密
	aesEncrypt, err := AESCBCEncrypt(plainText, key, iv)
	if err != nil {
		return "", nil
	}

	// Base32 编码
	base32Encode := base32.StdEncoding.EncodeToString(aesEncrypt)
	if err != nil {
		return "", nil
	}

	return base32Encode, nil
}

func AesBase32Decrypt(cipherText string, key, iv []byte) ([]byte, error) {
	// Base32 解码
	base32Decode, err := base32.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	// AES 解密
	aesDecrypt, err := AESCBCDecrypt(base32Decode, key, iv)
	if err != nil {
		return nil, err
	}

	return aesDecrypt, nil
}
