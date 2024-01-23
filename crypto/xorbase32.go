package crypto

import (
	"encoding/base32"
)

func XorBase32Encrypt(plainText, key []byte) (string, error) {
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	base32Encode := base32.StdEncoding.EncodeToString(xorEncode)

	return base32Encode, err
}

func XorBase32Decrypt(cipherText string, key []byte) ([]byte, error) {

	base32Decode, err := base32.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	xorDecode, err := XOREncodeDecode(base32Decode, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
