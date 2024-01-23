package crypto

import "errors"

func XOREncodeDecode(plainText, key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("key cannot be empty")
	}

	cipherText := make([]byte, len(plainText))
	keyLen := len(key)

	for i := 0; i < len(plainText); i++ {
		cipherText[i] = plainText[i] ^ key[i%keyLen]
	}

	return cipherText, nil
}
