package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSendTransaction(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var baseAmount sdk.Int = sdk.NewInt(1000000000)
	var transferAmount sdk.Int = sdk.NewInt(1000)
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
		return
	}
}

func TestSendRawTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx sdk.TxResponse

	select {
	case <-evtChan:
		var err error
		tx, err = nodes.RawTx(memCodec(), memCli, cb.GetAddress(), []byte{})
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
}

func TestStakeNode(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx *sdk.TxResponse
	var chains []string = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		tx, err = nodes.StakeTx(memCodec(), memCli, kb, chains, "https://myPocketNode:8080", sdk.NewInt(100000), kp, "test")
		// todo stake returns account does not exist even though it is contained within the CliCtx
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		got, err := nodes.QueryStakedValidators(memCodec(), memCli,0 )

		assert.Nil(t, err)
		assert.Equal(t, 2, len(got))
	}
}

func TestUnstakeNode(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx *sdk.TxResponse

	select {
	case <-evtChan:
		var err error
		tx, err = nodes.UnstakeTx(memCodec(), memCli, kb, cb.GetAddress(), "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		// todo unstake Tx marks success but it does not reflect itself on unstakedValidators
		got, err := nodes.QueryUnstakedValidators(memCodec(), memCli,0 )

		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
}

func TestStakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx *sdk.TxResponse
	var chains []string = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains,  sdk.NewInt(100000), kp, "test")
		// todo stake returns account does not exist even though it is contained within the CliCtx
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli,0)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
}
func TestUnstakeApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t)
	defer cleanup()
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeNewblock(t)
	defer stopCli()
	var tx *sdk.TxResponse
	var chains []string = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains,  sdk.NewInt(100000), kp, "test")
		// todo stake returns account does not exist even though it is contained within the CliCtx
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		tx, err = apps.UnstakeTx(memCodec(), memCli, kb, kp.GetAddress(), "test")

		assert.Nil(t, err)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplications(memCodec(), memCli,0)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(got))
	}
}
