package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// DesEncrypt 使用DES CFB模式进行加密, 需要8位的key和iv
func DesEncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 创建一个DES密码分组
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对原始数据进行填充
	plainText = pkcs7PaddingDes(plainText, block.BlockSize())

	// 创建一个CFB加密模式
	cfb := cipher.NewCFBEncrypter(block, iv)

	// 加密数据
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return cipherText, nil
}

// DesDecrypt 使用DES CFB模式进行解密, 需要8位的key和iv
func DesDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// 创建一个DES密码分组
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个CFB解密模式
	cfb := cipher.NewCFBDecrypter(block, iv)

	// 解密数据
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	// 对解密后的数据进行去填充
	plainText = pkcs7UnPaddingDes(plainText)

	return plainText, nil
}

// PKCS7填充
func pkcs7PaddingDes(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

// PKCS7去填充
func pkcs7UnPaddingDes(rawData []byte) []byte {
	length := len(rawData)
	unPadding := int(rawData[length-1])
	return rawData[:(length - unPadding)]
}
