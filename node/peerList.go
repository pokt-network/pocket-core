package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type PeerList struct {
	Map map[string]Node
	sync.Mutex
}

var (
	o  sync.Once
	pl *PeerList
)

// "GetPeerList" returns the global map of nodes.
func GetPeerList() *PeerList {
	o.Do(func() {
		pl = &PeerList{Map: make(map[string]Node)}
	})
	return pl
}

// "Add" adds a peer object to the global map.
func (pl *PeerList) Add(node Node) {
	pl.Lock()
	defer pl.Unlock()
	pl.Map[node.GID] = node
}

// "Remove" removes a peer object from the global map.
func (pl *PeerList) Remove(node Node) {
	pl.Lock()
	defer pl.Unlock()
	delete(pl.Map, node.GID)
}

// "Contains" returns true if node is within peerlist.
func (pl *PeerList) Contains(GID string) bool {
	pl.Lock()
	defer pl.Unlock()
	_, ok := pl.Map[GID]
	return ok
}

// "Count" returns the count of peers within the map.
func (pl *PeerList) Count() int {
	pl.Lock()
	defer pl.Unlock()
	return len(pl.Map)
}

// "Print" prints the peerlist to the CLI
func (pl *PeerList) Print() {
	fmt.Println(pl.Map)
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

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

func (pl *PeerList) CopyToDP() {
	pl.Lock()
	defer pl.Unlock()
	for _, peer := range pl.Map {
		GetDispatchPeers().Add(peer)
	}
}
