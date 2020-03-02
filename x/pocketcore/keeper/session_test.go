package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_Dispatch(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	bitcoin, err := types.NonNativeChain{
		Ticker:  "btc",
		Netid:   "1",
		Version: "0.19.0.1",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	// create a session header
	validHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 976,
	}
	// create an invalid session header
	invalidHeader := types.SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              bitcoin,
		SessionBlockHeight: 976,
	}
	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["pos"]).Return(ctx.KVStore(keys["pos"]))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("MustGetPrevCtx", validHeader.SessionBlockHeight).Return(ctx)
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())
	res, err := keeper.Dispatch(mockCtx, validHeader)
	assert.Nil(t, err)
	assert.Equal(t, res.SessionHeader.Chain, ethereum)
	assert.Equal(t, res.SessionHeader.SessionBlockHeight, int64(976))
	assert.Equal(t, res.SessionHeader.ApplicationPubKey, appPubKey)
	assert.Equal(t, res.SessionHeader, validHeader)
	assert.Len(t, res.SessionNodes, 5)
	_, err = keeper.Dispatch(mockCtx, invalidHeader)
	assert.NotNil(t, err)
}

func TestKeeper_IsSessionBlock(t *testing.T) {
	notSessionContext, _, _, _, keeper, _ := createTestInput(t, false)
	assert.False(t, keeper.IsSessionBlock(notSessionContext))
}

func TestKeeper_GetLatestSessionBlock(t *testing.T) {
	ctx, _, _, _, keeper, keys := createTestInput(t, false)

	mockCtx := new(Ctx)
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())
	mockCtx.On("MustGetPrevCtx", keeper.GetLatestSessionBlockHeight(mockCtx)).Return(ctx)
	assert.Equal(t, ctx.BlockHeight(), keeper.GetLatestSessionBlock(mockCtx).BlockHeight())
}

func TestKeeper_IsPocketSupportedBlockchain(t *testing.T) {
	ctx, _, _, _, keeper, _ := createTestInput(t, false)
	sb := []string{"ethereum"}
	notSB := "bitcoin"
	p := types.Params{
		SupportedBlockchains: sb,
	}
	keeper.SetParams(ctx, p)
	assert.True(t, keeper.IsPocketSupportedBlockchain(ctx, "ethereum"))
	assert.False(t, keeper.IsPocketSupportedBlockchain(ctx, notSB))
}
