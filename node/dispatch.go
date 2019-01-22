package node

import (
	"fmt"
	"net/http"
	"sync"
)

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

type DispatchPeers struct {
	Map map[Blockchain]map[string]Node // <BlockchainOBJ><GID><Node>
	sync.Mutex
}

var (
	dp  *DispatchPeers
	one sync.Once
)

// "GetDispatchPeers" gets Map structure.
func GetDispatchPeers() *DispatchPeers {
	one.Do(func() {
		dp = &DispatchPeers{Map: make(map[Blockchain]map[string]Node)}
	})
	return dp
}

// "Add" adds a peer to the dispatchPeers structure
func (dp *DispatchPeers) Add(n Node) {
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

// "getPeers" returns a map of peers by blockchain.
func getPeers(dp DispatchPeers, bc Blockchain) map[string]Node {
	return dp.Map[bc]
}

// "GetPeers" returns a map of peers by blockchain.
func (dp DispatchPeers) GetPeers(bc Blockchain) map[string]Node {
	dp.Lock()
	defer dp.Unlock()
	return getPeers(dp, bc)
}

// "Remove" deletes a peer from DispatchPeers.
func (dp *DispatchPeers) Delete(n Node) {
	dp.Lock()
	defer dp.Unlock()
	for _, chain := range n.Blockchains {
		delete(getPeers(*dp, chain), n.GID) // delete node from map via GID
	}
}

// "Print" outputs the dispatchPeer structure to the CLI
func (dp *DispatchPeers) Print() {
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
func (dp *DispatchPeers) Check() {
	pl := GetPeerList()
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
