package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/app"
)

var APIVersion = fmt.Sprintf("%s", app.AppVersion)

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
		Route{Name: "AppVersion", Method: "GET", Path: "/v1", HandlerFunc: Version},
		Route{Name: "Dispatch", Method: "POST", Path: "/v1/client/dispatch", HandlerFunc: Dispatch},
		Route{Name: "Service", Method: "POST", Path: "/v1/client/relay", HandlerFunc: Relay},
		Route{Name: "Challenge", Method: "POST", Path: "/v1/client/challenge", HandlerFunc: Challenge},
		Route{Name: "SendRawTx", Method: "POST", Path: "/v1/client/rawtx", HandlerFunc: SendRawTx},
		Route{Name: "QueryBlock", Method: "POST", Path: "/v1/query/block", HandlerFunc: Block},
		Route{Name: "QueryTX", Method: "POST", Path: "/v1/query/tx", HandlerFunc: Tx},
		Route{Name: "QueryAccountTXS", Method: "POST", Path: "/v1/query/accounttxs", HandlerFunc: AccountTxs},
		Route{Name: "QueryBlockTXS", Method: "POST", Path: "/v1/query/blocktxs", HandlerFunc: BlockTxs},
		Route{Name: "QueryHeight", Method: "POST", Path: "/v1/query/height", HandlerFunc: Height},
		Route{Name: "QueryBalance", Method: "POST", Path: "/v1/query/balance", HandlerFunc: Balance},
		Route{Name: "QueryAccount", Method: "POST", Path: "/v1/query/account", HandlerFunc: Account},
		Route{Name: "QueryNodes", Method: "POST", Path: "/v1/query/nodes", HandlerFunc: Nodes},
		Route{Name: "QueryNode", Method: "POST", Path: "/v1/query/node", HandlerFunc: Node},
		Route{Name: "QueryNodeParams", Method: "POST", Path: "/v1/query/nodeparams", HandlerFunc: NodeParams},
		Route{Name: "QueryNodeReceipts", Method: "POST", Path: "/v1/query/nodereceipts", HandlerFunc: NodeReceipts},
		Route{Name: "QueryNodeReceipt", Method: "POST", Path: "/v1/query/nodereceipt", HandlerFunc: NodeReceipt},
		Route{Name: "QueryApps", Method: "POST", Path: "/v1/query/apps", HandlerFunc: Apps},
		Route{Name: "QueryApp", Method: "POST", Path: "/v1/query/app", HandlerFunc: App},
		Route{Name: "QueryAppParams", Method: "POST", Path: "/v1/query/appparams", HandlerFunc: AppParams},
		Route{Name: "QueryPocketParams", Method: "POST", Path: "/v1/query/pocketparams", HandlerFunc: PocketParams},
		Route{Name: "QuerySupportedChains", Method: "POST", Path: "/v1/query/supportedchains", HandlerFunc: SupportedChains},
		Route{Name: "QuerySupply", Method: "POST", Path: "/v1/query/supply", HandlerFunc: Supply},
		Route{Name: "QueryDAOOwner", Method: "POST", Path: "/v1/query/daoowner", HandlerFunc: DAOOwner},
		Route{Name: "QueryUpgrade", Method: "POST", Path: "/v1/query/upgrade", HandlerFunc: Upgrade},
		Route{Name: "QueryACL", Method: "POST", Path: "/v1/query/acl", HandlerFunc: ACL},
		Route{Name: "QueryState", Method: "GET", Path: "/v1/query/state", HandlerFunc: State},
	}
	return routes
}

func WriteResponse(w http.ResponseWriter, jsn, path, ip string) {
	b, err := json.Marshal(jsn)
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

func WriteRaw(w http.ResponseWriter, jsn, path, ip string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := w.Write([]byte(jsn))
	if err != nil {
		panic(err)
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
