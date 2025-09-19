package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"errors"
)

var (
	ErrKeyLength = errors.New("invalid key length (must be 8/16/24 for DES/2TDEA/3TDEA)")
	ErrIVLength  = errors.New("invalid IV length (must be 8 bytes)")
	ErrBlockSize = errors.New("input not multiple of block size")
	ErrPadding   = errors.New("invalid padding")
)

// validateKey 检查密钥长度并返回实际使用的 cipher.Block
// 支持 DES(8字节)、2TDEA(16字节)、3TDEA(24字节)
func validateKey(key []byte) (cipher.Block, error) {
	switch len(key) {
	case 8:
		// DES
		return des.NewCipher(key)
	case 16:
		// 2TDEA
		return des.NewTripleDESCipher(append(key, key[:8]...))
	case 24:
		// 3TDEA
		return des.NewTripleDESCipher(key)
	default:
		return nil, ErrKeyLength
	}
}

// validateIV 检查 IV 长度是否为8字节
func validateIV(iv []byte) error {
	if len(iv) != des.BlockSize {
		return ErrIVLength
	}
	return nil
}

// pkcs7PadDES 对数据进行PKCS7填充
//
// 参数说明：
//   - data: 需要填充的原始数据
//   - blockSize: 块大小
//
// 返回值：
//   - []byte: 填充后的数据
func pkcs7PadDES(data []byte, blockSize int) []byte {
	// 计算需要填充的字节数
	padLen := blockSize - len(data)%blockSize
	// 创建填充字节切片
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	// 返回填充后的数据
	return append(data, padding...)
}

// pkcs7UnPadDES 对数据进行PKCS7去填充
//
// 参数说明：
//   - data: 需要去填充的数据
//
// 返回值：
//   - []byte: 去填充后的数据
//   - error: 错误信息
func pkcs7UnPadDES(data []byte) ([]byte, error) {
	// 检查数据是否为空
	if len(data) == 0 {
		return nil, ErrPadding
	}

	// 获取填充长度
	padLen := int(data[len(data)-1])

	// 验证填充长度的有效性
	if padLen == 0 || padLen > len(data) {
		return nil, ErrPadding
	}

	// 验证填充字节是否正确
	for _, v := range data[len(data)-padLen:] {
		if int(v) != padLen {
			return nil, ErrPadding
		}
	}

	// 返回去填充后的数据
	return data[:len(data)-padLen], nil
}

// DESCBCEncrypt 使用DES/TDEA CBC模式进行加密
//
// 参数说明：
//   - plainText: 原始明文数据
//   - key: 密钥（8/16/24字节分别对应DES/2TDEA/3TDEA）
//   - iv: 初始化向量（必须为8字节）
//
// 返回值：
//   - []byte: 加密后的密文
//   - error: 错误信息
func DESCBCEncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 验证IV长度
	if err := validateIV(iv); err != nil {
		return nil, err
	}

	// 验证密钥并创建密码分组
	block, err := validateKey(key)
	if err != nil {
		return nil, err
	}

	// 对原始数据进行PKCS7填充
	plainText = pkcs7PadDES(plainText, block.BlockSize())

	// 创建CBC加密模式
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return cipherText, nil
}

// DESCBCDecrypt 使用DES/TDEA CBC模式进行解密
//
// 参数说明：
//   - cipherText: 密文数据
//   - key: 密钥（8/16/24字节分别对应DES/2TDEA/3TDEA）
//   - iv: 初始化向量（必须为8字节）
//
// 返回值：
//   - []byte: 解密后的明文
//   - error: 错误信息
func DESCBCDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// 验证IV长度
	if err := validateIV(iv); err != nil {
		return nil, err
	}

	// 验证密钥并创建密码分组
	block, err := validateKey(key)
	if err != nil {
		return nil, err
	}

	// 验证密文长度必须是块大小的整数倍
	if len(cipherText)%block.BlockSize() != 0 {
		return nil, ErrBlockSize
	}

	// 创建CBC解密模式
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	// 对解密后的数据进行PKCS7去填充
	return pkcs7UnPadDES(plainText)
}

// DESCTREncrypt 使用DES/TDEA CTR模式进行加密
//
// 参数说明：
//   - plainText: 原始明文数据
//   - key: 密钥（8/16/24字节分别对应DES/2TDEA/3TDEA）
//   - iv: 初始化向量（必须为8字节）
//
// 返回值：
//   - []byte: 加密后的密文
//   - error: 错误信息
func DESCTREncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 验证IV长度
	if err := validateIV(iv); err != nil {
		return nil, err
	}

	// 验证密钥并创建密码分组
	block, err := validateKey(key)
	if err != nil {
		return nil, err
	}

	// 创建CTR加密模式
	stream := cipher.NewCTR(block, iv)

	// 加密数据（CTR模式无需填充）
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)

	return cipherText, nil
}

// DESCTRDecrypt 使用DES/TDEA CTR模式进行解密
//
// 参数说明：
//   - cipherText: 密文数据
//   - key: 密钥（8/16/24字节分别对应DES/2TDEA/3TDEA）
//   - iv: 初始化向量（必须为8字节）
//
// 返回值：
//   - []byte: 解密后的明文
//   - error: 错误信息
func DESCTRDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// CTR模式的加密和解密使用相同的操作
	return DESCTREncrypt(cipherText, key, iv)
}
