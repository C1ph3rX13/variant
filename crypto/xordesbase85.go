package crypto

// XorDesBase85Encrypt Des加密要求Key和IV的长度为8
func XorDesBase85Encrypt(plainText, key, iv []byte) (string, error) {
	// XOR 编码
	xorEncode, err := XOREncodeDecode(plainText, key)
	if err != nil {
		return "", nil
	}

	// Des 加密
	desEncrypt, err := DesCFBEncrypt(xorEncode, key, iv)
	if err != nil {
		return "", nil
	}

	// Base85 编码
	base85Encode, err := Base85Encode(desEncrypt)
	if err != nil {
		return "", nil
	}

	return base85Encode, nil
}

// XorDesBase85Decrypt Des加密要求Key和IV的长度为8
func XorDesBase85Decrypt(cipherText string, key, iv []byte) ([]byte, error) {
	// Base85 解码
	base85Decode, err := Base85Decode(cipherText)
	if err != nil {
		return nil, err
	}

	// Des 解密
	desDecrypt, err := DesCFBDecrypt(base85Decode, key, iv)
	if err != nil {
		return nil, err
	}

	// XOR 解码
	xorDecode, err := XOREncodeDecode(desDecrypt, key)
	if err != nil {
		return nil, err
	}

	return xorDecode, nil
}
