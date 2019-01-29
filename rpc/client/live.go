package client

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// DISCLAIMER: This is for the centralized dispatcher of Pocket core mvp, may be removed for production

// "Register" handles the localhost:<client-port>/v1/register call.
func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	n := &node.Node{}
	if err := shared.PopulateModelFromParams(w, r, ps, n); err != nil {
		shared.WriteJSONResponse(w, "500 error: "+err.Error())
		return
	}
	if node.EnsureWL(node.GetSWL(), n.GID) {
		// add to peerlist
		node.GetPeerList().Add(*n)
		// add to dispatch peers
		node.GetDispatchPeers().Add(*n)
		shared.WriteJSONResponse(w, "Success! Your node is now registered in the Pocket Network")
		return
	}
	shared.WriteJSONResponse(w, "Invalid credentials")
}

// "Register" handles the localhost:<client-port>/v1/register call.
func UnRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	n := &node.Node{}
	if err := shared.PopulateModelFromParams(w, r, ps, n); err != nil {
		shared.WriteJSONResponse(w, "500 error: "+err.Error())
		return
	}
	// remove from peerlist
	node.GetPeerList().Remove(*n)
	// remove from dispatch peers
	node.GetDispatchPeers().Delete(*n)
	shared.WriteJSONResponse(w, "Success! Your node is now unregistered in the Pocket Network")
}

func RegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "Register", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}

func UnRegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.CreateInfoStruct(r, "UnRegister", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
