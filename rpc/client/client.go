package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetClient" handles the localhost:<client-port>/v1/client call.
func GetClientInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetClientID" handles the localhost:<client-port>/v1/client/id call.
func GetClientID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetClientVersion" handles the localhost:<client-port>/v1/client/version call.
func GetClientVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetCliSyncStatus" handles the localhost:<client-port>/v1/client/syncing call.
func GetCliSyncStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}
