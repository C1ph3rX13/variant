package crypto

import (
	"encoding/hex"

	ecies "github.com/ecies/go/v2"
)

func EllipticCurveEncrypt(plaintext []byte, privKey []byte) ([]byte, error) {
	hexPrivKey := hex.EncodeToString(privKey)
	k, err := ecies.NewPrivateKeyFromHex(hexPrivKey)
	if err != nil {
		return nil, err
	}

	ciphertext, err := ecies.Encrypt(k.PublicKey, plaintext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func EllipticCurveDecrypt(ciphertext []byte, privKey []byte) ([]byte, error) {
	hexPrivKey := hex.EncodeToString(privKey)
	k, err := ecies.NewPrivateKeyFromHex(hexPrivKey)
	if err != nil {
		return nil, err
	}

	dec, err := ecies.Decrypt(k, ciphertext)
	if err != nil {
		return nil, err
	}

	return dec, nil
}
