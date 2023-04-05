package mesh

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	types4 "github.com/pokt-network/pocket-core/app/cmd/rpc/types"
	sdk "github.com/pokt-network/pocket-core/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// sdkErrorResponse - response error format for re-implemented endpoints.
type sdkErrorResponse struct {
	Code      sdk.CodeType      `json:"code"`
	Codespace sdk.CodespaceType `json:"codespace"`
	Error     string            `json:"message"`
}

// meshHealthResponse - response payload of /v1/mesh/health
type meshHealthResponse struct {
	Version   string `json:"version"`
	Servicers int    `json:"servicers"`
	FullNodes int    `json:"full_nodes"`
}

// meshRPCRelayResult response payload of /v1/client/relay
type meshRPCRelayResult struct {
	Success  bool                          `json:"signature"`
	Error    error                         `json:"error"`
	Dispatch *pocketTypes.DispatchResponse `json:"dispatch"`
}

// meshRPCSessionResult - response payload of /v1/private/mesh/session
type meshRPCSessionResult struct {
	Success         bool              `json:"success"`
	Error           *sdkErrorResponse `json:"error"`
	Dispatch        *dispatchResponse `json:"dispatch"`
	RemainingRelays json.Number       `json:"remaining_relays"`
}

// meshRPCRelayResponse - response payload of /v1/private/mesh/relay
type meshRPCRelayResponse struct {
	Success  bool              `json:"signature"`
	Error    *sdkErrorResponse `json:"error"`
	Dispatch *dispatchResponse `json:"dispatch"`
}

// meshCheckPayload - payload used to call /v1/private/mesh/check
type meshCheckPayload struct {
	Servicers []string `json:"servicers"`
	Chains    []string `json:"chains"`
}

// meshCheckResponse - response payload of /v1/private/mesh/check
type meshCheckResponse struct {
	Success        bool               `json:"success"`
	Error          *sdkErrorResponse  `json:"error"`
	Status         app.HealthResponse `json:"status"`
	Servicers      bool               `json:"servicers"`
	Chains         bool               `json:"chains"`
	WrongServicers []string           `json:"wrong_servicers"`
	WrongChains    []string           `json:"wrong_chains"`
}

// ++++++++++++++++++++ MESH CLIENT - PUBLIC ROUTES ++++++++++++++++++++

// meshNodeRelay - handle mesh node relay request, call handleRelay
// path: /v1/client/relay
func meshNodeRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if cors(&w, r) {
		return
	}
	var relay = pocketTypes.Relay{}

	if err := rpc.PopModel(w, r, ps, &relay); err != nil {
		response := meshRPCRelayResponse{
			Success: false,
			Error:   newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	logger.Debug(fmt.Sprintf("handling relay %s", relay.RequestHashString()))
	res, dispatch, err := handleRelay(&relay)

	if err != nil {
		response := meshRPCRelayResponse{
			Success:  false,
			Error:    newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := rpc.RPCRelayResponse{
		Signature: res.Signature,
		Response:  res.Response,
	}

	j, er := json.Marshal(response)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
	logger.Debug(fmt.Sprintf("relay %s done", relay.RequestHashString()))
}

// meshSimulateRelay - handle a simulated relay to test connectivity to the chains that this should be serving.
// this will only be enabled if start node with --simulateRelays
// path: /v1/client/sim
func meshSimulateRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = rpc.SimRelayParams{}
	if err := rpc.PopModel(w, r, ps, &params); err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	chain, err := chains.GetChain(params.RelayNetworkID)
	if err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	_url := strings.Trim(chain.URL, `/`)
	if len(params.Payload.Path) > 0 {
		_url = _url + "/" + strings.Trim(params.Payload.Path, `/`)
	}

	logger.Debug(
		fmt.Sprintf(
			"executing simulated relay of chain %s",
			chain.ID,
		),
	)
	// do basic http request on the relay
	res, er := executeBlockchainHTTPRequest(
		params.Payload.Data, _url, app.GlobalMeshConfig.UserAgent,
		chain.BasicAuth, params.Payload.Method, params.Payload.Headers,
	)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}
	rpc.WriteResponse(w, res, r.URL.Path, r.Host)
}

// ++++++++++++++++++++ MESH CLIENT - PRIVATE ROUTES ++++++++++++++++++++

// meshHealth - handle mesh health request
// path: /v1/private/mesh/health
func meshHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	res := meshHealthResponse{
		Version:   AppVersion,
		Servicers: servicerMap.Size(),
		FullNodes: nodesMap.Size(),
	}
	j, er := json.Marshal(res)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshChains - return load chains from app.GlobalMeshConfig.ChainsName file
// path: /v1/private/mesh/chains
func meshChains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	c := make([]pocketTypes.HostedBlockchain, 0)

	for _, chain := range chains.M {
		c = append(c, chain)
	}

	j, err := json.Marshal(c)
	if err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	rpc.WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNode - return servicer node configured by servicer_priv_key.json - return address
// path: /v1/private/mesh/servicer
func meshServicerNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	servicers := make([]types4.PublicPocketNode, 0)

	mutex.Lock()
	for _, a := range servicerList {
		servicers = append(servicers, types4.PublicPocketNode{
			Address: a,
		})
	}
	mutex.Unlock()

	j, err := json.Marshal(servicers)
	if err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	rpc.WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshUpdateChains - update chains in memory and also chains.json file.
// path: /v1/private/mesh/updatechains
func meshUpdateChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	var hostedChainsSlice []pocketTypes.HostedBlockchain
	if err := rpc.PopModel(w, r, ps, &hostedChainsSlice); err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}
	m := make(map[string]pocketTypes.HostedBlockchain)
	for _, chain := range hostedChainsSlice {
		if err := nodesTypes.ValidateNetworkIdentifier(chain.ID); err != nil {
			rpc.WriteErrorResponse(w, 400, fmt.Sprintf("invalid ID: %s in network identifier in json", chain.ID))
			return
		}
	}
	chains = &pocketTypes.HostedBlockchains{
		M: m,
		L: sync.RWMutex{},
	}

	j, er := json.Marshal(chains.M)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)

	updateChains(hostedChainsSlice)
}

