// This package is network code relating to pocket 'sessions'
package sessio

import (
	"github.com/pokt-network/pocket-core/net/peers"
	"github.com/pokt-network/pocket-core/node"
)

// "sessionPeer.go" holds the sessionPeer structure, enum, and functions

type Role int

/*
These constants are essentially an enum structure for Peer 'Role'
 */
const (
	VALIDATOR Role = iota + 1
	SERVICER
	DISPATCHER // TODO
)

/*
structure for sessionPeer
 */
type SessionPeer struct {
	Role      Role 	`json:"role"`		// the nodes specific role within the session
	node.Node 		`json:"node"`		// the node object
}

/***********************************************************************************************************************
Session Peer Functions
*/

/*
"AddSessionPeersToPeerList" adds sessionPeers from a slice to the peerlist
 */
func AddSessionPeersToPeerlist(spl []SessionPeer) {
	pl := peers.GetPeerList()			// get the peerlist
	for _, sp := range spl {			// for each SessionPeer
		pl.AddPeer(sp.Node)				// add to the list
	}
}

/*
"AddSessionPeerToPeerList" adds a single sessionPeer to the peerList
 */
func AddSessionPeerToPeerlist(sp SessionPeer) {
	pl := peers.GetPeerList()			// get the peerlist
	pl.AddPeer(sp.Node)					// add the peer to the list
}
