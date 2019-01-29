// This package contains handler functions for the client API
package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

<<<<<<< HEAD
// "Account" handles the localhost:<client-port>/v1/account call.
func Account(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
=======
// "GetAccount" handles the localhost:<client-port>/v1/account call.
func GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
>>>>>>> fixed all possible todos throughout package
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "IsAccountActive" handles the localhost:<client-port>/v1/account/active call.
func IsAccountActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

<<<<<<< HEAD
// "Balance" handles the localhost:<client-port>/v1/account/balance call.
func Balance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "DateJoined" handles the localhost:<client-port>/v1/account/joined call.
func DateJoined(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "AcountKarma" handles the localhost:<client-port>/v1/account/karma call.
func AcountKarma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "LastActive" handles the localhost:<client-port>/v1/account/last_active call.
func LastActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "AcctTXCount" handles the localhost:<client-port>/v1/account/transaction_count call.
func AcctTXCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "AccSessCount" handles the localhost:<client-port>/v1/account/session_count call.
func AccSessCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}

// "AccStatus" handles the localhost:<client-port>/v1/account/status call.
func AccStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
=======
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
>>>>>>> fixed all possible todos throughout package
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development")
}
