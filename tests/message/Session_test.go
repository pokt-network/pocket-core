package message

import (
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/peers"
	"github.com/pokt-network/pocket-core/session"
	"testing"
	"time"
)

func TestSessionMessage(t *testing.T) {
	const DEVID, SNODE, VNODE, IP = "DUMMYDEVID", "SNODE", "VNODE", "localhost"
	// run message servers
	message.RunMessageServers()
	// create the dummy peers to send over the wire
	var sessPeerList []session.SessionPeer	// create list
	sNode := session.SessionPeer{Role: session.SERVICER, Node: node.Node{GID: SNODE}}
	vNode := session.SessionPeer{Role: session.VALIDATOR, Node: node.Node{GID: VNODE}}
	sessPeerList = append(sessPeerList, sNode, vNode)	// add nodes to list
	// create the message structure
	m:=message.NewSessionMessage(message.NewSessionPayload{DevID: DEVID, Peers: sessPeerList})
	// send the message over the wire
	message.SendMessage(message.RELAY,m,IP, message.NewSessionPayload{})
	time.Sleep(time.Second)
	// check for session count
	if session.GetSessionList().Count() == 0 {
		t.Fatalf("No sessions within list")
	}
	// check for the correct devID within the session
	if !session.GetSessionList().Contains(DEVID) {
		t.Fatalf("The session for "+ DEVID+" doesn't exist")
	}
	// check for any peers within the peerlist
	if peers.GetPeerList().Count() == 0 {
		t.Fatalf("No peers within peerlist")
	}
	// check for the correct peers within the peerlist
	if !peers.GetPeerList().Contains(SNODE) || !peers.GetPeerList().Contains(VNODE) {
		t.Fatalf(SNODE + " and " + VNODE + " do not exist within peerlist")
	}
	// check for any sessionPeers
	s := session.GetSessionList().Get(DEVID)
	if s.Peers.Count() == 0 {
		t.Fatal("There are no peers within the session")
	}
	// check for proper sessionPeers
	if s.Peers.Get(SNODE).GID != SNODE || s.Peers.Get(VNODE).GID != VNODE{
		t.Fatal(SNODE+ " and " + VNODE + " do not exist within the session peer list")
	}
}
