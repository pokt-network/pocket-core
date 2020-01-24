package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	tmTypes "github.com/tendermint/tendermint/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRPC_Height(t *testing.T) {
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

func TestRPC_Block(t *testing.T) {
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
	if strings.Contains(string(b), "error") {
		return string(b)
	}
	resp, err := strconv.Unquote(string(b))
	if err != nil {
		panic("could not unquote resp: " + err.Error())
	}
	return resp
}
