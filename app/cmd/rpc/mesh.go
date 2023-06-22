package rpc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc/mesh"
	types4 "github.com/pokt-network/pocket-core/app/cmd/rpc/types"
	sdk "github.com/pokt-network/pocket-core/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"io"
	"net/http"
	"os"
	"strconv"
)

// ++++++++++++++++++++ MESH CLIENT - PUBLIC ROUTES ++++++++++++++++++++

// meshNodeRelay - handle mesh node relay request, call handleRelay
// path: /v1/client/relay
func meshNodeRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if cors(&w, r) {
		return
	}
	var relay = pocketTypes.Relay{}

	if err := PopModel(w, r, ps, &relay); err != nil {
		response := mesh.RPCRelayResponse{
			Success: false,
			Error:   mesh.NewSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	mesh.GetLogger().Debug(fmt.Sprintf("handling relay %s", relay.RequestHashString()))
	res, dispatch, err := mesh.HandleRelay(&relay)

	if err != nil {
		response := mesh.RPCRelayResponse{
			Success:  false,
			Error:    mesh.NewSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
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
	mesh.GetLogger().Debug(fmt.Sprintf("relay %s done", relay.RequestHashString()))
}

// meshSimulateRelay - handle a simulated relay to test connectivity to the chains that this should be serving.
// this will only be enabled if start node with --simulateRelays
// path: /v1/client/sim
func meshSimulateRelay(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var params = SimRelayParams{}
	if err := PopModel(w, r, ps, &params); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	chain, err := mesh.GetChains().GetChain(params.RelayNetworkID)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	mesh.GetLogger().Debug(
		fmt.Sprintf(
			"executing simulated relay of chain %s",
			chain.ID,
		),
	)
	// do basic http request on the relay
	res, er, _ := mesh.ExecuteBlockchainHTTPRequest(params.Payload, chain)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}
	WriteResponse(w, res, r.URL.Path, r.Host)
}

// ++++++++++++++++++++ MESH CLIENT - PRIVATE ROUTES ++++++++++++++++++++

// meshHealth - handle mesh health request
// path: /v1/private/mesh/health
func meshHealth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	res := mesh.HealthResponse{
		Version:   mesh.AppVersion,
		Servicers: mesh.ServicersSize(),
		FullNodes: mesh.NodesSize(),
	}
	j, er := json.Marshal(res)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshChains - return load chains from app.GlobalMeshConfig.ChainsName file
// path: /v1/private/mesh/chains
func meshChains(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	c := make([]pocketTypes.HostedBlockchain, 0)

	for _, chain := range mesh.GetChains().M {
		c = append(c, chain)
	}

	j, err := json.Marshal(c)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNode - return servicer node configured by servicer_priv_key.json - return address
// path: /v1/private/mesh/servicer
func meshServicerNode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	servicers := make([]types4.PublicPocketNode, 0)

	for _, a := range mesh.GetServicerLists() {
		servicers = append(servicers, types4.PublicPocketNode{
			Address: a,
		})
	}

	j, err := json.Marshal(servicers)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	WriteRaw(w, string(j), r.URL.Path, r.Host)
}

// meshUpdateChains - update chains in memory and also chains.json file.
// path: /v1/private/mesh/updatechains
func meshUpdateChains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}

	var hostedChainsSlice []pocketTypes.HostedBlockchain
	if err := PopModel(w, r, ps, &hostedChainsSlice); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	chains, e := mesh.UpdateChains(hostedChainsSlice)

	if e != nil {
		WriteErrorResponse(w, 400, e.Error())
		return
	}

	j, er := json.Marshal(chains.M)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)

}

