package crypto

import (
	"crypto"
)

// TODO convert to standard crypto library
func Hash(input []byte) []byte {
	hasher := crypto.SHA3_256.New()
	hasher.Write(input)
	return hasher.Sum(nil)
}
