package keeper

import (
	"encoding/hex"
	appsKeeper "github.com/pokt-network/pocket-core/x/apps/keeper"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
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
	ctx, _, _, _, keeper := createTestInput(t, false)
	ak := keeper.appKeeper.(appsKeeper.Keeper)
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := crypto.PublicKey(clientPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	appPrivateKey := getRandomPrivateKey()
	apk := appPrivateKey.PubKey().(ed25519.PubKeyEd25519)
	appPubKey := crypto.PublicKey(apk).String()
	// add app to world state
	app := appsTypes.NewApplication(sdk.ValAddress(apk.Address()), apk, []string{ethereum}, sdk.NewInt(10000000))
	// calculate relays
	app.MaxRelays = ak.CalculateAppRelays(ctx, app)
	// set the vals from the data
	ak.SetApplication(ctx, app)
	ak.SetAppByConsAddr(ctx, app)
	ak.SetStakedApplication(ctx, app)
	kp, _ := keeper.Keybase.GetCoinbase()
	npk := kp.PubKey
	nodePubKey := hex.EncodeToString(crypto.PublicKey(npk.(ed25519.PubKeyEd25519)).Bytes())
	validRelay := types.Relay{
		Payload: types.Payload{
			Data:    "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}",
			Method:  "",
			Path:    "",
			Headers: nil,
		},
		Proof: types.RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
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
	resp, err := keeper.HandleRelay(ctx, validRelay)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp)
	assert.Equal(t, resp.Response, "bar")
}
