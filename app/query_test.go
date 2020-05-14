package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	apps "github.com/pokt-network/pocket-core/x/apps"
	types3 "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	types2 "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/iavl/common"
	tmTypes "github.com/tendermint/tendermint/types"
	"gopkg.in/h2non/gock.v1"
)

func TestQueryBlock(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	height := int64(1)
	<-evtChan // Wait for block
	got, err := PCA.QueryBlock(&height)
	assert.Nil(t, err)
	assert.NotNil(t, got)

	cleanup()
	stopCli()
}

func TestQueryChainHeight(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryHeight()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), got) // should not be 0 due to empty blocks

	cleanup()
	stopCli()
}

func TestQueryTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	<-evtChan // Wait for block
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000))
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-evtChan // Wait for tx
	got, err := PCA.QueryTx(tx.TxHash)
	assert.Nil(t, err)
	validator, err := PCA.QueryBalance(kp.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.True(t, validator.Equal(sdk.NewInt(1000)))
	assert.NotNil(t, got)

	cleanup()
	stopCli()
}

func TestQueryValidators(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, twoValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryNodes(1, types2.QueryValidatorsParams{Page: 1, Limit: 1})
	assert.Nil(t, err)
	res := got.Result.([]types2.Validator)
	assert.Equal(t, 1, len(res))
	got, err = PCA.QueryNodes(1, types2.QueryValidatorsParams{Page: 2, Limit: 1})
	assert.Nil(t, err)
	res = got.Result.([]types2.Validator)
	assert.Equal(t, 1, len(res))
	got, err = PCA.QueryNodes(1, types2.QueryValidatorsParams{Page: 1, Limit: 1000})
	assert.Nil(t, err)
	res = got.Result.([]types2.Validator)
	assert.Equal(t, 2, len(res))

	cleanup()
	stopCli()
}
func TestQueryApps(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"00"}

	<-evtChan // Wait for block
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-evtChan // Wait for tx
	got, err := PCA.QueryApps(0, types3.QueryApplicationsWithOpts{
		Page:  1,
		Limit: 1,
	})
	assert.Nil(t, err)
	slice, ok := takeArg(got.Result, reflect.Slice)
	if !ok {
		t.Fatalf("couldn't convert arg to slice")
	}
	assert.Equal(t, 1, slice.Len())
	got, err = PCA.QueryApps(0, types3.QueryApplicationsWithOpts{
		Page:  2,
		Limit: 1,
	})
	assert.Nil(t, err)
	slice, ok = takeArg(got.Result, reflect.Slice)
	if !ok {
		t.Fatalf("couldn't convert arg to slice")
	}
	assert.Equal(t, 1, slice.Len())
	got, err = PCA.QueryApps(0, types3.QueryApplicationsWithOpts{
		Page:  1,
		Limit: 2,
	})
	assert.Nil(t, err)
	slice, ok = takeArg(got.Result, reflect.Slice)
	if !ok {
		t.Fatalf("couldn't convert arg to slice")
	}
	assert.Equal(t, 2, slice.Len())

	stopCli()
	cleanup()
}

func TestQueryValidator(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	if err != nil {
		t.Fatal(err)
	}
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryNode(cb.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.Equal(t, cb.GetAddress(), got.Address)
	assert.False(t, got.Jailed)
	assert.True(t, got.StakedTokens.Equal(sdk.NewInt(1000000000000000)))

	cleanup()
	stopCli()
}

func TestQueryDaoBalance(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := gov.QueryDAO(memCodec(), memCli, 0)
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(1000), got.BigInt())

	cleanup()
	stopCli()
}

func TestQueryACL(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := gov.QueryACL(memCodec(), memCli, 0)
	assert.Nil(t, err)
	assert.Equal(t, got, testACL)

	cleanup()
	stopCli()
}

func TestQueryDaoOwner(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kb := getInMemoryKeybase()
	cb, err := kb.GetCoinbase()
	if err != nil {
		t.Fatal(err)
	}
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := gov.QueryDAOOwner(memCodec(), memCli, 0)
	assert.Nil(t, err)
	assert.Equal(t, got.String(), cb.GetAddress().String())

	cleanup()
	stopCli()
}

func TestQueryUpgrade(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	var err error
	got, err := gov.QueryUpgrade(memCodec(), memCli, 0)
	assert.Nil(t, err)
	assert.Equal(t, got.UpgradeHeight(), int64(10000))

	cleanup()
	stopCli()
}

