package session

import (
	"github.com/pokt-network/pocket-core/crypto"
)

// The key that service nodes are XOR'ed against
type SessionKey []byte

// Generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func NewSessionKey(app SessionAppPubKey, nonNativeChain SessionBlockchain, blockID SessionBlockID) (SessionKey, error) {
	// validate session application
	if err := app.Validate(); err != nil {
		return nil, err
	}
	// get the public key from the app structure
	appPubKey, err := app.Bytes()
	if err != nil {
		return nil, err
	}
	if err = nonNativeChain.Validate(); err != nil {
		return nil, err
	}
	if err = blockID.Validate(); err != nil {
		return nil, err
	}
	// append them all together
	// in the order of appPubKey - > nonnativeChain -> blockID
	// TODO consider using amino buffer to find the session key
	seed := append(appPubKey, nonNativeChain...)
	seed = append(seed, blockID.Hash...)

	// return the hash of the result
	return crypto.SHA3FromBytes(seed), nil
}

func (sk SessionKey) Validate() error {
	// todo more validation
	if len(sk) == 0 {
		return EmptySessionKeyError
	}
	return nil
}
