package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerySupportedBlockchains(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	p := types.Params{
		SupportedBlockchains: []string{"ethereum"},
	}
	k.SetParams(ctx, p)
	sbbz, err := querySupportedBlockchains(ctx, abci.RequestQuery{}, k)
	assert.Nil(t, err)
	var sb []string
	er := makeTestCodec().UnmarshalJSON(sbbz, &sb)
	assert.Nil(t, er)
	assert.Equal(t, sb, []string{"ethereum"})
}

func TestQueryParameters(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	p := types.Params{
		SupportedBlockchains: []string{"ethereum"},
	}
	k.SetParams(ctx, p)
	sbbz, err := queryParameters(ctx, k)
	assert.Nil(t, err)
	var params types.Params
	er := makeTestCodec().UnmarshalJSON(sbbz, &params)
	assert.Nil(t, er)
	assert.Equal(t, params, p)
}

func TestQueryReceipt(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	ethereum := hex.EncodeToString([]byte{01})
	// create a session header
	validHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 976,
	}
	receipt := types.Receipt{
		SessionHeader:   validHeader,
		ServicerAddress: npk.Address().String(),
		Total:           2000,
		EvidenceType:    types.RelayEvidence,
	}
	addr := sdk.Address(sdk.Address(npk.Address()))
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", k.storeKey).Return(ctx.KVStore(k.storeKey))
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("Logger").Return(ctx.Logger())
	er := k.SetReceipt(mockCtx, addr, receipt)
	if er != nil {
		t.Fatal(er)
	}
	bz, er := types.ModuleCdc.MarshalJSON(types.QueryReceiptParams{
		Address: sdk.Address(npk.Address()),
		Header:  validHeader,
		Type:    "relay",
	})
	assert.Nil(t, er)
	request := abci.RequestQuery{
		Data:   bz,
		Path:   types.QueryReceipt,
		Height: ctx.BlockHeight(),
	}
	resbz, err := queryReceipt(ctx, request, k)
	assert.Nil(t, err)
	var stored types.Receipt
	er = types.ModuleCdc.UnmarshalJSON(resbz, &stored)
	assert.Nil(t, er)
	assert.Equal(t, stored, receipt)
	// receipts query
	var stored2 []types.Receipt
	bz2, er2 := types.ModuleCdc.MarshalJSON(types.QueryReceiptsParams{
		Address: sdk.Address(npk.Address()),
	})
	assert.Nil(t, er2)
	request2 := abci.RequestQuery{
		Data:                 bz2,
		Path:                 types.QueryReceipt,
		Height:               ctx.BlockHeight(),
		Prove:                false,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
	resbz2, err := queryReceipts(ctx, request2, k)
	assert.Nil(t, err)
	er = types.ModuleCdc.UnmarshalJSON(resbz2, &stored2)
	assert.Nil(t, er)
	assert.Equal(t, stored2, []types.Receipt{receipt})
}
