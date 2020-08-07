package auth

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignatureDepth(t *testing.T) {
	msg := []byte("foo")
	// 3 keys should have 5 signatures
	pk1 := crypto.GenerateEd25519PrivKey()
	pk2 := crypto.GenerateEd25519PrivKey()
	pk3 := crypto.GenerateEd25519PrivKey()
	pub1 := pk1.PublicKey()
	pub2 := pk2.PublicKey()
	pub3 := pk3.PublicKey()
	// sign
	sig1, _ := pk1.Sign(msg)
	sig2, _ := pk2.Sign(msg)
	sig3, _ := pk3.Sign(msg)
	// multisig 0 multisigpub0
	ms0 := crypto.MultiSignature{Sigs: [][]byte{sig1, sig2}}
	mspk0 := crypto.PublicKeyMultiSignature{PublicKeys: []crypto.PublicKey{pub1, pub2}}
	// multisig and multisigpub
	_ = crypto.MultiSignature{Sigs: [][]byte{ms0.Marshal(), sig3}}
	mspk := crypto.PublicKeyMultiSignature{PublicKeys: []crypto.PublicKey{mspk0, pub3}}
	assert.True(t, ValidateSignatureDepth(5, mspk))
	assert.False(t, ValidateSignatureDepth(4, mspk))
}
