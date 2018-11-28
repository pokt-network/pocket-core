// This package is contains the handler functions needed for the Relay API
package relay

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/net"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"github.com/pokt-network/pocket-core/session"
	"github.com/pokt-network/pocket-core/util"
	"math/big"
	"net/http"
	"sort"
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
	if session.SearchSessionList(dispatch.DevID) == nil {
		session.CreateNewSession(dispatch.DevID)
	}
	sessionKey := util.BytesToHex(session.GenerateSessionKey(dispatch.DevID)) // TODO should store the session key
	nodes := DispatchFind(sessionKey)
	res, err := json.Marshal(nodes)
	if err != nil {
		logs.NewLog("Couldn't convert node array to json array: "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
	shared.WriteRawJSONResponse(w, res)
	session.PrintSessionList()
}

/*
"DispatchFind" orders the nodes from smallest proximity from sessionKey to largest proximity to sessionKey
 */
func DispatchFind(sessionKey string) []node.Node {
	bigSessionKey := new(big.Int)                   				// create new big integer to store sessionKey in
	bigSessionKey.SetString(sessionKey,16)        				// convert hex string into big integer
	peerList := net.GetPeerList()                   				// get the global peerlist
	m := make(map[uint64]node.Node)                 				// map the nodes to the corresponding difference
	keys := make([]uint64, len(peerList))           				// store the keys (to easily sort)
	sortedNodes := make([]node.Node, len(peerList)) 				// resulting array that holds the sorted nodes ordered by difference
	var i = 0                                       				// loop count
	for gid, curNode := range peerList { 							// for each curNode in the peerlist
		id := new(big.Int)                                 			// setup a new big integer to hold the converted ID
		id.SetString(gid,16)                               	// convert the hex GID into a bigInteger for comparison
		difference := big.NewInt(0).Sub(bigSessionKey, id) 		// find the difference between the two
		difference.Abs(difference)                         			// take absolute of the difference for comparison
		m[difference.Uint64()] = curNode                   			// map the corresponding difference -> curNode
		keys[i] = difference.Uint64()                      			// store the difference in the keys array
		i++                                                			// increment the count
	}
	sort.Slice(keys,func(i,j int) bool {return keys[i]<keys[j]})	// sort the keys
	for i, k := range keys {										// after filling out the difference for all in peerList
		sortedNodes[i] = m[k]										// store the nodes in order by difference
	}
	return sortedNodes												// return the sorted order
}

/*
"DispatchServeInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
And provides the developers with an in-client reference to the API call
 */
func DispatchServeInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := shared.CreateInfoStruct(r, "DispatchServe", Dispatch{}, "sessionKey")
	shared.WriteInfoResponse(w, info)
}
