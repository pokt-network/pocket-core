package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "GetPersonalInfo" handles the localhost:<client-port>/v1/personal call.
 */
func GetPersonalInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "ListAccounts" handles the localhost:<client-port>/v1/personal/accounts call.
 */
func ListAccounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "PersonalNetOptions" handles the localhost:<client-port>/v1/personal/network call.
 */
func PersonalNetOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO <gives simple instructions>
}

/*
 "EnterNetwork" handles the localhost:<client-port>/v1/personal/network/enter call.
 */
func EnterNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "ExitNetwork" handles the localhost:<client-port>/v1/personal/network/exit call.
 */
func ExitNetwork(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetPrimaryAddr" handles the localhost:<client-port>/v1/personal/primary_address call.
 */
func GetPrimaryAddr(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "SendPOKT" handles the localhost:<client-port>/v1/personal/send call.
 */
func SendPOKT(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "SendPOKTRaw" handles the localhost:<client-port>/v1/personal/send/raw call.
 */
func SendPOKTRaw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "Sign" handles the localhost:<client-port>/v1/personal/sign call.
 */
func Sign(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "StakeOptions" handles the localhost:<client-port>/v1/personal/stake call.
 */
func StakeOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO <gives simple instructions>
}

/*
 "Stake" handles the localhost:<client-port>/v1/personal/stake/add call.
 */
func Stake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "UnStake" handles the localhost:<client-port>/v1/personal/stake/remove call.
 */
func UnStake(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
