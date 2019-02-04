package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "NetInfo" handles the localhost:<client-port>/v1/network call.
func NetInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "NetID" handles the localhost:<client-port>/v1/network/id call.
func NetID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "PeerCount" handles the localhost:<client-port>/v1/network/peer_count call.
func PeerCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "PeerList" handles the localhost:<client-port>/v1/network/peer_list call.
func PeerList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}

// "PL" handles the localhost:<client-port>/v1/network/peers call.
func Peers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteJSONResponse(w, "Hello! This endpoint is currently in development!")
}
