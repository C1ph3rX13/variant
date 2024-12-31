package crypto

import (
	"encoding/hex"
)

func XorSm4HexBase85Encrypt(plainText, key, iv []byte) (string, error) {

	// XOR 操作
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	// SM4 加密
	sm4Encrypt, err := Sm4CbcEncrypt(xorEncode, key, iv)
	if err != nil {
		return "", nil
	}

	// 十六进制 编码
	hexEncode := make([]byte, hex.EncodedLen(len(sm4Encrypt)))
	n := hex.Encode(hexEncode, sm4Encrypt)
	hexEncode = hexEncode[:n]

	// Base85 编码
	base85Encode, err := Base85Encode(hexEncode)
	if err != nil {
		return "", nil
	}

	return base85Encode, nil
}

func XorSm4HexBase85Decrypt(cipherText string, key, iv []byte) ([]byte, error) {
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

	// SM4 解密
	sm4Decrypt, err := Sm4CbcDecrypt(hexDecode, key, iv)
	if err != nil {
		return nil, err
	}

	// XOR 解码
	xorDecode, err := XOREncodeDecode(sm4Decrypt, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
