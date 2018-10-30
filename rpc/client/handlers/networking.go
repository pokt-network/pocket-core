// This package is contains the handler functions needed for the Client API
package handlers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/pocket_network/pocket-core/rpc/shared"
	"net/http"
)

// Define all API handlers that are under the 'networking' category within this file.

/*
 "GetNetworkInfo" handles the localhost:<client-port>/v1/network call.
 */
func GetNetworkInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetNetworkID" handles the localhost:<client-port>/v1/network/id call.
 */
func GetNetworkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetPeerCount" handles the localhost:<client-port>/v1/network/peer_count call.
 */
func GetPeerCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetPeerList" handles the localhost:<client-port>/v1/network/peer_list call.
 */
func GetPeerList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "GetPeers" handles the localhost:<client-port>/v1/network/peers call.
 */
func GetPeers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}