// meshStop - gracefully stop mesh rpc server. Also, this should stop new relays and wait/flush all pending relays, otherwise they will get loose.
// path: /v1/private/mesh/stop
func meshStop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}
	StopMeshRPC()
	fmt.Println("Stop Successful, PID:" + fmt.Sprint(os.Getpid()))
	os.Exit(0)
}

// ++++++++++++++++++++ POKT CLIENT - PRIVATE ROUTES ++++++++++++++++++++

// meshServicerNodeRelay - receive relays that was processed by mesh node
// path: /v1/private/mesh/relay
func meshServicerNodeRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var relay = pocketTypes.Relay{}

	if cors(&w, r) {
		return
	}

	token := r.URL.Query().Get("authtoken")
	if token != app.AuthToken.Value {
		rpc.WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := meshRPCRelayResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := rpc.PopModel(w, r, ps, &relay); err != nil {
		response := rpc.RPCRelayErrorResponse{
			Error: err,
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	_, dispatch, err := app.PCA.HandleRelay(relay, true)
	if err != nil {
		response := meshRPCRelayResult{
			Success:  false,
			Error:    err,
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := meshRPCRelayResult{
		Success:  true,
		Dispatch: dispatch,
	}
	j, er := json.Marshal(response)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNodeSession - receive requests from mesh node to validate a session for an app/servicer/blockchain on the servicer node data
// path: /v1/private/mesh/session
func meshServicerNodeSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var session pocketTypes.MeshSession

	token := r.URL.Query().Get("authtoken")
	if token != app.AuthToken.Value {
		rpc.WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := meshRPCSessionResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := rpc.PopModel(w, r, ps, &session); err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	res, err := app.PCA.HandleMeshSession(session)

	if err != nil {
		response := meshRPCSessionResult{
			Success: false,
			Error:   newSdkErrorFromPocketSdkError(err),
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	dispatch := dispatchResponse{
		BlockHeight: res.Session.BlockHeight,
		Session: dispatchSession{
			Header: res.Session.Session.SessionHeader,
			Key:    hex.EncodeToString(res.Session.Session.SessionKey),
			Nodes:  make([]dispatchSessionNode, 0),
		},
	}

	for i := range res.Session.Session.SessionNodes {
		sNode, ok := res.Session.Session.SessionNodes[i].(nodesTypes.Validator)
		if !ok {
			continue
		}
		dispatch.Session.Nodes = append(dispatch.Session.Nodes, dispatchSessionNode{
			Address:       sNode.Address.String(),
			Chains:        sNode.Chains,
			Jailed:        sNode.Jailed,
			OutputAddress: sNode.OutputAddress.String(),
			PublicKey:     sNode.PublicKey.String(),
			ServiceUrl:    sNode.ServiceURL,
			Status:        sNode.Status,
			Tokens:        sNode.GetTokens().String(),
			UnstakingTime: sNode.UnstakingCompletionTime,
		})
	}

	response := meshRPCSessionResult{
		Success:         true,
		Dispatch:        &dispatch,
		RemainingRelays: json.Number(strconv.FormatInt(res.RemainingRelays, 10)),
	}
	j, er := json.Marshal(response)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNodeCheck - receive requests from mesh node to validate servicers, chains and health status
// path: /v1/private/mesh/check
func meshServicerNodeCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var checkPayload meshCheckPayload

	token := r.URL.Query().Get("authtoken")
	if token != app.AuthToken.Value {
		rpc.WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	// useful just to test that mesh node is able to reach servicer - this payload should be ignored
	if verify == "true" {
		code := 200
		j, _ := json.Marshal(map[string]interface{}{})
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := rpc.PopModel(w, r, ps, &checkPayload); err != nil {
		rpc.WriteErrorResponse(w, 400, err.Error())
		return
	}

	chainsMap, err := app.PCA.QueryHostedChains()
	health := app.PCA.QueryHealth(rpc.APIVersion)

	if err != nil {
		response := meshRPCSessionResult{
			Success: false,
			Error:   newSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
		}
		j, _ := json.Marshal(response)
		rpc.WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := meshCheckResponse{
		Success:        true,
		Status:         health,
		Servicers:      true,
		Chains:         true,
		WrongServicers: make([]string, 0),
		WrongChains:    make([]string, 0),
	}

	for _, address := range checkPayload.Servicers {
		if err := servicerIsSupported(address); err != nil {
			response.Servicers = false
			response.WrongServicers = append(response.WrongServicers, address)
		}
	}

	for _, chain := range checkPayload.Chains {
		if _, ok := chainsMap[chain]; !ok {
			response.Chains = false
			response.WrongChains = append(response.WrongChains, chain)
		}
	}

	j, er := json.Marshal(response)
	if er != nil {
		rpc.WriteErrorResponse(w, 400, er.Error())
		return
	}

	rpc.WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}
