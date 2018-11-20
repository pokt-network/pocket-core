// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"math/rand"
)

/*
"RandBytes" returns a random string of bytes.
 */
func RandBytes(n int) []byte {
	output := make([]byte, n)
	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)
	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	return output
}
