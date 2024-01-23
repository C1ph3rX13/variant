package crypto

import (
	"encoding/hex"
)

func XorRc4HexBase85Encrypt(plainText, key []byte) (string, error) {

	// XOR 操作
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	// RC4 加密
	rc4Encrypt, err := Rc4encrypt(xorEncode, key)
	if err != nil {
		return "", nil
	}

	// 十六进制 编码
	hexEncode := make([]byte, hex.EncodedLen(len(rc4Encrypt)))
	n := hex.Encode(hexEncode, rc4Encrypt)
	hexEncode = hexEncode[:n]

	// Base85 编码
	base85Encode, err := Base85Encode(hexEncode)
	if err != nil {
		return "", nil
	}

	return base85Encode, nil
}

func XorRc4HexBase85Decrypt(cipherText string, key []byte) ([]byte, error) {
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

	// RC4 解密
	rc4Decrypt, err := Rc4decrypt(hexDecode, key)
	if err != nil {
		return nil, err
	}

	// XOR 解码
	xorDecode, err := XOREncodeDecode(rc4Decrypt, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