func TestQuerySupply(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	gotStaked, total, err := PCA.QueryTotalNodeCoins(0)
	fmt.Println(err)
	assert.Nil(t, err)
	fmt.Println(gotStaked, total)
	assert.True(t, gotStaked.Equal(sdk.NewInt(1000000000000000)))
	assert.True(t, total.Equal(sdk.NewInt(1000002010001000)))

	cleanup()
	stopCli()
}

func TestQueryPOSParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryNodeParams(0)
	assert.Nil(t, err)
	assert.Equal(t, uint64(100000), got.MaxValidators)
	assert.Equal(t, int64(1000000), got.StakeMinimum)
	assert.Equal(t, int64(10), got.DAOAllocation)
	assert.Equal(t, sdk.DefaultStakeDenom, got.StakeDenom)

	cleanup()
	stopCli()
}

func TestAccountBalance(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryBalance(cb.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, got.Int64(), int64(1000000000))

	cleanup()
	stopCli()
}

func TestQuerySigningInfo(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	cbAddr := cb.GetAddress()
	assert.Nil(t, err)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QuerySigningInfo(0, cbAddr.String())
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, got.Address.String(), cbAddr.String())

	cleanup()
	stopCli()
}

func TestQueryPocketSupportedBlockchains(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	var err error
	got, err := PCA.QueryPocketSupportedBlockchains(0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Contains(t, got, PlaceholderHash)

	cleanup()
	stopCli()
}

func TestQueryPocketParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryPocketParams(0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, int64(5), got.SessionNodeCount)
	assert.Equal(t, int64(3), got.ClaimSubmissionWindow)
	assert.Equal(t, int64(100), got.ClaimExpiration)
	assert.Contains(t, got.SupportedBlockchains, PlaceholderHash)

	cleanup()
	stopCli()
}

func TestQueryAccount(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	acc := getUnstakedAccount(kb)
	assert.NotNil(t, acc)
	<-evtChan // Wait for block
	got, err := PCA.QueryAccount(acc.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, acc.GetAddress(), (*got).GetAddress())

	cleanup()
	stopCli()
}

func TestQueryProofs(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	got, err := PCA.QueryReceipts(cb.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.Nil(t, got)

	cleanup()
	stopCli()
}

func TestQueryProof(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"00"}

	<-evtChan // Wait for block
	tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-evtChan // Wait for tx
	_, err = PCA.QueryReceipt(PlaceholderHash, kp.PublicKey.RawString(), kp.GetAddress().String(), "relay", 1, 0)
	assert.Nil(t, err)

	cleanup()
	stopCli()
}

func TestQueryStakedpp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"00"}
	<-evtChan // Wait for block
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventTx)
	tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	<-evtChan // Wait for  tx
	got, err := PCA.QueryApp(kp.GetAddress().String(), 0)
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, sdk.Staked, got.Status)
	assert.Equal(t, false, got.Jailed)

	cleanup()
	stopCli()
}

