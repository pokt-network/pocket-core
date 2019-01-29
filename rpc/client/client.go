package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetClient" handles the localhost:<client-port>/v1/client call.
func CliInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "CliID" handles the localhost:<client-port>/v1/client/id call.
func CliID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "CliVersion" handles the localhost:<client-port>/v1/client/version call.
func CliVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "CliSyncStatus" handles the localhost:<client-port>/v1/client/syncing call.
func CliSyncStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}
