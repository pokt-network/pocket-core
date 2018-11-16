// This package is contains the handler functions needed for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

// Define all API handlers that are under the 'transaction' category within this file.

/*
 "txOptions" handles the localhost:<client-port>/v1/transaction call.
 */
func TxOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "getTxByHash" handles the localhost:<client-port>/v1/transaction/hash call.
 */
func GetTxByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
