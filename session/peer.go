// This package is network code relating to pocket 'sessions'
package session

import (
	"github.com/pokt-network/pocket-core/node"
)

type Role int

const (
	VALIDATOR Role = iota + 1
	SERVICER
	DISPATCHER
)

type SessionPeer struct {
	Role      Role `json:"role"`
	node.Node `json:"node"`
}

// "AddSPeers" adds sessionPeers from a slice to the peerlist
func AddSPeers(spl []SessionPeer) {
	pl := node.GetPeerList() // get the peerlist
	for _, sp := range spl { // for each SessionPeer
		pl.Add(sp.Node) // add to the list
	}
}
