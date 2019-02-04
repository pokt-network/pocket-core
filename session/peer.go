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

type Peer struct {
	Role      Role `json:"role"`
	node.Node `json:"node"`
}

// "AddPeer" adds sessionPeers from a slice to the peerlist
func AddPeer(spl []Peer) {
	pl := node.PeerList()
	for _, sp := range spl {
		pl.Add(sp.Node)
	}
}
