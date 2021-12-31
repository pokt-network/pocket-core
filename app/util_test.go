package app

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov"
	"github.com/pokt-network/pocket-core/x/nodes"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
)

func TestBuildSignMultisig(t *testing.T) {
	codec.UpgradeHeight = 7000
	_, kb, cleanup := NewInMemoryTendermintNodeAmino(t, oneAppTwoNodeGenesis())
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
	bz, err := gov.BuildAndSignMulti(memCodec(), cb.GetAddress(), pms, &msg, getInMemoryTMClient(), kb, "test", 10000000, true)
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp2.GetAddress(), bz, kps, getInMemoryTMClient(), kb, "test", true)
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp3.GetAddress(), bz, nil, getInMemoryTMClient(), kb, "test", true)
	assert.Nil(t, err)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	<-evtChan // Wait for block
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), sdk.Address(pms.Address()), "test", sdk.NewInt(100000000), true)
	fmt.Println("HERE: ", tx)
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-evtChan // Wait for tx
	txRaw, err := nodes.RawTx(memCodec(), memCli, sdk.Address(pms.Address()), bz)
	assert.Nil(t, err)
	fmt.Println(txRaw)
	assert.Zero(t, txRaw.Code)

	cleanup()
	stopCli()
}
