package session

import (
	"github.com/pokt-network/pocket-core/crypto"
)

// The key that service nodes are XOR'ed against
type SessionKey []byte

// Generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func NewSessionKey(developer SessionDeveloper, nonNativeChain SessionBlockchain, blockID SessionBlockID) (SessionKey, error) {
	// get the public key from the developer structure
	devPubKey := developer.PubKey.Bytes()
	// check for empty params
	if len(devPubKey) == 0 {
		return nil, EmptyDevPubKeyError
	}
	if len(nonNativeChain) == 0 {
		return nil, EmptyNonNativeChainError
	}
	if len(blockID.Hash.Bytes()) == 0 {
		return nil, EmptyBlockIDError
	}
	// append them all together
	// in the order of devPubKey - > nonnativeChain -> blockID
	// TODO consider using amino buffer to find the session key
	seed := append(devPubKey, nonNativeChain...)
	seed = append(seed, blockID.Hash.Bytes()...)

	// return the hash of the result
	return crypto.Hash(seed), nil
}
