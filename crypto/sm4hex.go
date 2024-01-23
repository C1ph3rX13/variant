package crypto

import (
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm4"
	"log"
)

func Sm4CbcEncryptHex(rawData, keyHex, ivHex string) ([]byte, error) {
	plainText, err := hex.DecodeString(rawData)
	if err != nil {
		log.Fatalf("DecodeString cipherText Error: %v", err)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatalf("DecodeString key Error: %v", err)
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		log.Fatalf("DecodeString iv Error: %v", err)
	}

	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行加密
	cipherText, err := sm4.Sm4Cbc(key, plainText, true)
	if err != nil {
		log.Fatalf("Encrypt Error: %v", err)
	}

	return cipherText, nil
}

func Sm4CbcDecryptHex(cipherTextHex, keyHex, ivHex string) ([]byte, error) {
	// Decrypt: 解密的密文需要经过HEX编码后作为输入
	// Key, IV: 被HEX编码的长度为16位的字符串, 为兼容Key, IV中不可见字符
	cipherText, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		log.Fatalf("DecodeString cipherText Error: %v", err)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		log.Fatalf("DecodeString key Error: %v", err)
	}

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		log.Fatalf("DecodeString iv Error: %v", err)
	}

	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行解密
	plainText, err := sm4.Sm4Cbc(key, cipherText, false)
	if err != nil {
		log.Fatalf("DecodeString iv Error: %v", err)
	}

	return plainText, nil
}
