// This package is contains the handler functions needed for the Client API
package client

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
	"net/http"
)

// Define all API handlers that are under the 'client' category within this file.

/*
 "GetClient" handles the localhost:<client-port>/v1/client call.
 */
func GetClientInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetClientID" handles the localhost:<client-port>/v1/client/id call.
 */
func GetClientID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetClientVersion" handles the localhost:<client-port>/v1/client/version call.
 */
func GetClientVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetCliSyncStatus" handles the localhost:<client-port>/v1/client/syncing call.
 */
func GetCliSyncStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
