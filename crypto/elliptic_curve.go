package crypto

import (
	"encoding/hex"

	ecies "github.com/ecies/go/v2"
)

func EllipticCurveEncrypt(privKey []byte, plaintext []byte) ([]byte, error) {
	hexPrivKey := hex.EncodeToString(privKey)
	k, err := ecies.NewPrivateKeyFromHex(hexPrivKey)
	if err != nil {
		return []byte(""), err
	}

	ciphertext, err := ecies.Encrypt(k.PublicKey, plaintext)
	if err != nil {
		return []byte(""), err
	}

	return ciphertext, nil
}

func EllipticCurveDecrypt(privKey []byte, ciphertext []byte) ([]byte, error) {
	hexPrivKey := hex.EncodeToString(privKey)
	k, err := ecies.NewPrivateKeyFromHex(hexPrivKey)
	if err != nil {
		return []byte(""), err
	}

	dec, err := ecies.Decrypt(k, ciphertext)
	if err != nil {
		return []byte(""), err
	}

	return dec, nil
}
