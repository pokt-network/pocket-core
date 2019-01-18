package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetPersonalInfo" handles the localhost:<client-port>/v1/personal call.
func GetPersonalInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "ListAccounts" handles the localhost:<client-port>/v1/personal/accounts call.
func ListAccounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "PersonalNetOptions" handles the localhost:<client-port>/v1/personal/network call.
func PersonalNetOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "EnterNetwork" handles the localhost:<client-port>/v1/personal/network/enter call.
func EnterNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "ExitNetwork" handles the localhost:<client-port>/v1/personal/network/exit call.
func ExitNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetPrimaryAddr" handles the localhost:<client-port>/v1/personal/primary_address call.
func GetPrimaryAddr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "SendPOKT" handles the localhost:<client-port>/v1/personal/send call.

func SendPOKT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "SendPOKTRaw" handles the localhost:<client-port>/v1/personal/send/raw call.
func SendPOKTRaw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "Sign" handles the localhost:<client-port>/v1/personal/sign call.
func Sign(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "StakeOptions" handles the localhost:<client-port>/v1/personal/stake call.
func StakeOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "Stake" handles the localhost:<client-port>/v1/personal/stake/add call.
func Stake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "UnStake" handles the localhost:<client-port>/v1/personal/stake/remove call.
func UnStake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}
