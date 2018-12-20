// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"github.com/pokt-network/pocket-core/const"
)

// "session.go" specifies session related code for the crypto package.


/*
"SessionHash" returns the <SessionHashingAlgorithm> hash of a byte array.
 */
func SessionHash (s []byte) []byte {				// hashing algorithm of the session
	hasher := _const.SessionHashingAlgorithm.New()
	hasher.Write(s)
	return hasher.Sum(nil)
}

/*
"SessionNonce" generates a 32 byte random key. (Unused for now)
 */
func SessionNonce() []byte{							// one time random number from session
	return RandBytes(32)
}

/*
"GenerateSessionKey" function determines the sessionKey.
 */
func GenerateSessionKey(devID string) []byte {
	b1 := "block 1"								// Simulate block hash
	b2 := "block 2"
	block1 := SessionHash([]byte(b1))
	block2 := SessionHash([]byte(b2))
	dIDBytes := []byte(devID)					// Get Developer ID Bytes
	key := []byte{}
	key = append(key, block1...)				// Create the publicly verifiable key for the algorithm
	key = append(key, block2...)
	key = append(key, dIDBytes...)
	return SessionHash(key)						// Run through hashing algorithm
}


