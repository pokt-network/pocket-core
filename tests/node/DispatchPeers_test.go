package node

import (
	"github.com/pokt-network/pocket-core/node"
	"log"
	"testing"
)

func TestDispatchPeers(t *testing.T) {
	// create arbitrary blockchains
	ethereum := node.Blockchain{"ethereum", "1", "1.0"}
	rinkeby  := node.Blockchain{"ethereum","4", "1.0"}
	bitcoin := node.Blockchain{"bitcoin","1","1.0"}
	bitcoinv1 := node.Blockchain{"bitcoin","1","1.1"}
	bitcoinCash := node.Blockchain{"bitcoinCash","1","1.0"}
	// create arbitrary nodes
	node1:= node.Node{
		GID:"node1",
		Blockchains:[]node.Blockchain{ethereum, rinkeby, bitcoin}}
	node2:= node.Node{
		GID:"node2",
		Blockchains:[]node.Blockchain{ethereum, bitcoin, bitcoinv1}}
	node3:= node.Node{
		GID:"node3",
		Blockchains:[]node.Blockchain{bitcoinCash, ethereum, bitcoin}}
	// add them to dispatchPeers
	node.NewDispatchPeer(node1)
	node.NewDispatchPeer(node2)
	node.NewDispatchPeer(node3)
	// get node lists
	ethereumNodes:= node.GetPeersByBlockchain(ethereum)
	btcNodes:= node.GetPeersByBlockchain(bitcoin)
	btcNodesV1:= node.GetPeersByBlockchain(bitcoinv1)
	bchNodes := node.GetPeersByBlockchain(bitcoinCash)
	rinkebyNodes := node.GetPeersByBlockchain(rinkeby)
	// ensure each list has proper node count
	if len(ethereumNodes)!=3 || len(btcNodes)!=3 ||len(rinkebyNodes)!=1 && len(btcNodesV1)!=1 && len(bchNodes)!=1{
		log.Fatalf("Incorrect node count")
	}
	// ensure each node is on list
	if ethereumNodes[node1.GID].GID=="" || ethereumNodes[node2.GID].GID=="" || ethereumNodes[node3.GID].GID=="" {
		log.Fatalf("Missing nodes")
	}
	// print the dispatch peers for visibility
	node.PrintDispatchPeers()
}
