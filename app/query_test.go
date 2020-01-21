package app

import (
	"github.com/pokt-network/pocket-core/x/nodes"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestQueryBlock(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	_, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	height := int64(1)
	select {
	case <-evtChan:
		got, err := nodes.QueryBlock(getInMemoryTMClient(), &height)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		return
	}
}

func TestQueryChainHeight(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		got, err := nodes.QueryChainHeight(memCli)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), got) // should not be 0 due to empty blocks
		return
	}
}

func TestQueryTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000))
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		got, err := nodes.QueryTransaction(memCli, tx.TxHash)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		return
	}
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
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		got, err := nodes.QueryValidators(memCodec(), memCli, 1)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
}

func TestQueryValidator(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	if err != nil {
		t.Fatal(err)
	}
	codec := memCodec()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		got, err := nodes.QueryValidator(codec, memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.Equal(t, cb.GetAddress(), got.Address)
		assert.False(t, got.Jailed)
		assert.True(t, got.StakedTokens.Equal(sdk.NewInt(10000000)))
	}
}

func TestQueryDaoBalance(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		var err error
		got, err := nodes.QueryDAO(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, big.NewInt(0), got.BigInt())
		return
	}
}

func TestQuerySupply(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		gotStaked, gotUnstaked, err := nodes.QuerySupply(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.True(t, gotStaked.Equal(sdk.NewInt(10000000)))
		assert.True(t, gotUnstaked.Equal(sdk.NewInt(1000000000)))
		return
	}
}

func TestQueryPOSParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		got, err := nodes.QueryPOSParams(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100000), got.MaxValidators)
		assert.Equal(t, int64(1), got.StakeMinimum)
		assert.Equal(t, int8(90), got.ProposerRewardPercentage)
		assert.Equal(t, "stake", got.StakeDenom)
	}
}

func TestQueryStakedValidator(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		got, err := nodes.QueryStakedValidators(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
}

func TestAccountBalance(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	select {
	case <-evtChan:
		var err error
		got, err := nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, got.Int64(), int64(1000000000))
	}
}
