package node

import (
	"fmt"
	"net/http"
	"sync"
)

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

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

// "Peers" returns a map of peers by blockchain.
func (dp DPeers) Peers(bc Blockchain) map[string]Node {
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

// "Check" checks each service node's liveness.
func (dp *DPeers) Check() {
	pl := PeerList()
	for _, p := range pl.M {
		if !isAlive(p.(Node)) {
			// try again
			if !isAlive(p.(Node)) {
				pl.Remove(p.(Node))
				dp.Delete(p.(Node))
			}
		}
	}
}

// "Clear" removes all nodes from the list.
func (dp *DPeers) Clear() {
	dp.Lock()
	defer dp.Unlock()
	dp.Map = make(map[Blockchain]map[string]Node)
}

// "isAlive" checks a node and returns the status of that check.
func isAlive(n Node) bool { // TODO handle scenarios where the error is on the dispatch node side
	if resp, err := check(n); err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// "check" tests a node by doing an HTTP GET to API.
func check(n Node) (*http.Response, error) {
	return http.Get("http://" + n.IP + ":" + n.RelayPort + "/v1/")
}
