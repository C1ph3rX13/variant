package crypto

import (
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)

func Sm4CbcEncrypt(plainText, key, iv []byte) ([]byte, error) {
	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行加密
	cipherText, err := sm4.Sm4Cbc(key, plainText, true)
	if err != nil {
		return nil, fmt.Errorf("SM4CBC: encrypt error: %v", err)
	}

	return cipherText, nil
}

func Sm4CbcDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	// 设置 IV
	_ = sm4.SetIV(iv)

	// 设置 Key 进行解密
	plainText, err := sm4.Sm4Cbc(key, cipherText, false)
	if err != nil {
		return nil, fmt.Errorf("SM4CBC: decrypt error: %v", err)
	}

	return plainText, nil
}
