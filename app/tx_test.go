package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/types"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
	"testing"
	"time"
)

func TestUnstakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = apps.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test")
	}
	cleanup()
	stopCli()
}

func TestUnstakeNode(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	kp := *getUnstakedAccount(kb)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = nodes.StakeTx(memCodec(), memCli, kb, []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}, "https://myPocketNode:8080", sdk.NewInt(10000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = nodes.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		memCli, stopCli, evtChan = subscribeNewblock(t)
		for {
			select {
			case blck := <-evtChan:
				if blck.Data.(tmTypes.EventDataNewBlock).Block.Height > 25 { // validator isn't unstaked until after session ends
					got, err := nodes.QueryUnstakingValidators(memCodec(), memCli, 0)
					assert.Nil(t, err)
					assert.Equal(t, 1, len(got))
					cleanup()
					stopCli()
					return
				}
			default:
				continue
			}
		}
	}
}

func TestStakeNode(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	kp := *getUnstakedAccount(kb)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode:8080", sdk.NewInt(10000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := nodes.QueryStakedValidators(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(got))
	}
	cleanup()
	stopCli()
}

func TestStakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeNewTx(t)
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
	stopCli()
	cleanup()
}

func TestSendTransaction(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	var baseAmount = sdk.NewInt(1000000000)
	var transferAmount = sdk.NewInt(1000)
	var tx *sdk.TxResponse

	select {
	case <-evtChan:
		var err error
		tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", transferAmount)
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		validator, err := nodes.QueryAccountBalance(memCodec(), memCli, kp.GetAddress(), 0)
		assert.Nil(t, err)
		assert.True(t, validator.Equal(transferAmount))
		validator, err = nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.True(t, validator.Equal(baseAmount.Sub(transferAmount)))
	}
	cleanup()
	stopCli()
}

func TestSendRawTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	// create the transaction
	txBz, err := types.DefaultTxEncoder(memCodec())(types.NewTestTx(sdk.Context{}.WithChainID("pocket-test"),
		[]sdk.Msg{bank.MsgSend{
			FromAddress: cb.GetAddress(),
			ToAddress:   kp.GetAddress(),
			Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))),
		}},
		[]crypto.PrivateKey{pk},
		[]uint64{0},
		[]uint64{0},
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)))))
	assert.Nil(t, err)
	select {
	case <-evtChan:
		var err error
		txResp, err := nodes.RawTx(memCodec(), memCli, cb.GetAddress(), txBz)
		assert.Nil(t, err)
		assert.NotNil(t, txResp)
	}
	select {
	case <-evtChan: // todo needs empty block?
	}
	// next block
	select {
	case <-evtChan:
		res, err := nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.Equal(t, int64(999999999), res.Int64())
	}
	cleanup()
	stopCli()
}
