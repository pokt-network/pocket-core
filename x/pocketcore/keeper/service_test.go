package keeper

import (
	"encoding/hex"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestKeeper_HandleRelay(t *testing.T) {
	ethereum, err := types.NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	ctx, _, _, _, keeper, keys := createTestInput(t, false)
	mockCtx := new(Ctx)

	ak := keeper.appKeeper.(appsKeeper.Keeper)
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := getRandomPrivateKey()
	apk := appPrivateKey.PublicKey()
	appPubKey := apk.RawString()
	// add app to world state
	app := appsTypes.NewApplication(sdk.Address(apk.Address()), apk, []string{ethereum}, sdk.NewInt(10000000))
	// calculate relays
	app.MaxRelays = ak.CalculateAppRelays(ctx, app)
	// set the vals from the data
	ak.SetApplication(ctx, app)
	ak.SetStakedApplication(ctx, app)
	kp, _ := keeper.Keybase.GetCoinbase()
	npk := kp.PublicKey
	nodePubKey := npk.RawString()
	p := types.Payload{
		Data:    "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}",
		Method:  "",
		Path:    "",
		Headers: nil,
	}
	validRelay := types.Relay{
		Payload: p,
		Proof: types.RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1000,
			RequestHash:        p.HashString(),
			ServicerPubKey:     nodePubKey,
			Blockchain:         ethereum,
			Token: types.AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	appSig, er := appPrivateKey.Sign(validRelay.Proof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	clientSig, er := clientPrivateKey.Sign(validRelay.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Signature = hex.EncodeToString(clientSig)
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://www.google.com").
		Post("/").
		Reply(200).
		BodyString("bar")

	mockCtx.On("KVStore", keeper.storeKey).Return(ctx.KVStore(keeper.storeKey))
	mockCtx.On("KVStore", keys["pos"]).Return(ctx.KVStore(keys["pos"]))
	mockCtx.On("KVStore", keys["params"]).Return(ctx.KVStore(keys["params"]))
	mockCtx.On("KVStore", keys["application"]).Return(ctx.KVStore(keys["application"]))
	mockCtx.On("BlockHeight").Return(ctx.BlockHeight())
	mockCtx.On("MustGetPrevCtx", int64(1000)).Return(ctx)
	mockCtx.On("MustGetPrevCtx", keeper.GetLatestSessionBlockHeight(mockCtx)).Return(ctx)
	mockCtx.On("Logger").Return(ctx.Logger())

	resp, err := keeper.HandleRelay(mockCtx, validRelay)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp.Response, "bar")
}
