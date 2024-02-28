package crypto

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hex(data []byte) string {
	return hex.EncodeToString(Sha256(data))
}

func Sha256(data []byte) []byte {
	digest := sha256.New()
	digest.Write(data)
	return digest.Sum(nil)
}

func Sha1(data []byte) []byte {
	digest := sha1.New()
	digest.Write(data)
	return digest.Sum(nil)
}

func Sha1Hex(data []byte) string {
	return hex.EncodeToString(Sha1(data))
}
