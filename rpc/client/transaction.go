package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "txOptions" handles the localhost:<client-port>/v1/transaction call.
func TxOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "getTxByHash" handles the localhost:<client-port>/v1/transaction/hash call.
func GetTxByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}
