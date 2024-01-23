package crypto

import (
	"github.com/tjfoc/gmsm/sm4"
	"log"
)

func Sm4CbcEncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行加密
	cipherText, err := sm4.Sm4Cbc(key, plainText, true)
	if err != nil {
		log.Fatalf("Encrypt Error: %v", err)
	}

	return cipherText, nil
}

func Sm4CbcDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行解密
	plainText, err := sm4.Sm4Cbc(key, cipherText, false)
	if err != nil {
		log.Fatalf("DecodeString iv Error: %v", err)
	}

	return plainText, nil
}
