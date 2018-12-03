// This package deals with all things networking related.
package net

import (
	"github.com/pokt-network/pocket-core/node"
	"sync"
)

// "peers.go" specifies peer related code.
// TODO could convert to structure in the future to make more robust
var (
	once     sync.Once
	peerList map[string]node.Node
	lock sync.Mutex
)

func GetPeerList() map[string]node.Node {
	if peerList == nil {
		once.Do(func() {
			peerList = make(map[string]node.Node) // only make the peerlist once
		})
	}
	return peerList
}

func AddNodePeerList(node node.Node) {
	lock.Lock()										// concurrency protection 'only one thread can add at a time'
	defer lock.Unlock()
	if !peerlistContains(node.GID) { // if node not within peerlist
		peerList[node.GID] = node					// TODO could add update function
	}
}

func RemoveNodePeerList(node node.Node) {
	delete(peerList, node.GID)
}

func peerlistContains(GID string) bool{
	_, ok := peerList[GID]
	return ok
}

func PeerlistContains(GID string) bool{
	lock.Lock()										// concurrency protection 'only one thread can search at a time'
	defer lock.Unlock()
	_, ok := peerList[GID]
	return ok
}