// meshStop - gracefully stop mesh rpc server. Also, this should stop new relays and wait/flush all pending relays, otherwise they will get loose.
// path: /v1/private/mesh/stop
func meshStop(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !isAuthorized(w, r) {
		return
	}
	mesh.StopRPC()
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

	token := r.Header.Get("Authorization")
	if token != app.AuthToken.Value {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := mesh.RPCRelayResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
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

	_, dispatch, err := app.PCA.HandleRelay(relay, true)
	if err != nil {
		response := mesh.RPCRelayResult{
			Success:  false,
			Error:    err,
			Dispatch: dispatch,
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := mesh.RPCRelayResult{
		Success:  true,
		Dispatch: dispatch,
	}
	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNodeSession - receive requests from mesh node to validate a session for an app/servicer/blockchain on the servicer node data
// path: /v1/private/mesh/session
func meshServicerNodeSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var session pocketTypes.MeshSession

	token := r.Header.Get("Authorization")
	if token != app.AuthToken.Value {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	if verify == "true" {
		code := 200
		// useful just to test that mesh node is able to reach servicer
		response := mesh.RPCSessionResult{
			Success:  true,
			Error:    nil,
			Dispatch: nil,
		}

		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := PopModel(w, r, ps, &session); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	res, err := app.PCA.HandleMeshSession(session)

	if err != nil {
		response := mesh.RPCSessionResult{
			Success: false,
			Error:   mesh.NewSdkErrorFromPocketSdkError(err),
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	dispatch := mesh.DispatchResponse{
		BlockHeight: res.Session.BlockHeight,
		Session: mesh.DispatchSession{
			Header: res.Session.Session.SessionHeader,
			Key:    hex.EncodeToString(res.Session.Session.SessionKey),
			Nodes:  make([]mesh.DispatchSessionNode, 0),
		},
	}

	for i := range res.Session.Session.SessionNodes {
		sNode, ok := res.Session.Session.SessionNodes[i].(nodesTypes.Validator)
		if !ok {
			continue
		}
		dispatch.Session.Nodes = append(dispatch.Session.Nodes, mesh.DispatchSessionNode{
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

	response := mesh.RPCSessionResult{
		Success:         true,
		Dispatch:        &dispatch,
		RemainingRelays: json.Number(strconv.FormatInt(res.RemainingRelays, 10)),
	}
	j, er := json.Marshal(response)
	if er != nil {
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// meshServicerNodeCheck - receive requests from mesh node to validate servicers, chains and health status
// path: /v1/private/mesh/check
func meshServicerNodeCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var checkPayload mesh.CheckPayload

	token := r.Header.Get("Authorization")
	if token != app.AuthToken.Value {
		WriteErrorResponse(w, 401, "wrong authtoken: "+token)
		return
	}

	verify := r.URL.Query().Get("verify")
	// useful just to test that mesh node is able to reach servicer - this payload should be ignored
	if verify == "true" {
		code := 200
		j, _ := json.Marshal(map[string]interface{}{})
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, code)
		return
	}

	if err := PopModel(w, r, ps, &checkPayload); err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	chainsMap, err := app.PCA.QueryHostedChains()
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	health := app.PCA.QueryHealth(APIVersion)
	latestHeight := app.PCA.BaseApp.LastBlockHeight()

	paramReturn, err := app.PCA.QueryParam(latestHeight, "pos/BlocksPerSession")
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	blocksPerSession, err := strconv.ParseInt(paramReturn.Value, 10, 0)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
		return
	}

	if err != nil {
		response := mesh.RPCSessionResult{
			Success: false,
			Error:   mesh.NewSdkErrorFromPocketSdkError(sdk.ErrInternal(err.Error())),
		}
		j, _ := json.Marshal(response)
		WriteJSONResponseWithCode(w, string(j), r.URL.Path, r.Host, 400)
		return
	}

	response := mesh.CheckResponse{
		Success:          true,
		Status:           health,
		Servicers:        true,
		Chains:           true,
		BlocksPerSession: blocksPerSession,
		WrongServicers:   make([]string, 0),
		WrongChains:      make([]string, 0),
	}

	for _, address := range checkPayload.Servicers {
		if pocketTypes.GlobalPocketConfig.LeanPocket {
			// if lean pocket enabled, grab the targeted servicer through the relay proof
			nodeAddress, e1 := sdk.AddressFromHex(address)
			if e1 != nil {
				WriteErrorResponse(w, 400, "could not convert servicer hex")
				return
			}
			_, e2 := pocketTypes.GetPocketNodeByAddress(&nodeAddress)
			if e2 != nil {
				response.Servicers = false
				response.WrongServicers = append(response.WrongServicers, address)
			}
		} else {
			// get self node (your validator) from the current state
			node := pocketTypes.GetPocketNode()
			nodeAddress := node.GetAddress()
			if nodeAddress.String() != address {
				response.Servicers = false
				response.WrongServicers = append(response.WrongServicers, address)
			}
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
		WriteErrorResponse(w, 400, er.Error())
		return
	}

	WriteJSONResponse(w, string(j), r.URL.Path, r.Host)
}

// ReuseBody - transform request body in a reusable reader to allow multiple source read it.
func reuseBody(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rr, err := mesh.NewReusableReader(r.Body)
		if err != nil {
			WriteErrorResponse(w, 500, fmt.Sprintf("error in RPC Handler WriteErrorResponse: %v", err))
		} else {
			r.Body = io.NopCloser(rr)
			handler(w, r, ps)
		}
	}
}

// getMeshRoutes - return routes that will be handled/proxied by mesh rpc server
func getMeshRoutes(simulation bool) Routes {
	routes := Routes{
		// Proxy
		Route{Name: "AppVersion", Method: "GET", Path: "/v1", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "Health", Method: "GET", Path: "/v1/health", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "Challenge", Method: "POST", Path: "/v1/client/challenge", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "ChallengeCORS", Method: "OPTIONS", Path: "/v1/client/challenge", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "HandleDispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "HandleDispatchCORS", Method: "OPTIONS", Path: "/v1/client/dispatch", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "Stop", Method: "POST", Path: "/v1/private/stop", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryChains", Method: "POST", Path: "/v1/private/chains", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryAccounts", Method: "POST", Path: "/v1/query/accounts", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryAccountTxs", Method: "POST", Path: "/v1/query/accounttxs", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryACL", Method: "POST", Path: "/v1/query/acl", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryAllParams", Method: "POST", Path: "/v1/query/allparams", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryBlockTxs", Method: "POST", Path: "/v1/query/blocktxs", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryDAOOwner", Method: "POST", Path: "/v1/query/daoowner", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryNodeClaim", Method: "POST", Path: "/v1/query/nodeclaim", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryNodeClaims", Method: "POST", Path: "/v1/query/nodeclaims", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryParam", Method: "POST", Path: "/v1/query/param", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryState", Method: "POST", Path: "/v1/query/state", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryUpgrade", Method: "POST", Path: "/v1/query/upgrade", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QuerySigningInfo", Method: "POST", Path: "/v1/query/signinginfo", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "LocalNodes", Method: "POST", Path: "/v1/private/nodes", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryUnconfirmedTxs", Method: "POST", Path: "/v1/query/unconfirmedtxs", HandlerFunc: mesh.ProxyRequest},
		Route{Name: "QueryUnconfirmedTx", Method: "POST", Path: "/v1/query/unconfirmedtx", HandlerFunc: mesh.ProxyRequest},
		// mesh public route to handle relays
		Route{Name: "MeshService", Method: "POST", Path: "/v1/client/relay", HandlerFunc: reuseBody(meshNodeRelay)},
		// mesh private routes
		Route{Name: "MeshHealth", Method: "GET", Path: "/v1/private/mesh/health", HandlerFunc: meshHealth},
		Route{Name: "QueryMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/chains", HandlerFunc: meshChains},
		Route{Name: "MeshNodeServicer", Method: "POST", Path: "/v1/private/mesh/servicers", HandlerFunc: meshServicerNode},
		Route{Name: "UpdateMeshNodeChains", Method: "POST", Path: "/v1/private/mesh/updatechains", HandlerFunc: meshUpdateChains},
		Route{Name: "StopMeshNode", Method: "POST", Path: "/v1/private/mesh/stop", HandlerFunc: meshStop},
	}

	// check if simulation is turn on
	if simulation {
		simRoute := Route{Name: "SimulateRequest", Method: "POST", Path: "/v1/client/sim", HandlerFunc: meshSimulateRelay}
		routes = append(routes, simRoute)
	}

	return routes
}

// IsAuthorized - check if the request is authorized using authToken of the auth.json file
func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get(mesh.AuthorizationHeader)
	if token == mesh.GetAuthToken() {
		return true
	} else {
		WriteErrorResponse(w, 401, "wrong Authorization: "+token)
		return false
	}
}

// GetServicerMeshRoutes - return routes that need to be added to servicer to allow mesh node to communicate with.
func GetServicerMeshRoutes() Routes {
	routes := Routes{
		{Name: "MeshRelay", Method: "POST", Path: mesh.ServicerRelayEndpoint, HandlerFunc: meshServicerNodeRelay},
		{Name: "MeshSession", Method: "POST", Path: mesh.ServicerSessionEndpoint, HandlerFunc: meshServicerNodeSession},
		{Name: "MeshCheck", Method: "POST", Path: mesh.ServicerCheckEndpoint, HandlerFunc: meshServicerNodeCheck},
	}

	return routes
}

// StartMeshRPC - encapsulate mesh.StartRPC
func StartMeshRPC(simulation bool) {
	mesh.StartRPC(Router(getMeshRoutes(simulation)))
}
