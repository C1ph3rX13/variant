package crypto

import "crypto/rc4"

func Rc4encrypt(plainText, key []byte) ([]byte, error) {
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(plainText))
	c.XORKeyStream(ciphertext, plainText)
	return ciphertext, err
}

func Rc4decrypt(ciphertext, key []byte) ([]byte, error) {
	c, err := rc4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(ciphertext))
	c.XORKeyStream(plainText, ciphertext)
	return plainText, err
}
