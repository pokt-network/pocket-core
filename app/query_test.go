package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/iavl/common"
	tmTypes "github.com/tendermint/tendermint/types"
	"gopkg.in/h2non/gock.v1"
	"math/big"
	"testing"
	"time"
)

func TestQueryBlock(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	height := int64(1)
	select {
	case <-evtChan:
		got, err := nodes.QueryBlock(getInMemoryTMClient(), &height)
		assert.Nil(t, err)
		assert.NotNil(t, got)
	}
	cleanup()
	stopCli()
}

func TestQueryChainHeight(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := nodes.QueryChainHeight(memCli)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), got) // should not be 0 due to empty blocks
	}
	cleanup()
	stopCli()
}

func TestQueryTx(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = nodes.Send(memCodec(), memCli, kb, cb.GetAddress(), kp.GetAddress(), "test", sdk.NewInt(1000))
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case res := <-evtChan:
		time.Sleep(time.Second * 1)
		fmt.Println(res.Data.(tmTypes.EventDataTx))
		got, err := nodes.QueryTransaction(memCli, tx.TxHash)
		assert.Nil(t, err)
		validator, err := nodes.QueryAccountBalance(memCodec(), memCli, kp.GetAddress(), 0)
		assert.Nil(t, err)
		assert.True(t, validator.Equal(sdk.NewInt(1000)))
		assert.NotNil(t, got)
	}
	cleanup()
	stopCli()
}

func TestQueryNodeStatus(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	got, err := nodes.QueryNodeStatus(getInMemoryTMClient())
	assert.Nil(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "pocket-test", got.NodeInfo.Network)
	cleanup()
}

func TestQueryValidators(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := nodes.QueryValidators(memCodec(), memCli, 1)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
	cleanup()
	stopCli()
}

func TestQueryValidator(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	if err != nil {
		t.Fatal(err)
	}
	codec := memCodec()
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := nodes.QueryValidator(codec, memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.Equal(t, cb.GetAddress(), got.Address)
		assert.False(t, got.Jailed)
		assert.True(t, got.StakedTokens.Equal(sdk.NewInt(1000000000000000)))
	}
	cleanup()
	stopCli()
}

func TestQueryDaoBalance(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var err error
		got, err := nodes.QueryDAO(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, big.NewInt(0), got.BigInt())
	}
	cleanup()
	stopCli()
}

func TestQuerySupply(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		gotStaked, gotUnstaked, err := nodes.QuerySupply(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.True(t, gotStaked.Equal(sdk.NewInt(1000000000000000)))
		assert.True(t, gotUnstaked.Equal(sdk.NewInt(2000000000)))
	}
	cleanup()
	stopCli()
}

func TestQueryPOSParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := nodes.QueryPOSParams(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100000), got.MaxValidators)
		assert.Equal(t, int64(1000000), got.StakeMinimum)
		assert.Equal(t, int64(10), got.DAOAllocation)
		assert.Equal(t, sdk.DefaultStakeDenom, got.StakeDenom)
	}
	cleanup()
	stopCli()
}

func TestQueryStakedValidator(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := nodes.QueryStakedValidators(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(got))
	}
	cleanup()
	stopCli()
}

func TestAccountBalance(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var err error
		got, err := nodes.QueryAccountBalance(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, got.Int64(), int64(1000000000))
	}
	cleanup()
	stopCli()
}

func TestQuerySigningInfo(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	cbAddr := cb.GetAddress()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var err error
		got, err := nodes.QuerySigningInfo(memCodec(), memCli, 0, cbAddr)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, got.Address.String(), cbAddr.String())
	}
	cleanup()
	stopCli()
}

func TestQueryPocketSupportedBlockchains(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var err error
		got, err := pocket.QueryPocketSupportedBlockchains(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Contains(t, got, dummyChainsHash)
	}
	cleanup()
	stopCli()
}

func TestQueryPocketParams(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := pocket.QueryParams(memCodec(), memCli, 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, int64(5), got.SessionNodeCount)
		assert.Equal(t, int64(3), got.ProofWaitingPeriod)
		assert.Equal(t, int64(100), got.ClaimExpiration)
		assert.Contains(t, got.SupportedBlockchains, dummyChainsHash)
	}
	cleanup()
	stopCli()
}

func TestQueryAccount(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	acc := getUnstakedAccount(kb)
	assert.NotNil(t, acc)
	select {
	case <-evtChan:
		got, err := nodes.QueryAccount(memCodec(), memCli, acc.GetAddress(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, acc.GetAddress(), got.GetAddress())
	}
	cleanup()
	stopCli()
}

func TestQueryProofs(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		got, err := pocket.QueryProofs(memCodec(), memCli, cb.GetAddress(), 0)
		assert.Nil(t, err)
		assert.Nil(t, got)
	}
	cleanup()
	stopCli()
}

