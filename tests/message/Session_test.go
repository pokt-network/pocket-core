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
	const DEVID, DEVID2, SNODE, VNODE, IP = "SESSION1", "SESSION2", "SNODE", "VNODE", "localhost"
	// start two different sessions
	session1:=session.NewSession(DEVID)
	session2:=session.NewSession(DEVID2)
	// add them to the session list
	sessionList := session.GetSessionList()
	sessionList.AddSession(session1, session2)
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
	time.Sleep(time.Second*2)
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
	session1=session.GetSessionList().Get(DEVID)
	session2=session.GetSessionList().Get(DEVID2)
	if session1.Peers.Count() == 0 {
		t.Fatalf("There are no peers within the session")
	}
	// check for proper sessionPeers
	if session1.Peers.Get(SNODE).GID != SNODE || session1.Peers.Get(VNODE).GID != VNODE{
		t.Fatalf(SNODE+ " and " + VNODE + " do not exist within the session peer list")
	}
	if session2.Peers.Count() != 0 {
		t.Fatalf("Too many peers in session2")
	}
}

