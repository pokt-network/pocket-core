package app

import (
	"fmt"
	types2 "github.com/pokt-network/pocket-core/x/auth/types"
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
	bz, err := gov.BuildAndSignMulti(memCodec(), cb.GetAddress(), pms, &msg, getInMemoryTMClient(), kb, "test", 10000000)
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp2.GetAddress(), bz, kps, getInMemoryTMClient(), kb, "test")
	assert.Nil(t, err)
	bz, err = gov.SignMulti(memCodec(), kp3.GetAddress(), bz, nil, getInMemoryTMClient(), kb, "test")
	assert.Nil(t, err)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	<-evtChan // Wait for block
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), sdk.Address(pms.Address()), "test", sdk.NewInt(100000000))
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

func TestAminoToProtoTryCatch(t *testing.T) {
	cdc := memCodec()
	priv := crypto.GenerateEd25519PrivKey()
	pk := priv.PublicKey()
	account := types2.BaseAccount{
		Address: sdk.Address(pk.Address()),
		Coins:   sdk.NewCoins(sdk.NewCoin("pokt", sdk.NewInt(20))),
		PubKey:  pk,
	}
	// let's try amino bytes in the world state
	aminoBz, err := cdc.LegacyMarshalBinaryBare(&account)
	if err != nil {
		t.Fatalf("unable to marshal legacy: %s", err)
	}
	// ensure the marshaller works well
	aminoBz2, err := cdc.MarshalBinaryBare(&account)
	if err != nil {
		t.Fatalf("unable to marshal: %s", err)
	}
	assert.Equal(t, aminoBz, aminoBz2)
	// set upgrade after true
	cdc.SetAfterUpgradeMod(true)
	var res types2.BaseAccount
	// unmarshal amino bz after the upgrade
	err = cdc.UnmarshalBinaryBare(aminoBz, &res)
	if err != nil {
		t.Fatalf("unable to unmarshalBinaryBare: %s", err)
	}
	assert.Equal(t, account, res)
	// reset upgrade after
	cdc.SetAfterUpgradeMod(false)
	// again with binary length prefix
	aminoBz, err = cdc.LegacyMarshalBinaryLengthPrefixed(&account)
	if err != nil {
		t.Fatalf("unable to marshal legacy: %s", err)
	}
	// ensure the marshaller works well
	aminoBz2, err = cdc.MarshalBinaryLengthPrefixed(&account)
	if err != nil {
		t.Fatalf("unable to marshal: %s", err)
	}
	assert.Equal(t, aminoBz, aminoBz2)
	res = types2.BaseAccount{}
	// set upgrade after true
	cdc.SetAfterUpgradeMod(true)
	// unmarshal amino bz after the upgrade
	err = cdc.UnmarshalBinaryLengthPrefixed(aminoBz, &res)
	if err != nil {
		t.Fatalf("unable to unmarshalbinarylengthprefix: %s", err)
	}
	assert.Equal(t, account, res)
	// lets mix up some pointers shall we?
	// resset upgrade after
	cdc.SetAfterUpgradeMod(false)
	// let's try amino bytes in the world state
	aminoBz, err = cdc.LegacyMarshalBinaryBare(account)
	if err != nil {
		t.Fatalf("unable to marshal legacy: %s", err)
	}
	// ensure the marshaller works well
	aminoBz2, err = cdc.MarshalBinaryBare(&account)
	if err != nil {
		t.Fatalf("unable to marshal: %s", err)
	}
	assert.Equal(t, aminoBz, aminoBz2)
	// set upgrade after true
	res = types2.BaseAccount{}
	// set upgrade after true
	cdc.SetAfterUpgradeMod(true)
	// unmarshal amino bz after the upgrade
	err = cdc.UnmarshalBinaryBare(aminoBz, &res)
	if err != nil {
		t.Fatalf("unable to unmarshal mixed pointers: %s", err)
	}
	assert.Equal(t, res, account)
}
