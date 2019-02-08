package client

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// DISCLAIMER: This is for the centralized dispatcher of Pocket core mvp, may be removed for production

// "Register" handles the localhost:<client-port>/v1/register call.
func Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	n := node.Node{}
	if err := shared.PopModel(w, r, ps, &n); err != nil {
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	if node.EnsureWL(node.SWL(), n.GID) {
		// add to peerlist
		node.PeerList().Add(n)
		// add to dispatch peers
		node.DispatchPeers().Add(n)
		// write to db
		if _, err := db.NewDB().Add(n); err != nil {
			fmt.Println(err.Error())
			shared.WriteErrorResponse(w, 500, "unable to write peer to database")
			return
		}
		// write response
		shared.WriteJSONResponse(w, "Success! Your node is now registered in the Pocket Network")
		return
	}
	shared.WriteErrorResponse(w, 401, "Invalid credentials")
}

// "Register" handles the localhost:<client-port>/v1/register call.
func UnRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	n := node.Node{}
	if err := shared.PopModel(w, r, ps, &n); err != nil {
		shared.WriteErrorResponse(w, 400, err.Error())
		return
	}
	// remove from peerlist
	node.PeerList().Remove(n)
	// remove from dispatch peers
	node.DispatchPeers().Delete(n)
	// delete from database
	if _, err := db.NewDB().Remove(n); err != nil {
		shared.WriteErrorResponse(w, 500, "unable to remove peer from database")
		return
	}
	shared.WriteJSONResponse(w, "Success! Your node is now unregistered in the Pocket Network")
}

func RegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "Register", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}

func UnRegisterInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	info := shared.InfoStruct(r, "UnRegister", node.Node{}, "Success or failure message")
	shared.WriteInfoResponse(w, info)
}
