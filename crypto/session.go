// This package is for cryptography that is used in Pocket Core.
package crypto

import (
	"github.com/pokt-network/pocket-core/const"
)

/*
"SessionHash" returns the <SessionHashingAlgorithm> hash of a byte array.
 */
func SessionHash (s []byte) []byte {
	hasher := _const.SessionHashingAlgorithm.New()
	hasher.Write(s)
	return hasher.Sum(nil)
}

/*
"SessionNonce" generates a 32 byte random key. (Unused for now)
 */
func SessionNonce() []byte{
	return RandBytes(32)
}




