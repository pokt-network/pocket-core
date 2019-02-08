package node

import (
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/pokt-network/pocket-core/types"
)

type List types.List

var (
	o  sync.Once
	pl *List
)

// "PeerList" returns the global map of nodes.
func PeerList() *List {
	o.Do(func() {
		pl = (*List)(types.NewList())
	})
	return pl
}

// "Add" adds a peer object to the global map.
func (pl *List) Add(node Node) {
	(*types.List)(pl).Add(node.GID, node)
}

// "Remove" removes a peer object from the global map.
func (pl *List) Remove(node Node) {
	(*types.List)(pl).Remove(node.GID)
}

// "Contains" returns true if node is within peerlist.
func (pl *List) Contains(gid string) bool {
	return (*types.List)(pl).Contains(gid)
}

// "Count" returns the count of peers within the map.
func (pl *List) Count() int {
	return (*types.List)(pl).Count()
}

// "Print" prints the peerlist to the CLI.
func (pl *List) Print() {
	(*types.List)(pl).Print()
}

// "Clear" removes all nodes from the map.
func (pl *List) Clear() {
	(*types.List)(pl).Clear()
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
		pList := PeerList()
		pList.Add(n)
	}
	return nil
}

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

func (pl *List) CopyToDP() {
	pl.Mux.Lock()
	defer pl.Mux.Unlock()
	for _, peer := range pl.M {
		DispatchPeers().Add(peer.(Node))
	}
}
