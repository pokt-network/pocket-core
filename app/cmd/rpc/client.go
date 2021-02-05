package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// Dispatch supports CORS functionality
func Dispatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if cors(&w, r) {
		return
	}
	d := types.SessionHeader{}
	if err := PopModel(w, r, ps, &d); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.HandleDispatch(d)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

type RPCRelayResponse struct {
	Signature string `json:"signature"`
	Response  string `json:"response"`
	// remove proof object because client already knows about it
}

// Relay supports CORS functionality
func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = types.Relay{}
	if cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &relay); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.HandleRelay(relay)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	response := RPCRelayResponse{
		Signature: res.Signature,
		Response:  res.Response,
	}
	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// Stop
func Stop(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	value := r.URL.Query().Get("authtoken")
	if value == app.AuthToken.Value {
		app.ShutdownPocketCore()
		err := app.PCA.TMNode().Stop()
		if err != nil {
			fmt.Println(err)
			WriteErrorResponse(w, 400, err.Error())
			fmt.Println("Force Stop , PID:" + fmt.Sprint(os.Getpid()))
			os.Exit(1)
		}
		fmt.Println("Stop Successful, PID:" + fmt.Sprint(os.Getpid()))
		os.Exit(0)
	} else {
		WriteErrorResponse(w, 401, "wrong authtoken "+value)
	}
}

// Challenge supports CORS functionality
func Challenge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var challenge = types.ChallengeProofInvalidData{}
	if cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &challenge); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.HandleChallenge(challenge)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

type SendRawTxParams struct {
	Addr        string `json:"address"`
	RawHexBytes string `json:"raw_hex_bytes"`
}

func SendRawTx(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = SendRawTxParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	bz, err := hex.DecodeString(params.RawHexBytes)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	res, err := app.PCA.SendRawTx(params.Addr, bz)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	j, er := app.Codec().MarshalJSON(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

type simRelayParams struct {
	Url     string        `json:"chain_url"`
	Payload types.Payload `json:"payload"` // the data payload of the request
}

func SimRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = simRelayParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	// retrieve the hosted blockchain url requested
	url := strings.Trim(params.Url, `/`)
	if len(params.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(params.Payload.Path, `/`)
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(params.Payload.Data, url, params.Payload.Method, params.Payload.Headers)

	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(res), r.URL.Path, r.Host)
}

func executeHTTPRequest(payload string, url string, method string, headers map[string]string) (string, error) {
	// generate an http request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	// add headers if needed
	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	// execute the request
	resp, err := (&http.Client{Timeout: types.GetRPCTimeout() * time.Millisecond}).Do(req)
	if err != nil {
		return payload, err
	}
	// ensure code is 200
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Expected Code 200 from Request got %v", resp.StatusCode)
	}
	// read all bz
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// return
	return string(body), nil
}
