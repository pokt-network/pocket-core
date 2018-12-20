// This package deals with all things networking related.
package net

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/node"
	"io/ioutil"
	"log"
	"sync"
)

// "peers.go" specifies peer related code.

// TODO document and reorder message indexing
// TODO need to gracefully connect the following concepts: PEERLIST -> SESSIONLIST -> SESSIONPEERLIST
// 		Each session has a sessionPeerList
// 		Each sessionPeerList is made of persistent connections to peers (not peers as defined in this file but peers as defined
// 		in net/session/peer.go) <- confusing right? That's why this needs to be fixed ASAP
// TODO turn all panics into error correction (do research into this, next RC)
// TODO restructure packages (next RC)
// TODO add logging (next RC)

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

func GetPeerCount() int{
	return len(GetPeerList())
}

func AddNodePeerList(node node.Node) {
	if peerList==nil{
		GetPeerList()
	}
	lock.Lock()										// concurrency protection 'only one thread can add at a time'
	defer lock.Unlock()
	if !peerlistContains(node.GID) { // if node not within peerlist
		peerList[node.GID] = node					// TODO could add update function
	}
}

func RemoveNodePeerList(node node.Node) {
	if peerList==nil{
		GetPeerList()
	}
	delete(peerList, node.GID)
}

func peerlistContains(GID string) bool{
	_, ok := peerList[GID]
	return ok
}

func PeerlistContains(GID string) bool{
	lock.Lock()										// concurrency protection 'only one thread can search at a time'
	defer lock.Unlock()
	if peerList==nil{
		GetPeerList()
	}
	_, ok := peerList[GID]
	return ok
}

func ManualPeersFile(filepath string){
	if peerList==nil{
		GetPeerList()
	}
	file, _ := ioutil.ReadFile(filepath)	// TODO error handling -> use custom logs
	var data [] node.Node
	err := json.Unmarshal(file,&data)
	if err!=nil{
		log.Fatal("Unable to unmarshal json from " + filepath)
	}
	for _,n:= range data{
		AddNodePeerList(n)
	}
}

func ManualPeersJSON(b []byte){
	if peerList==nil{
		GetPeerList()	// TODO error handling
	}
	var data [] node.Node
	json.Unmarshal(b,&data)
	for _,n:= range data{
		AddNodePeerList(n)
	}
}
