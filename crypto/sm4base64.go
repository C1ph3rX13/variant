package crypto

import "encoding/base64"

// Sm4Base64Encrypt SM4加密要求Key和IV的长度为16
func Sm4Base64Encrypt(plainText, key, iv []byte) (string, error) {
	sm4Encrypt, err := Sm4CbcEncrypt(plainText, key, iv)
	if err != nil {
		return "", err
	}

	base64Encode := base64.StdEncoding.EncodeToString(sm4Encrypt)

	return base64Encode, nil
}

// Sm4Base64Decrypt SM4解密要求Key和IV的长度为16
func Sm4Base64Decrypt(cipherText string, key, iv []byte) ([]byte, error) {
	base64Decode, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	sm4Decrypt, err := Sm4CbcDecrypt(base64Decode, key, iv)
	if err != nil {
		return nil, err
	}

	return sm4Decrypt, nil
}
