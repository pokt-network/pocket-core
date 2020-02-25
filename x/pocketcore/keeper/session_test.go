package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestKeeper_Dispatch(t *testing.T) {
//	ctx, _, _, _, keeper := createTestInput(t, false)
//	appPrivateKey := getRandomPrivateKey()
//	appPubKey := appPrivateKey.PublicKey().RawString()
//	ethereum, err := types.NonNativeChain{
//		Ticker:  "eth",
//		Netid:   "4",
//		Version: "v1.9.9",
//		Client:  "geth",
//		Inter:   "",
//	}.HashString()
//	bitcoin, err := types.NonNativeChain{
//		Ticker:  "btc",
//		Netid:   "1",
//		Version: "0.19.0.1",
//		Client:  "",
//		Inter:   "",
//	}.HashString()
//	if err != nil {
//		t.Fatalf(err.Error())
//	}
//	// create a session header
//	validHeader := types.SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              ethereum,
//		SessionBlockHeight: 976,
//	}
//	// create an invalid session header
//	invalidHeader := types.SessionHeader{
//		ApplicationPubKey:  appPubKey,
//		Chain:              bitcoin,
//		SessionBlockHeight: 976,
//	}
//	res, err := keeper.Dispatch(ctx, validHeader)
//	assert.Nil(t, err)
//	assert.Equal(t, res.Chain, ethereum)
//	assert.Equal(t, res.SessionBlockHeight, int64(976))
//	assert.Equal(t, res.ApplicationPubKey, appPubKey)
//	assert.Equal(t, res.SessionHeader, validHeader)
//	assert.Len(t, res.SessionNodes, 5)
//	_, err = keeper.Dispatch(ctx, invalidHeader)
//	assert.NotNil(t, err)
//}

func TestKeeper_IsSessionBlock(t *testing.T) {
	notSessionContext, _, _, _, keeper := createTestInput(t, false)
	assert.False(t, keeper.IsSessionBlock(notSessionContext))
}

//func TestKeeper_GetLatestSessionBlock(t *testing.T) {
//	notSessionContext, _, _, _, keeper := createTestInput(t, false)
//	assert.Equal(t, keeper.GetLatestSessionBlock(notSessionContext).BlockHeight(), int64(976))
//}

func TestKeeper_IsPocketSupportedBlockchain(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	sb := []string{"ethereum"}
	notSB := "bitcoin"
	p := types.Params{
		SupportedBlockchains: sb,
	}
	keeper.SetParams(ctx, p)
	assert.True(t, keeper.IsPocketSupportedBlockchain(ctx, "ethereum"))
	assert.False(t, keeper.IsPocketSupportedBlockchain(ctx, notSB))
}
