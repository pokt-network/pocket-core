// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"github.com/pokt-network/pocket-core/const"
)

// "session.go" specifies session related code for the crypto package.

/*
"SessionHash" returns the <SessionHashingAlgorithm> hash of a byte array.
*/
func SessionHash(s []byte) []byte { // hashing algorithm of the session
	hasher := _const.SessionHashingAlgorithm.New()
	hasher.Write(s)
	return hasher.Sum(nil)
}

/*
"SessionNonce" generates a 32 byte random key. (Unused for now)
*/
func SessionNonce() []byte { // one time random number from session
	return RandBytes(32)
}
