package crypto

import (
	"encoding/hex"
)

func XorAesHexBase85Encrypt(plainText, key, iv []byte) (string, error) {

	// XOR 操作
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	// AES 加密
	aesEncrypt, err := AESCBCEncrypt(xorEncode, key, iv)
	if err != nil {
		return "", nil
	}

	// 十六进制 编码
	hexEncode := make([]byte, hex.EncodedLen(len(aesEncrypt)))
	n := hex.Encode(hexEncode, aesEncrypt)
	hexEncode = hexEncode[:n]

	// Base85 编码
	base85Encode, err := Base85Encode(hexEncode)
	if err != nil {
		return "", nil
	}

	return base85Encode, nil
}

func XorAesHexBase85Decrypt(cipherText string, key, iv []byte) ([]byte, error) {
	// Base85 解码
	base85Decode, err := Base85Decode(cipherText)
	if err != nil {
		return nil, err
	}

	// 十六进制 解码
	hexDecode := make([]byte, hex.DecodedLen(len(base85Decode)))
	n, err := hex.Decode(hexDecode, base85Decode)
	if err != nil {
		return nil, err
	}
	hexDecode = hexDecode[:n]

	// AES 解密
	aesDecrypt, err := AESCBCDecrypt(hexDecode, key, iv)
	if err != nil {
		return nil, err
	}

	// XOR 解码
	xorDecode, err := XOREncodeDecode(aesDecrypt, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
