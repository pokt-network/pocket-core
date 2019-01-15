// This package is contains the handler functions needed for the Relay API
package relay

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc/shared"
	"math/big"
	"net/http"
	"sort"
	"strings"
)

// TODO fix dispatch serve example APIInformation
// "dispatch.go" defines API handlers that are under the 'dispatch' category within this file.

/*
 "DispatchOptions" handles the localhost:<relay-port>/v1/dispatch call.
*/
func DispatchOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shared.WriteResponse(w, "Hello, World!")
}

/*
 "DispatchServe" handles the localhost:<relay-port>/v1/dispatch/serve call.
 NOTE: this call has been augmented for the Pocket Core MVP Centralized Dispatcher
// TODO see if this can be done more efficiently
*/
func DispatchServe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dispatch := &Dispatch{}
	result := make(map[string][]string)
	shared.PopulateModelFromParams(w, r, ps, dispatch)
	if node.EnsureWL(node.GetDeveloperWhiteList(), dispatch.DevID, w) {
		for _, blockchain := range dispatch.Blockchains {
			ips := make([]string,0)
			nodes := node.GetPeersByBlockchain(blockchain)
			for _,n:= range nodes {
				ips=append(ips, n.IP+":"+n.RelayPort)
			}
			result[strings.ToUpper(blockchain.Name)+"V"+blockchain.Version+" | NetID "+blockchain.NetID]=ips
		}
		res, err := json.MarshalIndent(result,"","  ")
		if err != nil {
			fmt.Println(err)
			logs.NewLog("Couldn't convert node array to json array: "+err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		}
		shared.WriteRawJSONResponse(w, res)
	}
}

/*
"DispatchFind" orders the nodes from smallest proximity from sessionKey to largest proximity to sessionKey
// TODO convert to P2P -> currently just searches the peerlist
// TODO NEED a separate dispatch file with calls like these
*/
func DispatchFind(sessionKey string) []node.Node {
	bigSessionKey := new(big.Int)           					// create new big integer to store sessionKey in
	bigSessionKey.SetString(sessionKey, 16) 				// convert hex string into big integer
	peerList := node.GetPeerList()         						// get the global peerlist
	peerList.Lock()                         					// TODO currently locking the peerlist, however this will all change when p2p is integerated
	defer peerList.Unlock()
	m := make(map[uint64]node.Node)                      		// map the nodes to the corresponding difference
	keys := make([]uint64, len(peerList.List))           		// store the keys (to easily sort)
	sortedNodes := make([]node.Node, len(peerList.List)) 		// resulting array that holds the sorted nodes ordered by difference
	var i = 0                                            		// loop count
	for gid, curNode := range peerList.List {            		// for each curNode in the peerlist
		id := new(big.Int)                                 		// setup a new big integer to hold the converted ID
		id.SetString(gid, 16)                             // convert the hex GID into a bigInteger for comparison
		difference := big.NewInt(0).Sub(bigSessionKey, id) 	// find the difference between the two
		difference.Abs(difference)                         		// take absolute of the difference for comparison
		m[difference.Uint64()] = curNode                   		// map the corresponding difference -> curNode
		keys[i] = difference.Uint64()                      		// store the difference in the keys array
		i++                                                		// increment the count
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] }) // sort the keys
	for i, k := range keys {                                    // after filling out the difference for all in peerList
		sortedNodes[i] = m[k] 									// store the nodes in order by difference
	}
	return sortedNodes 											// return the sorted order
}

/*
"DispatchServeInfo" handles a get request to localhost:<relay-port>/v1/dispatch/serve call.
And provides the developers with an in-client reference to the API call
*/
func DispatchServeInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	info := shared.CreateInfoStruct(r, "DispatchServe", Dispatch{}, "list of service nodes")
	shared.WriteInfoResponse(w, info)
}
