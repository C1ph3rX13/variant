package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// DesCFBEncrypt 使用DES CFB模式进行加密, 需要8位的key和iv
func DesCFBEncrypt(plainText, key, iv []byte) ([]byte, error) {
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

// DesCFBDecrypt 使用DES CFB模式进行解密, 需要8位的key和iv
func DesCFBDecrypt(cipherText, key, iv []byte) ([]byte, error) {
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

func DesCBCEncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 创建一个DES密码分组
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 对原始数据进行填充
	plainText = pkcs7PaddingDes(plainText, block.BlockSize())

	// 创建一个CFB加密模式
	cfb := cipher.NewCBCDecrypter(block, iv)

	// 加密数据
	cipherText := make([]byte, len(plainText))
	cfb.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

func DesCBCDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// 创建一个DES密码分组
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建一个CFB解密模式
	cfb := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plainText := make([]byte, len(cipherText))
	cfb.CryptBlocks(plainText, cipherText)

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
func pkcs7UnPaddingDes(plainText []byte) []byte {
	length := len(plainText)
	unPadding := int(plainText[length-1])
	return plainText[:(length - unPadding)]
}
