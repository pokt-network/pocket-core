package app

import (
	"github.com/pokt-network/pocket-core/x/nodes"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendTransaction(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	res, err := nodes.Send(memCodec(), getInMemoryTMClient(), kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000))
	assert.Nil(t, err)
	assert.NotNil(t, res)
	// todo assert that accountValue of coinbase is `1000 less` than account value before
	// todo assert that receiver accountValue = 1000
}
