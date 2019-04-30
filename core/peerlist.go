package core

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/pokt-network/pocket-core/types"
)

type PeerList types.List

var (
	o              sync.Once
	globalPeerList *PeerList
)

// "PeerList" returns the global map of nodes.
func GetPeerList() *PeerList {
	o.Do(func() {
		globalPeerList = (*PeerList)(types.NewList())
	})
	return globalPeerList
}

// "Add" adds a peer object to the global map.
func (pl *PeerList) Add(node Node) {
	(*types.List)(pl).Add(node.GID, node)
}

// "Remove" removes a peer object from the global map.
func (pl *PeerList) Remove(node Node) {
	(*types.List)(pl).Remove(node.GID)
}

// "Contains" returns true if node is within peerlist.
func (pl *PeerList) Contains(gid string) bool {
	return (*types.List)(pl).Contains(gid)
}

// "Count" returns the count of peers within the map.
func (pl *PeerList) Count() int {
	return (*types.List)(pl).Count()
}

// "Print" prints the peerlist to the CLI.
func (pl *PeerList) Print() {
	(*types.List)(pl).Print()
}

// "Clear" removes all nodes from the map.
func (pl *PeerList) Clear() {
	(*types.List)(pl).Clear()
}

// "Set" clears all nodes from the map and sets the peerlist as the node slice.
func (pl *PeerList) Set(nodes []Node) {
	pl.Clear()
	for _, n := range nodes {
		pl.Add(n)
	}
}

// "ManualPeersFile" adds Map from a peers.json to the peerlist
func ManualPeersFile(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return manualPeersJSON(file)
}

// "manualPeersJSON" adds Map from a json []byte to the peerlist
func manualPeersJSON(b []byte) error {
	var nSlice []Node
	if err := json.Unmarshal(b, &nSlice); err != nil {
		return err
	}
	for _, n := range nSlice {
		pList := GetPeerList()
		pList.Add(n)
	}
	return nil
}
