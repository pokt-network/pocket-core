// This package deals with all things networking related.
package net

import (
	"github.com/pokt-network/pocket-core/node"
)

// "peers.go" specifies peer related code.
// TODO could convert to structure in the future to make more robust
// TODO add concurrency protection
var (
	peerList = make(map[string]node.Node) // only make the peerlist once
)

func GetPeerList() map[string]node.Node {
	return peerList
}

func AddNodePeerList(node node.Node) {
	peerList[node.GID] = node
}

func RemoveNodePeerList(node node.Node) {
	delete(peerList, node.GID)
}
