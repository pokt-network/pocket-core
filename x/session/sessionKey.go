package session

import (
	"github.com/pokt-network/pocket-core/crypto"
	model "github.com/pokt-network/pocket-core/types"
)

// The key that service nodes are XOR'ed against
type SessionKey []byte

// Generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func NewSessionKey(devPubKey model.AminoBuffer, nonNativeChain model.AminoBuffer, blockID model.AminoBuffer) (SessionKey, error) {
	if len(devPubKey) == 0 {
		return nil, EmptyDevPubKeyError
	}
	if len(nonNativeChain) == 0 {
		return nil, EmptyNonNativeChainError
	}
	if len(blockID) == 0 {
		return nil, EmptyBlockIDError
	}
	seed := devPubKey.Append(nonNativeChain, blockID)

	return crypto.Hash(seed), nil
}