func TestRelayGenerator(t *testing.T) {
	const appPrivKey = "e63df045400c136dae909c6bfabfe632dd37e44abbabea3e9fb1f672bd21c04567f0446d45f3e1ba9f3edc957018174cb82871521ca793acdb45898ec4b1c479"
	const nodePublicKey = "7eb410b363df8f71caf6d3a88f11360b74abbcf7e1293cfbc88a021d966110d5"
	const sessionBlockheight = 1
	const query = `{"jsonrpc":"2.0","method":"net_version","params":[],"id":67}`
	const supportedBlockchain = "49aff8a9f51b268f6fc485ec14fb08466c3ec68c8d86d9b5810ad80546b65f29"
	apkBz, err := hex.DecodeString(appPrivKey)
	if err != nil {
		panic(err)
	}
	var appPrivateKey crypto.Ed25519PrivateKey
	copy(appPrivateKey[:], apkBz)
	aat := types.AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivateKey.PublicKey().RawString(),
		ClientPublicKey:      appPrivateKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	sig, err := appPrivateKey.Sign(aat.Hash())
	if err != nil {
		panic(err)
	}
	aat.ApplicationSignature = hex.EncodeToString(sig)
	payload := types.Payload{
		Data: query,
	}
	// setup relay
	relay := types.Relay{
		Payload: payload,
		Proof: types.RelayProof{
			Entropy:            int64(common.RandInt()),
			SessionBlockHeight: sessionBlockheight,
			ServicerPubKey:     nodePublicKey,
			Blockchain:         supportedBlockchain,
			Token:              aat,
			Signature:          "",
		},
	}
	relay.Proof.RequestHash = relay.RequestHashString()
	sig, err = appPrivateKey.Sign(relay.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay.Proof.Signature = hex.EncodeToString(sig)
	_, err = json.MarshalIndent(relay, "", "  ")
	if err != nil {
		panic(err)
	}
}

func TestQueryRelay(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	const headerKey = "foo"
	const headerVal = "bar"
	genBz, _, validators, app := fiveValidatorsOneAppGenesis()
	// setup relay endpoint
	expectedRequest := `"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":64`
	expectedResponse := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
	gock.New(PlaceholderURL).
		Post("").
		BodyString(expectedRequest).
		MatchHeader(headerKey, headerVal).
		Reply(200).
		BodyString(expectedResponse)
	_, kb, cleanup := NewInMemoryTendermintNode(t, genBz)
	appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
	assert.Nil(t, err)
	// setup AAT
	aat := types.AAT{
		Version:              "0.0.1",
		ApplicationPublicKey: appPrivateKey.PublicKey().RawString(),
		ClientPublicKey:      appPrivateKey.PublicKey().RawString(),
		ApplicationSignature: "",
	}
	sig, err := appPrivateKey.Sign(aat.Hash())
	if err != nil {
		panic(err)
	}
	aat.ApplicationSignature = hex.EncodeToString(sig)
	payload := types.Payload{
		Data:    expectedRequest,
		Headers: map[string]string{headerKey: headerVal},
	}
	// setup relay
	relay := types.Relay{
		Payload: payload,
		Meta:    types.RelayMeta{BlockHeight: 5}, // todo race condition here
		Proof: types.RelayProof{
			Entropy:            32598345349034509,
			SessionBlockHeight: 1,
			ServicerPubKey:     validators[0].PublicKey.RawString(),
			Blockchain:         PlaceholderHash,
			Token:              aat,
			Signature:          "",
		},
	}
	relay.Proof.RequestHash = relay.RequestHashString()
	sig, err = appPrivateKey.Sign(relay.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay.Proof.Signature = hex.EncodeToString(sig)
	_, _, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	res, err := PCA.QueryRelay(relay)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedResponse, res.Response)
	gock.New(PlaceholderURL).
		Post("").
		BodyString(expectedRequest).
		Reply(200).
		BodyString(expectedResponse)

	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
<<<<<<< HEAD
	<-evtChan // Wait for tx
	inv, found := types.GetEvidence(types.SessionHeader{
		ApplicationPubKey:  aat.ApplicationPublicKey,
		Chain:              relay.Proof.Blockchain,
		SessionBlockHeight: relay.Proof.SessionBlockHeight,
	}, types.RelayEvidence)
	assert.True(t, found)
	assert.NotNil(t, inv)
	assert.Equal(t, inv.NumOfProofs, int64(1))
	cleanup()
	stopCli()
	gock.Off()
=======
	select {
	case <-evtChan:
		res, err := PCA.QueryRelay(relay)
		assert.Nil(t, err, err)
		assert.Equal(t, expectedResponse, res.Response)
		gock.New(PlaceholderURL).
			Post("").
			BodyString(expectedRequest).
			Reply(200).
			BodyString(expectedResponse)
		_, stopCli, evtChan = subscribeTo(t, tmTypes.EventNewBlock)
		select {
		case <-evtChan:
			inv, err := types.GetEvidence(types.SessionHeader{
				ApplicationPubKey:  aat.ApplicationPublicKey,
				Chain:              relay.Proof.Blockchain,
				SessionBlockHeight: relay.Proof.SessionBlockHeight,
			}, types.RelayEvidence, 10000)
			assert.Nil(t, err)
			assert.NotNil(t, inv)
			assert.Equal(t, inv.NumOfProofs, int64(1))
			cleanup()
			stopCli()
			gock.Off()
			return
		}
	}
>>>>>>> #873
}

func TestQueryDispatch(t *testing.T) {
	genBz, _, validators, app := fiveValidatorsOneAppGenesis()
	_, kb, cleanup := NewInMemoryTendermintNode(t, genBz)
	appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
	assert.Nil(t, err)
	// Setup HandleDispatch Request
	key := types.SessionHeader{
		ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
		Chain:              PlaceholderHash,
		SessionBlockHeight: 1,
	}
	// setup the query
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	<-evtChan // Wait for block
	res, err := PCA.QueryDispatch(key)
	assert.Nil(t, err)
	for _, val := range validators {
		assert.Contains(t, res.Session.SessionNodes, val)
	}
	cleanup()
	stopCli()
}
