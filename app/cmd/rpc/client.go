package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
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

type RPCRelayErrorResponse struct {
	Error    error                   `json:"error"`
	Dispatch *types.DispatchResponse `json:"dispatch"`
}

// Relay supports CORS functionality
func Relay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = types.Relay{}
	if cors(&w, r) {
		return
	}
	if err := PopModel(w, r, ps, &relay); err != nil {
		response := RPCRelayErrorResponse{
			Error: err,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}
	res, dispatch, err := app.PCA.HandleRelay(relay)
	if err != nil {
		response := RPCRelayErrorResponse{
			Error:    err,
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
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

// UpdateChains
func UpdateChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	value := r.URL.Query().Get("authtoken")
	if value == app.AuthToken.Value {
		var hostedChainsSlice []types.HostedBlockchain
		if err := PopModel(w, r, ps, &hostedChainsSlice); err != nil {
			WriteErrorResponse(w, 400, err.Error())
			return
		}
		m := make(map[string]types.HostedBlockchain)
		for _, chain := range hostedChainsSlice {
			if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
				WriteErrorResponse(w, 400, fmt.Sprintf("invalid ID: %s in network identifier in json", chain.ID))
				return
			}
			m[chain.ID] = chain
		}
		result, err := app.PCA.SetHostedChains(m)
		if err != nil {
			WriteErrorResponse(w, 400, err.Error())
		} else {
			j, er := json.Marshal(result)
			if er != nil {
				WriteErrorResponse(w, 400, er.Error())
				return
			}
			WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
		}
	} else {
		WriteErrorResponse(w, 401, "wrong authtoken "+value)
	}
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
	RelayNetworkID string        `json:"relay_network_id"` // RelayNetworkID
	Payload        types.Payload `json:"payload"`          // the data payload of the request
}

func SimRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = simRelayParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	hostedChains := app.NewHostedChains(false)

	chain, err := hostedChains.GetChain(params.RelayNetworkID)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}
	url := strings.Trim(chain.URL, `/`)
	if len(params.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(params.Payload.Path, `/`)
	}
	// do basic http request on the relay
	res, er := executeHTTPRequest(params.Payload.Data, url, types.GlobalPocketConfig.UserAgent, chain.BasicAuth, params.Payload.Method, params.Payload.Headers)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, string(res), r.URL.Path, r.Host)
}

func executeHTTPRequest(payload, url, userAgent string, basicAuth types.BasicAuth, method string, headers map[string]string) (string, error) {
	// generate an http request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if basicAuth.Username != "" {
		req.SetBasicAuth(basicAuth.Username, basicAuth.Password)
	}
	if userAgent == "" {
		req.Header.Set("User-Agent", userAgent)
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
