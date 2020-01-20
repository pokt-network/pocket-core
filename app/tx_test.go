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
	memCli := getInMemoryTMClient()
	memCodec := memCodec()

	fromAddr := cb.GetAddress()
	toAddr := kp.GetAddress()

	res, err := nodes.Send(memCodec, memCli, kb, fromAddr, toAddr, "test", sdk.NewInt(1000))
	assert.Nil(t, err)
	assert.NotNil(t, res)

	got, err := nodes.QueryAccountBalance(memCodec, memCli, toAddr, res.Height)
	assert.Nil(t, err)
	assert.True(t, got.Equal(sdk.NewInt(1000)))
	// todo assert that accountValue of coinbase is `1000 less` than account value before
	// todo assert that receiver accountValue = 1000
}
