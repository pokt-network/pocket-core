package keeper

import (
	"github.com/pokt-network/posmint/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestAATGeneration(t *testing.T) {
	passphrase := "test"
	kb := NewTestKeybase()
	kp, err := kb.Create(passphrase)
	assert.Nil(t, err)
	appPubKey := kp.PubKey
	res, err := AATGeneration(crypto.PublicKey(appPubKey.(ed25519.PubKeyEd25519)).String(),
		crypto.PublicKey(appPubKey.(ed25519.PubKeyEd25519)).String(), passphrase, kb)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Validate())
}
