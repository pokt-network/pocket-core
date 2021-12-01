package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
)

var APIVersion = app.AppVersion

func StartRPC(port string, timeout int64, simulation, debug, allBlockTxs, hotReloadChains bool) {
	routes := GetRoutes()
	if simulation {
		simRoute := Route{Name: "SimulateRequest", Method: "POST", Path: "/v1/client/sim", HandlerFunc: SimRequest}
		routes = append(routes, simRoute)
	}

	if debug {
		routes = append(routes, Route{Name: "DebugBlock", Method: "GET", Path: "/debug/pprof/block", HandlerFunc: wrapperHandler(pprof.Handler(("block")))})
		routes = append(routes, Route{Name: "DebugCmd", Method: "GET", Path: "/debug/pprof/cmdline", HandlerFunc: wrapperHandlerFunc(pprof.Cmdline)})
		routes = append(routes, Route{Name: "DebugGoroutine", Method: "GET", Path: "/debug/pprof/goroutine", HandlerFunc: wrapperHandler(pprof.Handler(("goroutine")))})
		routes = append(routes, Route{Name: "DebugHeap", Method: "GET", Path: "/debug/pprof/heap", HandlerFunc: wrapperHandler(pprof.Handler(("heap")))})
		routes = append(routes, Route{Name: "DebugIndex", Method: "GET", Path: "/debug/pprof", HandlerFunc: wrapperHandlerFunc(pprof.Index)})
		routes = append(routes, Route{Name: "DebugProfile", Method: "GET", Path: "/debug/pprof/profile", HandlerFunc: wrapperHandlerFunc(pprof.Profile)})
		routes = append(routes, Route{Name: "DebugSymbol", Method: "GET", Path: "/debug/pprof/symbol", HandlerFunc: wrapperHandlerFunc(pprof.Symbol)})
		routes = append(routes, Route{Name: "DebugThreadCreate", Method: "GET", Path: "/debug/pprof/threadcreate", HandlerFunc: wrapperHandler(pprof.Handler(("threadcreate")))})
		routes = append(routes, Route{Name: "DebugTrace", Method: "GET", Path: "/debug/pprof/trace", HandlerFunc: wrapperHandlerFunc(pprof.Trace)})
		routes = append(routes, Route{Name: "FreeOsMemory", Method: "GET", Path: "/debug/freememory", HandlerFunc: FreeMemory})
		routes = append(routes, Route{Name: "MemStats", Method: "GET", Path: "/debug/memstats", HandlerFunc: MemStats})
		routes = append(routes, Route{Name: "QuerySecondUpgrade", Method: "POST", Path: "/debug/second", HandlerFunc: SecondUpgrade})
		routes = append(routes, Route{Name: "QueryValidatorByChain", Method: "POST", Path: "/debug/vbc", HandlerFunc: QueryValidatorsByChain})
	}

	if allBlockTxs {
		routes = append(routes, Route{Name: "QueryAllBlockTxs", Method: "POST", Path: "/v1/query/allblocktxs", HandlerFunc: AllBlockTxs})
	}

	//if hot reload is not enabled, enable manual reload.
	if !hotReloadChains {
		routes = append(routes, Route{Name: "UpdateChains", Method: "POST", Path: "/v1/private/updatechains", HandlerFunc: UpdateChains})
	}

	srv := &http.Server{
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      60 * time.Second,
		Addr:              ":" + port,
		Handler:           http.TimeoutHandler(Router(routes), time.Duration(timeout)*time.Millisecond, "Server Timeout Handling Request"),
	}
	log.Fatal(srv.ListenAndServe())
}

func Router(routes Routes) *httprouter.Router {
	router := httprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	return router
}

func cors(w *http.ResponseWriter, r *http.Request) (isOptions bool) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	return ((*r).Method == "OPTIONS")
}

type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc httprouter.Handle
}

type Routes []Route

