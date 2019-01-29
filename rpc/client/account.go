// This package contains handler functions for the client API
package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetAccount" handles the localhost:<client-port>/v1/account call.
func GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "IsAccountActive" handles the localhost:<client-port>/v1/account/active call.
func IsAccountActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetAccountBalance" handles the localhost:<client-port>/v1/account/balance call.
func GetAccountBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetDateJoined" handles the localhost:<client-port>/v1/account/joined call.
func GetDateJoined(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetAccountKarma" handles the localhost:<client-port>/v1/account/karma call.
func GetAccountKarma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetLastActive" handles the localhost:<client-port>/v1/account/last_active call.
func GetLastActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetAccTxCount" handles the localhost:<client-port>/v1/account/transaction_count call.
func GetAccTxCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetAccSessCount" handles the localhost:<client-port>/v1/account/session_count call.
func GetAccSessCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "GetAccStatus" handles the localhost:<client-port>/v1/account/status call.
func GetAccStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}
