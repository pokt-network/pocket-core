package crypto

import (
	crypt "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// wrapper around tendermint public key
type PublicKey crypt.PubKey

// wrapper around tendermint private key
type PrivateKey crypt.PrivKey

func NewKeypair() (privateKey PrivateKey, publicKey PublicKey) {
	privateKey = NewPrivateKey()
	publicKey = privateKey.PubKey()
	return
}

func NewPrivateKey() PrivateKey {
	return secp256k1.GenPrivKey()
}

// todo need function to convert hex string or bytes into public key