func GetRoutes() Routes {
	routes := Routes{
		Route{Name: "AppVersion", Method: "GET", Path: "/v1", HandlerFunc: Version},
		Route{Name: "Challenge", Method: "POST", Path: "/v1/client/challenge", HandlerFunc: Challenge},
		Route{Name: "ChallengeCORS", Method: "OPTIONS", Path: "/v1/client/challenge", HandlerFunc: Challenge},
		Route{Name: "HandleDispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: Dispatch},
		Route{Name: "HandleDispatchCORS", Method: "OPTIONS", Path: "/v1/client/dispatch", HandlerFunc: Dispatch},
		Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: SendRawTx},
		Route{Name: "Service", Method: "POST", Path: "/v1/client/relay", HandlerFunc: Relay},
		Route{Name: "Stop", Method: "POST", Path: "/v1/private/stop", HandlerFunc: Stop},
		Route{Name: "ServiceCORS", Method: "OPTIONS", Path: "/v1/client/relay", HandlerFunc: Relay},
		Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: Account},
		Route{Name: "QueryAccountTxs", Method: "POST", Path: "/v1/query/accounttxs", HandlerFunc: AccountTxs},
		Route{Name: "QueryACL", Method: "POST", Path: "/v1/query/acl", HandlerFunc: ACL},
		Route{Name: "QueryAllParams", Method: "POST", Path: "/v1/query/allparams", HandlerFunc: AllParams},
		Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: App},
		Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: AppParams},
		Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: Apps},
		Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: Balance},
		Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: Block},
		Route{Name: "QueryBlockTxs", Method: "POST", Path: "/v1/query/blocktxs", HandlerFunc: BlockTxs},
		Route{Name: "QueryDAOOwner", Method: "POST", Path: "/v1/query/daoowner", HandlerFunc: DAOOwner},
		Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: Height},
		Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: Node},
		Route{Name: "QueryNodeClaim", Method: "POST", Path: "/v1/query/nodeclaim", HandlerFunc: NodeClaim},
		Route{Name: "QueryNodeClaims", Method: "POST", Path: "/v1/query/nodeclaims", HandlerFunc: NodeClaims},
		Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: NodeParams},
		Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: Nodes},
		Route{Name: "QueryParam", Method: "POST", Path: "/v1/query/param", HandlerFunc: Param},
		Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: PocketParams},
		Route{Name: "QueryState", Method: "POST", Path: "/v1/query/state", HandlerFunc: State},
		Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: Supply},
		Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: SupportedChains},
		Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: Tx},
		Route{Name: "QueryUpgrade", Method: "POST", Path: "/v1/query/upgrade", HandlerFunc: Upgrade},
		Route{Name: "QuerySigningInfo", Method: "POST", Path: "/v1/query/signinginfo", HandlerFunc: SigningInfo},
		Route{Name: "QueryChains", Method: "POST", Path: "/v1/private/chains", HandlerFunc: Chains},
	}
	return routes
}

func FreeMemory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	debug.FreeOSMemory()
	WriteResponse(w, "MemoryFreed", r.URL.Path, r.Host)
}
func MemStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b, err := json.Marshal(m)
	if err != nil {
		WriteErrorResponse(w, 400, err.Error())
	}
	WriteResponse(w, string(b), r.URL.Path, r.Host)
}

func WriteResponse(w http.ResponseWriter, jsn, path, ip string) {
	b, err := json.Marshal(jsn)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, err := w.Write(b)
		if err != nil {
			fmt.Println(fmt.Errorf("error in RPC Handler WriteResponse: %v", err))
		}
	}
}

func WriteRaw(w http.ResponseWriter, jsn, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(jsn))
	if err != nil {
		fmt.Println(fmt.Errorf("error in RPC Handler WriteRaw: %v", err))
	}
}

func WriteJSONResponse(w http.ResponseWriter, jsn, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsn), &raw); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println(fmt.Errorf("error in RPC Handler WriteJSONResponse: %v", err))
		return
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(raw)
	if err != nil {
		fmt.Println(fmt.Errorf("error in RPC Handler WriteJSONResponse: %v", err))
		return
	}
}

func WriteJSONResponseWithCode(w http.ResponseWriter, jsn, path, ip string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsn), &raw); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println(fmt.Errorf("error in RPC Handler WriteJSONResponse: %v", err))
		return
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(raw)
	if err != nil {
		fmt.Println(fmt.Errorf("error in RPC Handler WriteJSONResponse: %v", err))
		return
	}
}

func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(&rpcError{
		Code:    errorCode,
		Message: errorMsg,
	})
	if err != nil {
		fmt.Println(fmt.Errorf("error in RPC Handler WriteErrorResponse: %v", err))
	}
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func PopModel(_ http.ResponseWriter, r *http.Request, _ httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}

func wrapperHandlerFunc(f func(http.ResponseWriter, *http.Request)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		f(w, r)
	}
}

func wrapperHandler(h http.Handler) httprouter.Handle {
	f := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
	return wrapperHandlerFunc(f)
}
