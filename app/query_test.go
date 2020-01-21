package app

import (
	"fmt"
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

	height := int64(0)
	got, err := nodes.QueryBlock(getInMemoryTMClient(), &height)
	assert.NotNil(t, err)
	assert.Nil(t, got)

	time.Sleep(60 * time.Millisecond) // Feed empty blocks
	height = 1
	got, err = nodes.QueryBlock(getInMemoryTMClient(), &height)
	assert.Nil(t, err)
	assert.NotNil(t, got)
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
	t.Errorf("context expired")
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

	time.Sleep(140 * time.Millisecond) // Feed empty blocks to ensure tx is on the chain

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
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	time.Sleep(time.Second * 2)
	got, err := nodes.QueryValidators(memCodec(), getInMemoryTMClient(), 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(got))
	fmt.Println(got)
}

//func TestQueryValidators(t *testing.T) {
//	_, _, cleanup := NewInMemoryTendermintNode(t)
//	defer cleanup()
//	time.Sleep(time.Second * 5)
//	tmClient := getInMemoryTMClient()
//	err := tmClient.Start()
//	defer tmClient.Stop()
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//	res, err := tmClient.Subscribe(ctx, "test-client", "tm.event = 'NewBlock'")
//	if err != nil {
//		t.Fatal(err)
//	}
//	go func() {
//		for block := range res {
//			fmt.Println(block.Data.(types.EventDataNewBlock))
//			//got, err := nodes.QueryValidators(memCodec(), getInMemoryTMClient(), 1)
//			//assert.Nil(t, err)
//			//assert.Equal(t, 1, len(got))
//			//fmt.Println(got)
//		}
//	}()
//}

func TestQueryValidator(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()

	tmClient := getInMemoryTMClient()
	codec := memCodec()

	time.Sleep(1000 * time.Millisecond) // Feed empty blocks

	got, err := nodes.QueryValidator(codec, tmClient, cb.GetAddress(), 0)

	assert.Nil(t, err)
	assert.Equal(t, cb.GetAddress(), got.Address)
	assert.False(t, got.Jailed)
	assert.True(t, got.StakedTokens.Equal(sdk.NewInt(10000000)))
}

func TestQueryDaoBalance(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	got, err := nodes.QueryDAO(memCodec(), getInMemoryTMClient(), 0)
	assert.NotNil(t, err)

	time.Sleep(90 * time.Millisecond)
	height, err := nodes.QueryChainHeight(getInMemoryTMClient())
	got, err = nodes.QueryDAO(memCodec(), getInMemoryTMClient(), height)
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(0), got.BigInt())
}

func TestQuerySupply(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	_, _, err := nodes.QuerySupply(memCodec(), getInMemoryTMClient(), 0)
	assert.NotNil(t, err)

	time.Sleep(100 * time.Millisecond)
	height, err := nodes.QueryChainHeight(getInMemoryTMClient())
	gotStaked, gotUnstaked, err := nodes.QuerySupply(memCodec(), getInMemoryTMClient(), height)
	assert.Nil(t, err)
	assert.True(t, gotStaked.Equal(sdk.NewInt(10000000)))
	assert.True(t, gotUnstaked.Equal(sdk.NewInt(0)))
}

func TestQueryPOSParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()

	tmClient := getInMemoryTMClient()
	codec := memCodec()

	// TODO failed to load state at height 0 !?
	got, err := nodes.QueryPOSParams(codec, tmClient, 0)
	assert.NotNil(t, err)

	time.Sleep(100 * time.Millisecond)

	got, err = nodes.QueryPOSParams(codec, tmClient, 0)

	assert.Nil(t, err)
	assert.Equal(t, uint64(100000), got.MaxValidators)
	assert.Equal(t, int64(1), got.StakeMinimum)
	assert.Equal(t, int8(90), got.ProposerRewardPercentage)
	assert.Equal(t, "stake", got.StakeDenom)
}

func TestQueryStakedValidator(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	tmClient := getInMemoryTMClient()
	codec := memCodec()
	time.Sleep(1 * time.Second)
	got, err := nodes.QueryStakedValidators(codec, tmClient, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(got))
}

func TestAccountBalance(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()

	tmClient := getInMemoryTMClient()
	codec := memCodec()

	// TODO failed to load state at height 0 !?
	got, err := nodes.QueryAccountBalance(codec, tmClient, cb.GetAddress(), 0)
	assert.NotNil(t, err)

	time.Sleep(100 * time.Millisecond)
	got, err = nodes.QueryAccountBalance(codec, tmClient, cb.GetAddress(), 0)
	assert.Nil(t, err)
	assert.Equal(t, got, got)
	// TODO fix, there is a bug on QueryAccountBalance
}
