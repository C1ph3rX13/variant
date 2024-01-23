package crypto

func XorRc4Base85Encrypt(plainText, key []byte) (string, error) {

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

	// Base85 编码
	base85Encode, err := Base85Encode(rc4Encrypt)
	if err != nil {
		return "", nil
	}

	return base85Encode, nil
}

func XorRc4Base85Decrypt(cipherText string, key []byte) ([]byte, error) {
	// Base85 解码
	base85Decode, err := Base85Decode(cipherText)
	if err != nil {
		return nil, err
	}

	// RC4 解密
	rc4Decrypt, err := Rc4decrypt(base85Decode, key)
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
