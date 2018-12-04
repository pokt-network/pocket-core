// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"math/rand"
)

// "util.go" specifies utility functions for the crypto package.

/*
"RandBytes" returns a random string of bytes.
*/
func RandBytes(n int) []byte { // generates random bytes from the seed specified
	output := make([]byte, n)   // create n bytes
	_, err := rand.Read(output) // read all random
	if err != nil {
		panic(err) // if error panic
	}
	return output // return the random bytes.
}
