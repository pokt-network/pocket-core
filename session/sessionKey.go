// This package is for all 'session' related code.
package session

import "github.com/pokt-network/pocket-core/crypto"

// "sessionKey.go" defines sessionKey specific methods

/*
"GenerateSessionKey" function determines the sessionKey.
 */
func GenerateSessionKey(devID string) []byte {
	b1 := "block 1"								// Simulate block hash
	b2 := "block 2"
	block1 := crypto.SessionHash([]byte(b1))
	block2 := crypto.SessionHash([]byte(b2))
	dIDBytes := []byte(devID)					// Get Developer ID Bytes
	key := []byte{}
	key = append(key, block1...)				// Create the publicly verifiable key for the algorithm
	key = append(key, block2...)
	key = append(key, dIDBytes...)
	return crypto.SessionHash(key)				// Run through hashing algorithm
}
