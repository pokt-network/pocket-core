// This package is for all 'session' related code.
package session

import "github.com/pokt-network/pocket-core/crypto"

/*
"SessionKeyAlgo" function determines the sessionKey.
 */
func SessionKeyAlgo(devID string) []byte {
	// Simulate block hash
	b1 := "block 1"
	b2 := "block 2"
	block1 := crypto.SessionHash([]byte(b1))
	block2 := crypto.SessionHash([]byte(b2))
	// Get Developer ID Bytes
	dIDBytes := []byte(devID)
	key := []byte{}
	// Create the publicly verifiable key for the algorithm
	key = append(key, block1...)
	key = append(key, block2...)
	key = append(key, dIDBytes...)
	// Run through hashing algorithm
	return crypto.SessionHash(key)
}
