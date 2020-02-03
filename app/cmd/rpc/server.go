package rpc

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	APIVersion = "0.0.1"
)

func StartRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, Router(GetRoutes())))
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
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*r).Method == "OPTIONS" {
		return false
	}
	return true
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
		Route{Name: "Version", Method: "GET", Path: "/v1", HandlerFunc: Version},
		Route{Name: "Dispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: Dispatch},
		Route{Name: "Service", Method: "POST", Path: "/v1/client/relay", HandlerFunc: Relay},
		Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: SendRawTx},
		Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: Block},
		Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: Tx},
		Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: Height},
		Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: Balance},
		Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: Account},
		Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: Nodes},
		Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: Node},
		Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: NodeParams},
		Route{Name: "QueryNodeProofs", Method: "POST", Path: "/v1/query/nodeproofs", HandlerFunc: NodeProofs},
		Route{Name: "QueryNodeProof", Method: "POST", Path: "/v1/query/nodeproof", HandlerFunc: NodeProof},
		Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: Apps},
		Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: App},
		Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: AppParams},
		Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: PocketParams},
		Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: SupportedChains},
		Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: Supply},
	}
	return routes
}

func WriteResponse(w http.ResponseWriter, jsn, path, ip string) {
	b, err := json.MarshalIndent(jsn, "", "\t")
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		log.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		_, err := w.Write(b)
		if err != nil {
			panic(err)
		}
	}
}
func WriteJSONResponse(w http.ResponseWriter, jsn, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsn), &raw); err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		log.Println(err.Error())
	}
	json.NewEncoder(w).Encode(raw)
}

func WriteErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	err := json.NewEncoder(w).Encode(&rpcError{
		Code:    errorCode,
		Message: errorMsg,
	})
	log.Print(err)
	if err != nil {
		panic(err)
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
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
