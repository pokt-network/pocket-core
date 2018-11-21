// This package is contains the handler functions needed for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/session"
	"github.com/pokt-network/pocket-core/util"
	"net/http"
)

// Define all API handlers that are under the 'dispatch' category within this file.

/*
 "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
 */
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
const(
	dispatchServeMethod string 	= "POST"
	dispatchServerReturn string = "DATA - Session ID"
	dispatchServeExample string	= "curl --data {devid:1234}' http://localhost:8546/v1/dispatch/serve"
)
var(
	dispatchServeParams =[]string{"devid"}
)
/*
 "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
 */
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dispatch := &Dispatch{}
	shared.PopulateModelFromParams(w,r,ps,dispatch)
	if session.SearchSessionList(dispatch.DevID)!=nil{
		// Session Found
		// Write the sessionKey
		// TODO should store sessionKey and return value if found
		shared.WriteResponse(w,util.BytesToHex(session.GenerateSessionKey(dispatch.DevID)))
	} else {
		// Session not found
		session.CreateNewSession(dispatch.DevID)
		session.SearchSessionList(dispatch.DevID)
		shared.WriteResponse(w,util.BytesToHex(session.GenerateSessionKey(dispatch.DevID)))
		// TODO store sessionKey
	}
	session.PrintSessionList()
}

/*
"DispatchServeInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
And provides the devlopers with an in-client reference to the API call
 */
func DispatchServeInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	info:=shared.APIReference{r.URL.String(), dispatchServeMethod,
	dispatchServeParams, dispatchServerReturn, dispatchServeExample}
	shared.WriteInfoResponse(w,info)
}
