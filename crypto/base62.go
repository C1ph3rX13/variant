package crypto

import (
	"github.com/eknkc/basex"
)

var base62 *basex.Encoding

func init() {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base62, _ = basex.NewEncoding(alphabet)
}

func Base62Encode(plainText []byte) (string, error) {
	// base62 编码
	cipherText := base62.Encode(plainText)
	return cipherText, nil
}

func Base62Decode(cipherText string) ([]byte, error) {
	// base62 解码
	plainText, err := base62.Decode(cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
