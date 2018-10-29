// This package is contains the handler functions needed for the Client API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Define all API handlers that are under the 'account' category within this file.

/*
 "GetAccount" handles the localhost:<client-port>/v1/account call.
 */
func GetAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "IsAccountActive" handles the localhost:<client-port>/v1/account/active call.
 */
func IsAccountActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetAccountBalance" handles the localhost:<client-port>/v1/account/balance call.
 */
func GetAccountBalance(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetDateJoined" handles the localhost:<client-port>/v1/account/joined call.
 */
func GetDateJoined(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetAccountKarma" handles the localhost:<client-port>/v1/account/karma call.
 */
func GetAccountKarma(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetLastActive" handles the localhost:<client-port>/v1/account/last_active call.
 */
func GetLastActive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetAccTxCount" handles the localhost:<client-port>/v1/account/transaction_count call.
 */
func GetAccTxCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetAccSessCount" handles the localhost:<client-port>/v1/account/session_count call.
 */
func GetAccSessCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetAccStatus" handles the localhost:<client-port>/v1/account/status call.
 */
func GetAccStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