func TestQueryProof(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
		time.Sleep(time.Second / 2)
	}
	select {
	case <-evtChan:
		got, err := pocket.QueryProof(memCodec(), kp.GetAddress(), memCli, dummyChainsHash, kp.PublicKey.RawString(), 1, 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
	}
	cleanup()
	stopCli()
}

func TestQueryStakedApp(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	kp, err := kb.GetCoinbase()
	assert.Nil(t, err)
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	var tx *sdk.TxResponse
	var chains = []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"}

	select {
	case <-evtChan:
		var err error
		memCli, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		tx, err = apps.StakeTx(memCodec(), memCli, kb, chains, sdk.NewInt(1000000), kp, "test")
		assert.Nil(t, err)
		assert.NotNil(t, tx)
	}
	select {
	case <-evtChan:
		got, err := apps.QueryApplication(memCodec(), memCli, kp.GetAddress(), 0)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, sdk.Staked, got.Status)
		assert.Equal(t, false, got.Jailed)
	}
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
	// setup relay
	relay := types.Relay{
		Payload: types.Payload{
			Data: query,
		},
		Proof: types.RelayProof{
			Entropy:            int64(common.RandInt()),
			SessionBlockHeight: sessionBlockheight,
			ServicerPubKey:     nodePublicKey,
			Blockchain:         supportedBlockchain,
			Token:              aat,
			Signature:          "",
		},
	}
	sig, err = appPrivateKey.Sign(relay.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay.Proof.Signature = hex.EncodeToString(sig)
	res, err := json.MarshalIndent(relay, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func TestQueryRelay(t *testing.T) {
	genBz, validators, app := fiveValidatorsOneAppGenesis()
	// setup relay endpoint
	expectedRequest := `"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":64`
	expectedResponse := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
	gock.New(dummyChainsURL).
		Post("").
		BodyString(expectedRequest).
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
	// setup relay
	relay := types.Relay{
		Payload: types.Payload{
			Data: expectedRequest,
		},
		Proof: types.RelayProof{
			Entropy:            32598345349034509,
			SessionBlockHeight: 1,
			ServicerPubKey:     validators[0].PublicKey.RawString(),
			Blockchain:         dummyChainsHash,
			Token:              aat,
			Signature:          "",
		},
	}
	sig, err = appPrivateKey.Sign(relay.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay.Proof.Signature = hex.EncodeToString(sig)
	// setup relay 2
	relay2 := types.Relay{
		Payload: types.Payload{
			Data: expectedRequest,
		},
		Proof: types.RelayProof{
			Entropy:            325983445349034509,
			SessionBlockHeight: 1,
			ServicerPubKey:     validators[0].PublicKey.RawString(),
			Blockchain:         dummyChainsHash,
			Token:              aat,
			Signature:          "",
		},
	}
	sig, err = appPrivateKey.Sign(relay2.Proof.Hash())
	if err != nil {
		panic(err)
	}
	relay2.Proof.Signature = hex.EncodeToString(sig)
	// setup the query
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		res, err := pocket.QueryRelay(memCodec(), memCli, relay)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, res.Response)
		gock.New(dummyChainsURL).
			Post("").
			BodyString(expectedRequest).
			Reply(200).
			BodyString(expectedResponse)
		res, err = pocket.QueryRelay(memCodec(), memCli, relay2)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, res.Response)
		inv, found := types.GetAllInvoices().GetInvoice(types.SessionHeader{
			ApplicationPubKey:  aat.ApplicationPublicKey,
			Chain:              relay.Proof.Blockchain,
			SessionBlockHeight: relay.Proof.SessionBlockHeight,
		})
		assert.True(t, found)
		assert.NotNil(t, inv)
		assert.Equal(t, inv.TotalRelays, int64(2))
		cleanup()
		stopCli()
		gock.Off()
		return
	}
}

func TestQueryDispatch(t *testing.T) {
	genBz, validators, app := fiveValidatorsOneAppGenesis()
	_, kb, cleanup := NewInMemoryTendermintNode(t, genBz)
	appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
	assert.Nil(t, err)
	// Setup Dispatch Request
	key := types.SessionHeader{
		ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
		Chain:              dummyChainsHash,
		SessionBlockHeight: 1,
	}
	// setup the query
	memCli, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		res, err := pocket.QueryDispatch(memCodec(), memCli, key)
		assert.Nil(t, err)
		for _, val := range validators {
			assert.Contains(t, res.SessionNodes, val)
		}
		cleanup()
		stopCli()
	}
}
