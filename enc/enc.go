package enc

import (
	"errors"
	"fmt"
	"variant/crypto"
)

func (params Payload) EncSetKey() (string, error) {
	binRaw, err := BinReader(params.PlainText)
	if err != nil {
		return "", fmt.Errorf("<enc.BinReader()> err: %w", err)
	}

	cipherText, err := crypto.XorRc4Base85Encrypt(binRaw, params.Key)
	if err != nil {
		return "", fmt.Errorf("encrypt binary failed: %w", err)
	}

	return cipherText, nil
}

func (params Payload) SKEASetKeyIv() (string, error) {
	if len(params.Key) != 16 || len(params.IV) != 16 {
		return "", fmt.Errorf("the length of key and iv should be greater equal to 16")
	}

	binRaw, err := BinReader(params.PlainText)
	if err != nil {
		return "", fmt.Errorf("<enc.BinReader()> err: %w", err)
	}

	cipherText, err := crypto.XorSm4HexBase85Encrypt(binRaw, params.Key, params.IV)
	if err != nil {
		return "", fmt.Errorf("encrypt binary failed: %w", err)
	}

	return cipherText, nil
}

func (params Payload) SignSetKeyIV(signFn interface{}) (string, error) {
	binRaw, err := BinReader(params.PlainText)
	if err != nil {
		return "", fmt.Errorf("<BinReader()> err: %w", err)
	}

	switch signFn := signFn.(type) {
	// 形参为: (cipherText, key []byte)
	case func([]byte, []byte) (string, error):
		cipherText, err := signFn(binRaw, params.Key)
		if err != nil {
			return "", fmt.Errorf("encrypt failed: %w", err)
		}
		return cipherText, nil

	// 形参为: (cipherText, key, iv []byte)
	case func([]byte, []byte, []byte) (string, error):
		if len(params.Key) != 16 || len(params.IV) != 16 {
			return "", fmt.Errorf("the length of key and iv should be greater equal to 16")
		}
		cipherText, err := signFn(binRaw, params.Key, params.IV)
		if err != nil {
			return "", fmt.Errorf("encrypt failed: %w", err)
		}
		return cipherText, nil

	default:
		return "", errors.New("unsupported encryption function signature")
	}
}
