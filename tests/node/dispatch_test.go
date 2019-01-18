package node

import (
	"log"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/rpc"
)

func TestDispatchLiveness(t *testing.T) {
	// start API servers
	go rpc.StartRelayRPC(config.GetInstance().RRPCPort)
	time.Sleep(time.Second)
	// get peer list
	pl := node.GetPeerList()
	// create arbitrary nodes
	self := node.Node{GID: "self", IP: "localhost", RelayPort: config.GetInstance().RRPCPort}
	dead := node.Node{GID: "deadNode", IP: "0.0.0.0", RelayPort: "0"}
	// add self to peerlist
	pl.Add(self)
	// add dead node to peerlist
	pl.Add(dead)
	// check liveness port of self
	node.GetDispatchPeers().Check()
	// ensure that dead node is deleted and live node is still within list
	if !pl.Contains(self.GID) || pl.Contains(dead.GID) {
		t.Fatalf("Peerlist result not correct, expected: " + self.GID + " only, and not: " + dead.GID)
	}
	pl.Print()
}

func TestDispatchPeers(t *testing.T) {
	// create arbitrary blockchains
	ethereum := node.Blockchain{"ethereum", "1", "1.0"}
	rinkeby := node.Blockchain{"ethereum", "4", "1.0"}
	bitcoin := node.Blockchain{"bitcoin", "1", "1.0"}
	bitcoinv1 := node.Blockchain{"bitcoin", "1", "1.1"}
	bitcoinCash := node.Blockchain{"bitcoinCash", "1", "1.0"}
	// create arbitrary nodes
	node1 := node.Node{
		GID:         "node1",
		Blockchains: []node.Blockchain{ethereum, rinkeby, bitcoin}}
	node2 := node.Node{
		GID:         "node2",
		Blockchains: []node.Blockchain{ethereum, bitcoin, bitcoinv1}}
	node3 := node.Node{
		GID:         "node3",
		Blockchains: []node.Blockchain{bitcoinCash, ethereum, bitcoin}}
	// add them to dispatchPeers
	dp := node.GetDispatchPeers()
	dp.Add(node1)
	dp.Add(node2)
	dp.Add(node3)
	// get node lists
	ethereumNodes := dp.GetPeers(ethereum)
	btcNodes := dp.GetPeers(bitcoin)
	btcNodesV1 := dp.GetPeers(bitcoinv1)
	bchNodes := dp.GetPeers(bitcoinCash)
	rinkebyNodes := dp.GetPeers(rinkeby)
	// ensure each list has proper node count
	if len(ethereumNodes) != 3 || len(btcNodes) != 3 || len(rinkebyNodes) != 1 && len(btcNodesV1) != 1 && len(bchNodes) != 1 {
		log.Fatalf("Incorrect node count")
	}
	// ensure each node is on list
	if ethereumNodes[node1.GID].GID == "" || ethereumNodes[node2.GID].GID == "" || ethereumNodes[node3.GID].GID == "" {
		log.Fatalf("Missing nodes")
	}
	// print the dispatch peers for visibility
	dp.Print()
}
