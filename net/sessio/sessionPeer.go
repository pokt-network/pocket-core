// This package deals with all things networking related.
package sessio

import (
	"github.com/pokt-network/pocket-core/net/peers"
	"github.com/pokt-network/pocket-core/node"
)
type Role int

const (
	VALIDATOR Role = iota+1
	SERVICER
	DISPATCHER // TODO
)

type SessionPeer struct {
	Role Role				`json:"role"`
	node.Node				`json:"node"`
}

/***********************************************************************************************************************
Session Peer Functions
 */

func AddSessionPeersToPeerlist(spl []SessionPeer){
	pl := peers.GetPeerList()
	for _,sp := range spl {
		pl.AddPeer(sp.Node)
	}
}

func AddSessionPeerToPeerlist(sp SessionPeer){
	pl := peers.GetPeerList()
	pl.AddPeer(sp.Node)
}



