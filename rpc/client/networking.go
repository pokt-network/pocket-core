package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "GetNetworkInfo" handles the localhost:<client-port>/v1/network call.
func GetNetworkInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetNetworkID" handles the localhost:<client-port>/v1/network/id call.
func GetNetworkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetPeerCount" handles the localhost:<client-port>/v1/network/peer_count call.
func GetPeerCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetPeerList" handles the localhost:<client-port>/v1/network/peer_list call.
func GetPeerList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}

// "GetPeers" handles the localhost:<client-port>/v1/network/peers call.
func GetPeers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello! This endpoint is currently in development!")
}
