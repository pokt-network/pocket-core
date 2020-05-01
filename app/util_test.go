package app

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
	"testing"
)

func TestBuildSignMultisig(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp2, err := kb.Create("test")
	assert.Nil(t, err)
	kp3, err := kb.Create("test")
	assert.Nil(t, err)
	kps := []crypto.PublicKey{cb.PublicKey, kp2.PublicKey, kp3.PublicKey}
	pms := crypto.PublicKeyMultiSignature{PublicKeys: kps}
	msg := types.MsgSend{
		FromAddress: sdk.Address(pms.Address()),
		ToAddress:   kp2.GetAddress(),
		Amount:      sdk.NewInt(1),
	}
	bz, err := gov.BuildAndSignMulti(memCodec(), cb.GetAddress(), pms, msg, getInMemoryTMClient(), kb, "test")
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp2.GetAddress(), bz, kps, getInMemoryTMClient(), kb, "test")
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp3.GetAddress(), bz, nil, getInMemoryTMClient(), kb, "test")
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), sdk.Address(pms.Address()), "test", sdk.NewInt(100000))
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		tx, err := nodes.RawTx(memCodec(), memCli, sdk.Address(pms.Address()), bz)
		assert.Nil(t, err)
		fmt.Println(tx)
		assert.Zero(t, tx.Code)
	}
	cleanup()
	stopCli()
}

func TestExportState(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		res, err := testPCA.ExportAppState(false, nil)
		assert.Nil(t, err)
		assert.NotNil(t, res)
	}
	cleanup()
	stopCli()
}
