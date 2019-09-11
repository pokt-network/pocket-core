package crypto

import (
	"golang.org/x/crypto/sha3"
)

// Converts []byte to SHA3-256 hashed []byte
func SHA3FromBytes(b []byte) []byte {
	hasher := sha3.New256()
	hasher.Write(b)
	return hasher.Sum(nil)
}

// Converts string to SHA3-256 hashed []byte
func SHA3FromString(s string) []byte {
	hasher := sha3.New256()
	hasher.Write([]byte(s))
	return hasher.Sum(nil)
}
