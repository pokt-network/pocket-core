// This package is contains the handler functions needed for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"net/http"
)

// "personal.go" defines API handlers that are under the 'personal' category within this file.

/*
 "GetPersonalInfo" handles the localhost:<client-port>/v1/personal call.
*/
func GetPersonalInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "ListAccounts" handles the localhost:<client-port>/v1/personal/accounts call.
*/
func ListAccounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "PersonalNetOptions" handles the localhost:<client-port>/v1/personal/network call.
*/
func PersonalNetOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "EnterNetwork" handles the localhost:<client-port>/v1/personal/network/enter call.
*/
func EnterNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "ExitNetwork" handles the localhost:<client-port>/v1/personal/network/exit call.
*/
func ExitNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetPrimaryAddr" handles the localhost:<client-port>/v1/personal/primary_address call.
*/
func GetPrimaryAddr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "SendPOKT" handles the localhost:<client-port>/v1/personal/send call.
*/
func SendPOKT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "SendPOKTRaw" handles the localhost:<client-port>/v1/personal/send/raw call.
*/
func SendPOKTRaw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "Sign" handles the localhost:<client-port>/v1/personal/sign call.
*/
func Sign(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "StakeOptions" handles the localhost:<client-port>/v1/personal/stake call.
*/
func StakeOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "Stake" handles the localhost:<client-port>/v1/personal/stake/add call.
*/
func Stake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "UnStake" handles the localhost:<client-port>/v1/personal/stake/remove call.
*/
func UnStake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
