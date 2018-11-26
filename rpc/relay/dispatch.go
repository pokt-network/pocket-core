// This package is contains the handler functions needed for the Relay API
package relay

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/session"
	"github.com/pokt-network/pocket-core/util"
	"net/http"
)

// "dispatch.go" defines API handlers that are under the 'dispatch' category within this file.

/*
 "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
 */
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
 */
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dispatch := &Dispatch{}
	shared.PopulateModelFromParams(w, r, ps, dispatch)
	if session.SearchSessionList(dispatch.DevID) != nil {
		// Session Found
		// Write the sessionKey
		// TODO should store sessionKey and return value if found
		shared.WriteResponse(w, util.BytesToHex(session.GenerateSessionKey(dispatch.DevID)))
	} else {
		// Session not found
		session.CreateNewSession(dispatch.DevID)
		session.SearchSessionList(dispatch.DevID)
		shared.WriteResponse(w, util.BytesToHex(session.GenerateSessionKey(dispatch.DevID)))
		// TODO store sessionKey
	}
	session.PrintSessionList()
}

/*
"DispatchServeInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
And provides the developers with an in-client reference to the API call
 */
func DispatchServeInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := shared.CreateInfoStruct(r, "DispatchServe", Dispatch{}, "sessionKey")
	shared.WriteInfoResponse(w, info)
}
