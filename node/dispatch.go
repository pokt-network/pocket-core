package node

import (
	"fmt"
	"sync"
)

type DPeers struct {
	Map map[Blockchain]map[string]Node // <BlockchainOBJ><GID><Node>
	sync.Mutex
}

var (
	dp  *DPeers
	one sync.Once
)

// "DispatchPeers" gets Map structure.
func DispatchPeers() *DPeers {
	one.Do(func() {
		dp = &DPeers{Map: make(map[Blockchain]map[string]Node)}
	})
	return dp
}

// "Add" adds a peer to the dispatchPeers structure
func (dp *DPeers) Add(n Node) {
	dp.Lock()
	defer dp.Unlock()
	for _, bchain := range n.Blockchains {
		// type map[GID]Node
		nodes := dp.Map[bchain]
		// if bchain not within Map
		if nodes == nil {
			// add new node to empty map
			dp.Map[bchain] = map[string]Node{n.GID: n}
			continue
		}
		// add node to inner map
		nodes[n.GID] = n
		// update outer map
		dp.Map[bchain] = nodes
	}
}

// "peers" returns a map of peers by blockchain.
func peers(dp DPeers, bc Blockchain) map[string]Node {
	return dp.Map[bc]
}

// "PeersByChain" returns a map of peers by blockchain.
func (dp DPeers) PeersByChain(bc Blockchain) map[string]Node {
	dp.Lock()
	defer dp.Unlock()
	return peers(dp, bc)
}

// "Remove" deletes a peer from DPeers.
func (dp *DPeers) Delete(n Node) {
	dp.Lock()
	defer dp.Unlock()
	for _, chain := range n.Blockchains {
		delete(peers(*dp, chain), n.GID) // delete node from map via GID
	}
}

// "Print" outputs the dispatchPeer structure to the CLI
func (dp *DPeers) Print() {
	dp.Lock()
	defer dp.Unlock()
	// [blockchain]Map of nodes
	for bc, nMap := range dp.Map {
		fmt.Println(bc.Name, "Version:", bc.Version, "NetID:", bc.NetID)
		fmt.Println("  GID's:")
		// GID in Map
		for gid := range nMap {
			fmt.Println("   ", gid)
		}
		fmt.Println("")
	}
}

// "Clear" removes all nodes from the map.
func (dp *DPeers) Clear() {
	dp.Lock()
	defer dp.Unlock()
	dp.Map = make(map[Blockchain]map[string]Node)
}
