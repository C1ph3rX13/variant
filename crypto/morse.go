package crypto

import (
	"strings"

	"github.com/alwindoss/morse"
)

func MorseEncode(input string) (string, error) {
	h := morse.NewHacker()
	morseCode, err := h.Encode(strings.NewReader(input))
	if err != nil {
		return "", err
	}
	return string(morseCode), nil
}

func MorseDecode(input string) ([]byte, error) {
	h := morse.NewHacker()
	morseCode, err := h.Decode(strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	return morseCode, nil
}
