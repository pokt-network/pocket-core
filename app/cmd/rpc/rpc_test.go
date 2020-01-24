package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/x/nodes"
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"
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
		resp := getResponse(rec)
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
		resp := getResponse(rec)
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, cb.GetAddress().String()))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, cb.GetAddress().String()))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, app.GetAddress().String()))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, app.GetAddress().String()))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, "max_validators"))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, "max_applications"))
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
		resp := getResponse(rec)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp)
		assert.True(t, strings.Contains(resp, "chains"))
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
		resp := getResponse(rec)
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

func newBody(params interface{}) io.Reader {
	bz, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(bz)
	return reader
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
	if strings.Contains(string(b), "error"){
		return string(b)
	}
	resp, err := strconv.Unquote(string(b))
	if err != nil {
		panic("could not unquote resp: " + err.Error())
	}
	return resp
}
