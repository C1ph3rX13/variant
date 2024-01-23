package crypto

import (
	"github.com/eknkc/basex"
)

var base85 *basex.Encoding

func init() {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~"
	base85, _ = basex.NewEncoding(alphabet)
}

func Base85Encode(plainText []byte) (string, error) {
	// Base85 编码
	cipherText := base85.Encode(plainText)
	return cipherText, nil
}

func Base85Decode(cipherText string) ([]byte, error) {
	// Base85 解码
	plainText, err := base85.Decode(cipherText)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
