// This package is contains the handler functions needed for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

// "account.go" defines all API handlers that are under the 'account' category.

/*
 "GetAccount" handles the localhost:<client-port>/v1/account call.
 */
func GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "IsAccountActive" handles the localhost:<client-port>/v1/account/active call.
 */
func IsAccountActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetAccountBalance" handles the localhost:<client-port>/v1/account/balance call.
 */
func GetAccountBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetDateJoined" handles the localhost:<client-port>/v1/account/joined call.
 */
func GetDateJoined(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetAccountKarma" handles the localhost:<client-port>/v1/account/karma call.
 */
func GetAccountKarma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetLastActive" handles the localhost:<client-port>/v1/account/last_active call.
 */
func GetLastActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetAccTxCount" handles the localhost:<client-port>/v1/account/transaction_count call.
 */
func GetAccTxCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetAccSessCount" handles the localhost:<client-port>/v1/account/session_count call.
 */
func GetAccSessCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetAccStatus" handles the localhost:<client-port>/v1/account/status call.
 */
func GetAccStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
