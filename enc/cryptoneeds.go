package enc

import "variant/crypto"

type Sign struct {
	Sm4Base64Encrypt       func(plainText, key, iv []byte) (string, error)
	Sm4CbcEncryptHex       func(rawData, keyHex, ivHex string) ([]byte, error)
	XorBase32Encrypt       func([]byte, []byte) (string, error)
	XorAesHexBase85Encrypt func(plainText, key, iv []byte) (string, error)
}

func NewEnc() *Sign {
	return &Sign{
		Sm4Base64Encrypt:       crypto.Sm4Base64Encrypt,
		Sm4CbcEncryptHex:       crypto.Sm4CbcEncryptHex,
		XorBase32Encrypt:       crypto.XorBase32Encrypt,
		XorAesHexBase85Encrypt: crypto.XorAesHexBase85Encrypt,
	}
}
