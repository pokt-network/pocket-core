package handlers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
 "GetNetworkInfo" handles the localhost:<client-port>/v1/network call.
 */
func GetNetworkInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetNetworkID" handles the localhost:<client-port>/v1/network/id call.
 */
func GetNetworkID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetPeerCount" handles the localhost:<client-port>/v1/network/peer_count call.
 */
func GetPeerCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetPeerList" handles the localhost:<client-port>/v1/network/peer_list call.
 */
func GetPeerList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}

/*
 "GetPeers" handles the localhost:<client-port>/v1/network/peers call.
 */
func GetPeers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO
}
