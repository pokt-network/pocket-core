package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/x/nodes"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	authTypes "github.com/pokt-network/posmint/x/auth/types"
	"github.com/pokt-network/posmint/x/bank"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
	"gopkg.in/h2non/gock.v1"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRPC_QueryHeight(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		q := newQueryRequest("height", nil)
		rec := httptest.NewRecorder()
		Height(rec, q, httprouter.Params{})
		resp := getResponse(rec)
		assert.Equal(t, "1", resp)
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryBlock(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 1,
		}
		q := newQueryRequest("block", newBody(params))
		rec := httptest.NewRecorder()
		Block(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		var blk core_types.ResultBlock
		err := memCodec().UnmarshalJSON([]byte(resp), &blk)
		assert.Nil(t, err)
		assert.NotEmpty(t, blk.Block.Height)
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryTX(t *testing.T) {
	var tx *types.TxResponse
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	memCLI, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var err error
		_, stopCli, evtChan = subscribeTo(t, tmTypes.EventTx)
		kb := getInMemoryKeybase()
		cb, err := kb.GetCoinbase()
		assert.Nil(t, err)
		tx, err = nodes.Send(memCodec(), memCLI, kb, cb.GetAddress(), cb.GetAddress(), "test", types.NewInt(100))
		assert.Nil(t, err)
	}
	select {
	case <-evtChan:
		var params = hashParams{
			Hash: tx.TxHash,
		}
		q := newQueryRequest("tx", newBody(params))
		rec := httptest.NewRecorder()
		Tx(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		var resTX core_types.ResultTx
		err := json.Unmarshal([]byte(resp), &resTX)
		assert.Nil(t, err)
		assert.NotEmpty(t, resTX.Height)
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryBalance(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		kb := getInMemoryKeybase()
		cb, err := kb.GetCoinbase()
		assert.Nil(t, err)
		var params = heightAddrParams{
			Height:  0,
			Address: cb.GetAddress().String(),
		}
		q := newQueryRequest("balance", newBody(params))
		rec := httptest.NewRecorder()
		Balance(rec, q, httprouter.Params{})
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		var balance types.Int
		err = json.Unmarshal([]byte(resp), &balance)
		assert.Nil(t, err)
		assert.NotZero(t, balance)
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryAccount(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	type Coins struct {
		denom  string `json:string`
		amount int    `json:amount`
	}
	select {
	case <-evtChan:
		kb := getInMemoryKeybase()
		cb, err := kb.GetCoinbase()
		assert.Nil(t, err)
		var params = heightAddrParams{
			Height:  0,
			Address: cb.GetAddress().String(),
		}
		q := newQueryRequest("account", newBody(params))
		rec := httptest.NewRecorder()
		Account(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.Regexp(t, "upokt", string(resp))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryNodes(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		kb := getInMemoryKeybase()
		cb, err := kb.GetCoinbase()
		assert.Nil(t, err)
		var params = heightAndStakingStatusParams{
			Height:        0,
			StakingStatus: "staked",
		}
		q := newQueryRequest("nodes", newBody(params))
		rec := httptest.NewRecorder()
		Nodes(rec, q, httprouter.Params{})
		assert.True(t, strings.Contains(rec.Body.String(), cb.GetAddress().String()))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryNode(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		kb := getInMemoryKeybase()
		cb, err := kb.GetCoinbase()
		assert.Nil(t, err)
		var params = heightAddrParams{
			Height:  0,
			Address: cb.GetAddress().String(),
		}
		q := newQueryRequest("node", newBody(params))
		rec := httptest.NewRecorder()
		Node(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(rec.Body.String(), cb.GetAddress().String()))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryApp(t *testing.T) {
	gBZ, _, app := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, gBZ)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightAddrParams{
			Height:  0,
			Address: app.GetAddress().String(),
		}
		q := newQueryRequest("app", newBody(params))
		rec := httptest.NewRecorder()
		App(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(rec.Body.String(), app.GetAddress().String()))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryApps(t *testing.T) {
	gBZ, _, app := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, gBZ)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightAndStakingStatusParams{
			Height:        0,
			StakingStatus: "staked",
		}
		q := newQueryRequest("apps", newBody(params))
		rec := httptest.NewRecorder()
		Apps(rec, q, httprouter.Params{})
		assert.True(t, strings.Contains(rec.Body.String(), app.GetAddress().String()))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryNodeParams(t *testing.T) {
	gBZ, _, _ := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, gBZ)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("nodeparams", newBody(params))
		rec := httptest.NewRecorder()
		NodeParams(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(rec.Body.String(), "max_validators"))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryAppParams(t *testing.T) {
	gBZ, _, _ := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, gBZ)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("appparams", newBody(params))
		rec := httptest.NewRecorder()
		AppParams(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(rec.Body.String(), "max_applications"))
	}
	cleanup()
	stopCli()
}

func TestRPC_QueryPocketParams(t *testing.T) {
	gBZ, _, _ := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, gBZ)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("pocketparams", newBody(params))
		rec := httptest.NewRecorder()
		PocketParams(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(rec.Body.String(), "chains"))
	}
	cleanup()
	stopCli()
}

func TestRPC_QuerySupportedChains(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("supportedchains", newBody(params))
		rec := httptest.NewRecorder()
		SupportedChains(rec, q, httprouter.Params{})
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		fmt.Println(resp)
		assert.True(t, strings.Contains(resp, dummyChainsHash))
	}
	cleanup()
	stopCli()
}
func TestRPC_QuerySupply(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("supply", newBody(params))
		rec := httptest.NewRecorder()
		Supply(rec, q, httprouter.Params{})

		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)

		var supply querySupplyResponse
		err := json.Unmarshal([]byte(resp), &supply)
		assert.Nil(t, err)
		assert.NotZero(t, supply.Total)
	}
	cleanup()
	stopCli()
}

func TestRPC_StakeNode(t *testing.T) {
	_, _, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		var params = heightParams{
			Height: 0,
		}
		q := newQueryRequest("supply", newBody(params))
		rec := httptest.NewRecorder()
		Supply(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		var supply querySupplyResponse
		err := json.Unmarshal([]byte(resp), &supply)
		assert.Nil(t, err)
		assert.NotZero(t, supply.Total)
	}
	cleanup()
	stopCli()
}

func TestRPC_Relay(t *testing.T) {
	kb := getInMemoryKeybase()
	genBZ, validators, app := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, genBZ)
	// setup relay endpoint
	defer gock.Off()
	expectedRequest := `"jsonrpc":"2.0","method":"web3_sha3","params":["0x68656c6c6f20776f726c64"],"id":64`
	expectedResponse := "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad"
	gock.New(dummyChainsURL).
		Post("").
		BodyString(expectedRequest).
		Reply(200).
		BodyString(expectedResponse)
	appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
	assert.Nil(t, err)
	// setup AAT
	aat := pocketTypes.AAT{
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
	relay := pocketTypes.Relay{
		Payload: pocketTypes.Payload{
			Data: expectedRequest,
		},
		Proof: pocketTypes.RelayProof{
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
	// setup the query
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		q := newClientRequest("relay", newBody(relay))
		rec := httptest.NewRecorder()
		Relay(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		var response pocketTypes.RelayResponse
		err := json.Unmarshal(resp, &response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, response.Response)
		cleanup()
		stopCli()
	}
}

func TestRPC_Dispatch(t *testing.T) {
	kb := getInMemoryKeybase()
	genBZ, validators, app := fiveValidatorsOneAppGenesis()
	_, _, cleanup := NewInMemoryTendermintNode(t, genBZ)
	appPrivateKey, err := kb.ExportPrivateKeyObject(app.Address, "test")
	assert.Nil(t, err)
	// Setup Dispatch Request
	key := pocketTypes.SessionHeader{
		ApplicationPubKey:  appPrivateKey.PublicKey().RawString(),
		Chain:              dummyChainsHash,
		SessionBlockHeight: 1,
	}
	// setup the query
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	select {
	case <-evtChan:
		q := newClientRequest("dispatch", newBody(key))
		rec := httptest.NewRecorder()
		Dispatch(rec, q, httprouter.Params{})
		resp := getJSONResponse(rec)
		rawResp := string(resp)
		assert.Regexp(t, key.ApplicationPubKey, rawResp)
		assert.Regexp(t, key.Chain, rawResp)

		for _, validator := range validators {
			assert.Regexp(t, validator.Address.String(), rawResp)
		}
		cleanup()
		stopCli()
	}
}

func TestRPC_RawTX(t *testing.T) {
	_, kb, cleanup := NewInMemoryTendermintNode(t, oneValTwoNodeGenesisState())
	cb, err := kb.GetCoinbase()
	assert.Nil(t, err)
	kp, err := kb.Create("test")
	assert.Nil(t, err)
	pk, err := kb.ExportPrivateKeyObject(cb.GetAddress(), "test")
	assert.Nil(t, err)
	_, stopCli, evtChan := subscribeTo(t, tmTypes.EventNewBlock)
	// create the transaction
	txBz, err := auth.DefaultTxEncoder(memCodec())(authTypes.NewTestTx(types.Context{}.WithChainID("pocket-test"),
		[]types.Msg{bank.MsgSend{
			FromAddress: cb.GetAddress(),
			ToAddress:   kp.GetAddress(),
			Amount:      types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(1))),
		}},
		[]crypto.PrivateKey{pk},
		common.RandInt64(),
		types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(100000)))))
	assert.Nil(t, err)
	select {
	case <-evtChan:
		var err error
		params := sendRawTxParams{
			Addr:        cb.GetAddress().String(),
			RawHexBytes: hex.EncodeToString(txBz),
		}
		q := newClientRequest("rawtx", newBody(params))
		rec := httptest.NewRecorder()
		SendRawTx(rec, q, httprouter.Params{})
		resp := getResponse(rec)
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		var response types.TxResponse
		err = memCodec().UnmarshalJSON([]byte(resp), &response)
		assert.Nil(t, err)
		assert.True(t, strings.Contains(response.Logs.String(), `"success":true`))
	}
	cleanup()
	stopCli()
}

func newBody(params interface{}) io.Reader {
	bz, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(bz)
	return reader
}

func newClientRequest(query string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", "localhost:8081/v1/client/"+query, body)
	if err != nil {
		panic("could not create request: %v")
	}
	return req
}

func newQueryRequest(query string, body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", "localhost:8081/v1/query/"+query, body)
	if err != nil {
		panic("could not create request: %v")
	}
	return req
}

func getResponse(rec *httptest.ResponseRecorder) string {
	res := rec.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("could not read response: " + err.Error())
	}
	if strings.Contains(string(b), "error") {
		return string(b)
	}
	resp, err := strconv.Unquote(string(b))
	if err != nil {
		panic("could not unquote resp: " + err.Error())
	}
	return resp
}

func getJSONResponse(rec *httptest.ResponseRecorder) []byte {
	res := rec.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("could not read response: " + err.Error())
	}
	return b
}
