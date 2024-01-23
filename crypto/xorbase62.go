package crypto

func XorBase62Encrypt(plainText, key []byte) (string, error) {
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	base62Encode := base62.Encode(xorEncode)

	return base62Encode, err
}

func XorBase62Decrypt(cipherText string, key []byte) ([]byte, error) {

	base32Decode, err := base62.Decode(cipherText)
	if err != nil {
		return nil, err
	}

	xorDecode, err := XOREncodeDecode(base32Decode, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
