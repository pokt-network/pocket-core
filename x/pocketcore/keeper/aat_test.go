package keeper

import (
	"testing"

	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	"github.com/stretchr/testify/assert"
)

func TestAATGeneration(t *testing.T) {
	passphrase := "test"
	kb := NewTestKeybase()
	kp, err := kb.Create(passphrase)
	assert.Nil(t, err)
	privkey, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	assert.Nil(t, err)
	appPubKey := kp.PublicKey
	res, err := AATGeneration(appPubKey.RawString(), appPubKey.RawString(), privkey)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Validate())
}
