package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// 预定义错误
var (
	// ErrInvalidKeyLength 表示密钥长度不符合要求
	ErrInvalidKeyLength = errors.New("KEY must be 16, 24, or 32 bytes long")
	// ErrInvalidIVLength 表示IV长度不符合要求
	ErrInvalidIVLength = fmt.Errorf("IV must be %d bytes long", aes.BlockSize)
	// ErrInvalidNonceLength 表示nonce长度不符合要求
	ErrInvalidNonceLength = errors.New("nonce must be 12 bytes long for GCM")
)

// checkKeyLength 验证密钥长度是否为16、24或32字节。
// 返回值：
//   - nil: 密钥长度合法
//   - error: 密钥长度不合法
func checkKeyLength(key []byte) error {
	switch len(key) {
	case 16, 24, 32:
		return nil
	default:
		return ErrInvalidKeyLength
	}
}

// checkIVLength 验证IV长度是否等于AES块大小（16字节）。
// 返回值：
//   - nil: IV长度合法
//   - error: IV长度不合法
func checkIVLength(iv []byte) error {
	if len(iv) != aes.BlockSize {
		return ErrInvalidIVLength
	}
	return nil
}

// checkNonceLength 验证nonce长度是否为12字节（推荐用于GCM）
// 返回值：
//   - nil: nonce长度合法
//   - error: nonce长度不合法
func checkNonceLength(nonce []byte) error {
	if len(nonce) != 12 {
		return ErrInvalidNonceLength
	}
	return nil
}

// pkcs7PadAES 对数据进行PKCS7填充。
// 参数：
//
//	text     - 需要填充的数据
//	blockSize - 块大小（如AES为16字节）
//
// 返回值：
//   - 填充后的数据
func pkcs7PadAES(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

// pkcs7UnPadAES 去除PKCS7填充。
// 参数：
//
//	text - 填充后的数据
//
// 返回值：
//   - 去除填充后的数据
//   - error: 填充无效
func pkcs7UnPadAES(text []byte) ([]byte, error) {
	if len(text) == 0 {
		return nil, errors.New("empty input")
	}
	unPadding := int(text[len(text)-1])
	if unPadding > len(text) || unPadding <= 0 {
		return nil, errors.New("invalid PKCS7 padding")
	}
	return text[:len(text)-unPadding], nil
}

// AESCBCEncrypt 使用AES CBC模式对明文进行加密。
// 参数：
//
//	plainText - 明文数据
//	key       - 加密密钥（必须为16、24或32字节）
//	iv        - 初始化向量（必须为16字节）
//
// 返回值：
//   - 加密后的密文
//   - error: 错误信息
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
	paddedText := pkcs7PadAES(plainText, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(paddedText))
	blockMode.CryptBlocks(cipherText, paddedText)

	return cipherText, nil
}

// AESCBCDecrypt 使用AES CBC模式对密文进行解密。
// 参数：
//
//	cipherText - 密文数据
//	key        - 解密密钥（必须为16、24或32字节）
//	iv         - 初始化向量（必须为16字节）
//
// 返回值：
//   - 解密后的明文
//   - error: 错误信息
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

	// 去除PKCS7填充
	unpadded, err := pkcs7UnPadAES(result)
	if err != nil {
		return nil, err
	}
	return unpadded, nil
}

// AESCTREncrypt 使用AES CTR模式对明文进行加密。
// CTR 模式无需填充，加密和解密使用相同的流生成器。
// 参数：
//
//	plainText - 明文数据
//	key       - 加密密钥（必须为16、24或32字节）
//	iv        - 初始化向量（必须为16字节）
//
// 返回值：
//   - 加密后的密文
//   - error: 错误信息
func AESCTREncrypt(plainText, key, iv []byte) ([]byte, error) {
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

	blockMode := cipher.NewCTR(block, iv)
	cipherText := make([]byte, len(plainText))
	blockMode.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

// AESCTRDecrypt 使用AES CTR模式对密文进行解密。
// CTR 模式无需填充，加密和解密使用相同的流生成器。
// 参数：
//
//	cipherText - 密文数据
//	key        - 解密密钥（必须为16、24或32字节）
//	iv         - 初始化向量（必须为16字节）
//
// 返回值：
//   - 解密后的明文
//   - error: 错误信息
func AESCTRDecrypt(cipherText, key, iv []byte) ([]byte, error) {
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

	blockMode := cipher.NewCTR(block, iv)
	result := make([]byte, len(cipherText))
	blockMode.XORKeyStream(result, cipherText)
	return result, nil
}

// AESGCMEncrypt 使用AES GCM模式对明文进行加密。
// GCM 模式提供认证加密（AEAD），加密结果包含认证标签。
// 参数：
//
//	plainText   - 明文数据
//	key         - 加密密钥（必须为16、24或32字节）
//	nonce       - 随机数（通常为12字节）
//	additional  - 附加数据（可选，用于认证，不影响加密内容）
//
// 返回值：
//   - 加密后的密文（包含认证标签）
//   - error: 错误信息
func AESGCMEncrypt(plainText, key, nonce []byte, additional []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkNonceLength(nonce); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText := aesgcm.Seal(nil, nonce, plainText, additional)
	return cipherText, nil
}

// AESGCMDecrypt 使用AES GCM模式对密文进行解密。
// GCM 模式会验证认证标签，若验证失败则返回错误。
// 参数：
//
//	cipherText  - 密文数据（包含认证标签）
//	key         - 解密密钥（必须为16、24或32字节）
//	nonce       - 随机数（通常为12字节）
//	additional  - 附加数据（必须与加密时一致，用于认证）
//
// 返回值：
//   - 解密后的明文
//   - error: 错误信息
func AESGCMDecrypt(cipherText, key, nonce []byte, additional []byte) ([]byte, error) {
	if err := checkKeyLength(key); err != nil {
		return nil, err
	}
	if err := checkNonceLength(nonce); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesgcm.Open(nil, nonce, cipherText, additional)
	if err != nil {
		return nil, errors.New("GCM decryption failed: invalid tag or corrupted data")
	}
	return plainText, nil
}
