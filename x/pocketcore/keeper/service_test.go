package keeper

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func setupHandleRelayTest(t *testing.T) (
	ctx sdk.Ctx,
	keeper Keeper,
	kvkeys map[string]*sdk.KVStoreKey,
	clientPrivateKey, appPrivateKey crypto.Ed25519PrivateKey,
	nodePubKey crypto.PublicKey,
	chain string,
) {
	var kb keys.Keybase
	ctx, _, _, _, keeper, kvkeys, kb = createTestInput(t, false)

	chain = hex.EncodeToString([]byte{01})
	clientPrivateKey = getRandomPrivateKey()

	kp, _ := kb.GetCoinbase()
	nodePubKey = kp.PublicKey

	appPrivateKey = getRandomPrivateKey()

	appPubKey := appPrivateKey.PublicKey()
	app := appsTypes.NewApplication(
		sdk.Address(appPubKey.Address()),
		appPubKey,
		[]string{chain},
		sdk.NewInt(1000000000),
	)

	// Stake app
	ak := keeper.appKeeper.(appsKeeper.Keeper)
	app.MaxRelays = ak.CalculateAppRelays(ctx, app)
	ak.SetApplication(ctx, app)

	return
}

func testRelayAt(
	t *testing.T,
	ctx sdk.Ctx,
	keeper Keeper,
	clientBlockHeight int64,
	clientPrivateKey, appPrivateKey crypto.Ed25519PrivateKey,
	nodePubKey crypto.PublicKey,
	chain string,
) (*types.RelayResponse, sdk.Error) {
	clientPubKey := clientPrivateKey.PublicKey()
	appPubKey := appPrivateKey.PublicKey()

	blocksPerSesssion := keeper.BlocksPerSession(ctx)
	clientSessionHeight :=
		((clientBlockHeight-1)/blocksPerSesssion)*blocksPerSesssion + 1

	validRelay := types.Relay{
		Payload: types.Payload{
			Data: `{
			"jsonrpc":"2.0",
			"method":"web3_clientVersion",
			"params":[],
			"id":67
		}`,
			Method:  "",
			Path:    "",
			Headers: nil,
		},
		Meta: types.RelayMeta{BlockHeight: clientSessionHeight},
		Proof: types.RelayProof{
			Entropy:            rand.Int63(),
			SessionBlockHeight: clientSessionHeight,
			ServicerPubKey:     nodePubKey.RawString(),
			Blockchain:         chain,
			Token: types.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey.RawString(),
				ClientPublicKey:      clientPubKey.RawString(),
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}

	validRelay.Proof.RequestHash = validRelay.RequestHashString()
	appSig, er := appPrivateKey.Sign(validRelay.Proof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)

	clientSig, er := clientPrivateKey.Sign(validRelay.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	assert.Nil(t, er)
	validRelay.Proof.Signature = hex.EncodeToString(clientSig)

	httpClient := types.GetChainsClient()
	defer gock.Off() // Flush pending mocks after test execution
	defer gock.RestoreClient(httpClient)
	gock.InterceptClient(httpClient)
	gock.New("https://www.google.com:443").
		Post("/").
		Reply(200).
		BodyString("bar")
	return keeper.HandleRelay(ctx, validRelay)
}

func TestKeeper_HandleRelay(t *testing.T) {
	ctx, keeper, kvkeys, clientPrivateKey, appPrivateKey, nodePubKey, chain :=
		setupHandleRelayTest(t)

	// Store the original allowances to clean up at the end of this test
	originalClientBlockSyncAllowance := types.GlobalPocketConfig.ClientBlockSyncAllowance
	originalClientSessionSyncAllowance := types.GlobalPocketConfig.ClientSessionSyncAllowance

	// Eliminate the impact of ClientBlockSyncAllowance to avoid any relay meta validation errors (OutOfSyncError)
	types.GlobalPocketConfig.ClientBlockSyncAllowance = 10000

	nodeBlockHeight := ctx.BlockHeight()
	blocksPerSesssion := keeper.BlocksPerSession(ctx)
	latestSessionHeight := keeper.GetLatestSessionBlockHeight(ctx)

	t.Cleanup(func() {
		types.GlobalPocketConfig.ClientBlockSyncAllowance = originalClientBlockSyncAllowance
		types.GlobalPocketConfig.ClientSessionSyncAllowance = originalClientSessionSyncAllowance
		gock.Off() // Flush pending mocks after test execution
	})

	mockCtx := new(Ctx)
	mockCtx.On("KVStore", kvkeys["pos"]).Return(ctx.KVStore(kvkeys["pos"]))
	mockCtx.On("KVStore", kvkeys["params"]).Return(ctx.KVStore(kvkeys["params"]))
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())
	mockCtx.On("Logger").Return(ctx.Logger())
	mockCtx.On("PrevCtx", nodeBlockHeight).Return(ctx, nil)

	allSessionRangesTests := 4 // The range of block heights we will mock

	// Set up mocks for heights we'll query later.
	for i := int64(1); i <= blocksPerSesssion*int64(allSessionRangesTests); i++ {
		mockCtx.On("PrevCtx", nodeBlockHeight-i).Return(ctx, nil)
		mockCtx.On("PrevCtx", nodeBlockHeight+i).Return(ctx, nil)
	}

	// Client is synced with Node --> Success
	resp, err := testRelayAt(
		t,
		mockCtx,
		keeper,
		nodeBlockHeight,
		clientPrivateKey,
		appPrivateKey,
		nodePubKey,
		chain,
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp.Response, "bar")

	// TC 1:
	// Client is behind or advanced beyond Node's height with ClientSessionSyncAllowance 0
	// --> CodeInvalidBlockHeightError
	types.GlobalPocketConfig.ClientSessionSyncAllowance = 0
	for i := 1; i <= allSessionRangesTests; i++ {
		resp, err = testRelayAt(
			t,
			mockCtx,
			keeper,
			latestSessionHeight-blocksPerSesssion*int64(i),
			clientPrivateKey,
			appPrivateKey,
			nodePubKey,
			chain,
		)
		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, err.Codespace(), sdk.CodespaceType(types.ModuleName))
		assert.Equal(t, err.Code(), sdk.CodeType(types.CodeInvalidBlockHeightError))
		resp, err = testRelayAt(
			t,
			mockCtx,
			keeper,
			latestSessionHeight+blocksPerSesssion*int64(i),
			clientPrivateKey,
			appPrivateKey,
			nodePubKey,
			chain,
		)
		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, err.Codespace(), sdk.CodespaceType(types.ModuleName))
		assert.Equal(t, err.Code(), sdk.CodeType(types.CodeInvalidBlockHeightError))
	}

	// TC2:
	// Test a relay while one session behind and forward,
	// while ClientSessionSyncAllowance = 1
	// --> Success on one session behind
	// --> InvalidBlockHeightError on one session forward
	sessionRangeTc := 1
	types.GlobalPocketConfig.ClientSessionSyncAllowance = int64(sessionRangeTc)

	// First test the minimum boundary
	resp, err = testRelayAt(
		t,
		mockCtx,
		keeper,
		latestSessionHeight-blocksPerSesssion*int64(sessionRangeTc),
		clientPrivateKey,
		appPrivateKey,
		nodePubKey,
		chain,
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp.Response, "bar")

	// Second test the maximum boundary - Error
	resp, err = testRelayAt(
		t,
		mockCtx,
		keeper,
		latestSessionHeight+blocksPerSesssion*int64(sessionRangeTc),
		clientPrivateKey,
		appPrivateKey,
		nodePubKey,
		chain,
	)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, err.Codespace(), sdk.CodespaceType(types.ModuleName))
	assert.Equal(t, err.Code(), sdk.CodeType(types.CodeInvalidBlockHeightError))

	// TC2:
	// Test a relay while two sessions behind and forward,
	// while ClientSessionSyncAllowance = 1
	// --> InvalidBlockHeightError on two sessions behind and forwards
	sessionRangeTc = 2
	types.GlobalPocketConfig.ClientSessionSyncAllowance = 1

	// First test two sessions back - Error
	resp, err = testRelayAt(
		t,
		mockCtx,
		keeper,
		latestSessionHeight-blocksPerSesssion*int64(sessionRangeTc),
		clientPrivateKey,
		appPrivateKey,
		nodePubKey,
		chain,
	)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, err.Codespace(), sdk.CodespaceType(types.ModuleName))
	assert.Equal(t, err.Code(), sdk.CodeType(types.CodeInvalidBlockHeightError))

	// Second test two sessions forward - Error
	resp, err = testRelayAt(
		t,
		mockCtx,
		keeper,
		latestSessionHeight+blocksPerSesssion*int64(sessionRangeTc),
		clientPrivateKey,
		appPrivateKey,
		nodePubKey,
		chain,
	)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, err.Codespace(), sdk.CodespaceType(types.ModuleName))
	assert.Equal(t, err.Code(), sdk.CodeType(types.CodeInvalidBlockHeightError))
}
