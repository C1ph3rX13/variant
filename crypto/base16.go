package crypto

import (
	"github.com/eknkc/basex"
)

var base16 *basex.Encoding

func init() {
	alphabet := "0123456789abcdef"
	base16, _ = basex.NewEncoding(alphabet)
}

func Base16Encode(plainText []byte) (string, error) {
	// base16 编码
	cipherText := base16.Encode(plainText)
	return cipherText, nil
}

func Base16Decode(cipherText string) ([]byte, error) {
	// base16 解码
	plainText, err := base16.Decode(cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
