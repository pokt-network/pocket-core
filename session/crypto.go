package session

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/crypto"
)

// "Hash" returns the <SESSIONHASHINGALGORITHM> hash of a byte array.
func Hash(s []byte) []byte {
	hasher := _const.SESSIONHASHINGALGORITHM.New()
	hasher.Write(s)
	return hasher.Sum(nil)
}

// "Nonce" generates a 32 byte random key. (Unused for now)
func Nonce() ([]byte, error) {
	return crypto.RandBytes(32)
}

// "Key" function determines the sessionKey.
func Key(devID string) []byte {
	b1 := "block 1"
	b2 := "block 2"
	block1 := Hash([]byte(b1))
	block2 := Hash([]byte(b2))
	dIDBytes := []byte(devID)
	key := []byte{}
	key = append(key, block1...)
	key = append(key, block2...)
	key = append(key, dIDBytes...)
	return Hash(key)
}
