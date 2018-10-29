// This package is contains the handler functions needed for the Client API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Define all API handlers that are under the 'transaction' category within this file.

/*
 "txOptions" handles the localhost:<client-port>/v1/transaction call.
 */
func TxOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "getTxByHash" handles the localhost:<client-port>/v1/transaction/hash call.
 */
func GetTxByHash(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
