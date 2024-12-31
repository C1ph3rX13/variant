package crypto

import (
	"encoding/hex"
	"fmt"

	sgn "github.com/EgeBalci/sgn/pkg"
)

func SgnEncoder(file []byte) {
	encoder, err := sgn.NewEncoder(64)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the proper architecture
	err = encoder.SetArchitecture(64)
	if err != nil {
		return
	}

	// Encode the binary
	encodedBinary, err := encoder.Encode(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print out the hex dump of the encoded binary
	fmt.Println(hex.Dump(encodedBinary))
}
