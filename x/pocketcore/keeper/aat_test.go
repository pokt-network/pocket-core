package keeper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAATGeneration(t *testing.T) {
	passphrase := "test"
	kb := NewTestKeybase()
	kp, err := kb.Create(passphrase)
	assert.Nil(t, err)
	appPubKey := kp.PublicKey
	res, err := AATGeneration(appPubKey.RawString(),
		appPubKey.RawString(), passphrase, kb)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Nil(t, res.Validate())
}
