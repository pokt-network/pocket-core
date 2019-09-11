package crypto

import (
	"crypto/elliptic"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func S256() elliptic.Curve {
	return secp256k1.S256()
}
