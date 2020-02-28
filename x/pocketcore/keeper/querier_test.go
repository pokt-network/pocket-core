package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestQuerySupportedBlockchains(t *testing.T) {
	ctx, _, _, _, k := createTestInput(t, false)
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
	ctx, _, _, _, k := createTestInput(t, false)
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

//
//func TestQueryEvidence(t *testing.T) {
//	ctx, _, _, _, k := createTestInput(t, false)
//	appPrivateKey := getRandomPrivateKey()
//	appPubKey := appPrivateKey.PublicKey().RawString()
//	npk := getRandomPubKey()
//	ethereum, err := types.NonNativeChain{
//		Ticker:  "eth",
//		Netid:   "4",
//		Version: "v1.9.9",
//		Client:  "geth",
//		Inter:   "",
//	}.HashString()
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	// create a session header
//	validHeader := types.SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: 1,
//	}
//	storedEvidence := types.StoredEvidence{
//		SessionHeader:   validHeader,
//		ServicerAddress: npk.Address().String(),
//		TotalRelays:     2000,
//	}
//	k.SetEvidence(ctx, sdk.Address(npk.Address()), storedEvidence)
//	bz, er := types.ModuleCdc.MarshalJSON(types.QueryEvidenceParams{
//		Address: sdk.Address(npk.Address()),
//		Header: types.SessionHeader{
//			ApplicationPubKey:  appPubKey,
//			Chain:              ethereum,
//			SessionBlockHeight: 1000,
//		},
//	})
//	assert.Nil(t, er)
//	request := abci.RequestQuery{
//		Data:                 bz,
//		Path:                 types.QueryEvidence,
//		Height:               ctx.BlockHeight(),
//		Prove:                false,
//		XXX_NoUnkeyedLiteral: struct{}{},
//		XXX_unrecognized:     nil,
//		XXX_sizecache:        0,
//	}
//	resbz, err := queryEvidence(ctx, request, k)
//	assert.Nil(t, err)
//	var stored types.StoredEvidence
//	er = types.ModuleCdc.UnmarshalJSON(resbz, &stored)
//	assert.Nil(t, er)
//	assert.Equal(t, stored, storedEvidence)
//	// evidences query
//	var stored2 []types.StoredEvidence
//	bz2, er2 := types.ModuleCdc.MarshalJSON(types.QueryEvidencesParams{
//		Address: sdk.Address(npk.Address()),
//	})
//	assert.Nil(t, er2)
//	request2 := abci.RequestQuery{
//		Data:                 bz2,
//		Path:                 types.QueryEvidences,
//		Height:               ctx.BlockHeight(),
//		Prove:                false,
//		XXX_NoUnkeyedLiteral: struct{}{},
//		XXX_unrecognized:     nil,
//		XXX_sizecache:        0,
//	}
//	resbz2, err := queryEvidences(ctx, request2, k)
//	assert.Nil(t, err)
//	er = types.ModuleCdc.UnmarshalJSON(resbz2, &stored2)
//	assert.Nil(t, er)
//	assert.Equal(t, stored2, []types.StoredEvidence{storedEvidence})
//}
