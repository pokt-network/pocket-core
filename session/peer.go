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

// "AddPeers" adds sessionPeers from a slice to the Global peerlist
func AddPeers(spl []Peer) {
	pl := node.PeerList()
	for _, sp := range spl {
		pl.Add(sp.Node)
	}
}
