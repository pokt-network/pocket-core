package crypto
import (
	crypt "github.com/tendermint/tendermint/crypto"
)

// wrapper around tendermint public key
type PublicKey crypt.PubKey
