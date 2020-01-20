package app

import (
	"github.com/pokt-network/pocket-core/x/nodes"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueryNodes(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	height := int64(0)
	got, err := nodes.QueryBlock(getInMemoryTMClient(), &height)
	assert.NotNil(t, err)
	assert.Nil(t, got)

	time.Sleep(60*time.Millisecond) // Feed empty blocks
	height = 1
	got, err = nodes.QueryBlock(getInMemoryTMClient(), &height)
	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestQueryChainHeight(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	got, err := nodes.QueryChainHeight(getInMemoryTMClient())
	assert.Nil(t, err)
	assert.Equal(t, got, int64(0))

	time.Sleep(50*time.Millisecond) // end block
	got, err = nodes.QueryChainHeight(getInMemoryTMClient())
	assert.Nil(t, err)
	assert.Equal(t, int64(1), got) // should not be 0 due to empty blocks
}

func TestQueryTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)

	memCli := getInMemoryTMClient()

	var hash string
	got, err := nodes.QueryTransaction(memCli, hash)
	assert.NotNil(t, err)
	assert.Nil(t, got)

	tx, err := nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000))
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	got, err = nodes.QueryTransaction(memCli, tx.TxHash)
	assert.NotNil(t, err) // Needs to be committed to the chain
	assert.Nil(t, got)

	time.Sleep(140 *time.Millisecond) // Feed empty blocks to ensure tx is on the chain

	got, err = nodes.QueryTransaction(memCli, tx.TxHash)
	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestQueryNodeStatus(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	got, err := nodes.QueryNodeStatus(getInMemoryTMClient())
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "pocket-test", got.NodeInfo.Network)
}

func TestQueryValidators(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	kp, err := kb.Create("test")
	assert.Nil(t, err)


	time.Sleep(70*time.Millisecond) // Feed empty blocks
	got, err := nodes.QueryValidators(memCodec(), getInMemoryTMClient(), 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(got))
}
