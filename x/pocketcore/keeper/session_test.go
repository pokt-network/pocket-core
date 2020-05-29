package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_Dispatch(t *testing.T) {
	ctx, _, _, _, keeper, keys, _ := createTestInput(t, false)
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
	bitcoin := hex.EncodeToString([]byte{02})
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
	mockCtx.On("PrevCtx", validHeader.SessionBlockHeight).Return(ctx, nil)
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())
	mockCtx.On("Logger").Return(ctx.Logger())
	res, err := keeper.HandleDispatch(mockCtx, validHeader)
	assert.Nil(t, err)
	assert.Equal(t, res.Session.SessionHeader.Chain, ethereum)
	assert.Equal(t, res.Session.SessionHeader.SessionBlockHeight, int64(976))
	assert.Equal(t, res.Session.SessionHeader.ApplicationPubKey, appPubKey)
	assert.Equal(t, res.Session.SessionHeader, validHeader)
	assert.Len(t, res.Session.SessionNodes, 5)
	_, err = keeper.HandleDispatch(mockCtx, invalidHeader)
	assert.NotNil(t, err)
}

func TestKeeper_IsSessionBlock(t *testing.T) {
	notSessionContext, _, _, _, keeper, _, _ := createTestInput(t, false)
	assert.False(t, keeper.IsSessionBlock(notSessionContext.WithBlockHeight(977)))
}

func TestKeeper_IsPocketSupportedBlockchain(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	sb := []string{"ethereum"}
	notSB := "bitcoin"
	p := types.Params{
		SupportedBlockchains: sb,
	}
	keeper.SetParams(ctx, p)
	assert.True(t, keeper.IsPocketSupportedBlockchain(ctx, "ethereum"))
	assert.False(t, keeper.IsPocketSupportedBlockchain(ctx, notSB))
}
