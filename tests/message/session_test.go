package message

import (
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
)

func TestSessionMessage(t *testing.T) {
	const DEVID, DEVID2, SNODE, VNODE, IP = "SESSION1", "SESSION2", "SNODE", "VNODE", "localhost"
	// start two different sessions
	session1 := session.NewSession(DEVID)
	session2 := session.NewSession(DEVID2)
	// add them to the session list
	sessionList := session.SessionList()
	sessionList.AddMulti(session1, session2)
	// run message servers
	message.StartServers()
	// create the dummy peers to send over the wire
	var sessPeerList []session.Peer // create list
	sNode := session.Peer{Role: session.SERVICER, Node: node.Node{GID: SNODE}}
	vNode := session.Peer{Role: session.VALIDATOR, Node: node.Node{GID: VNODE}}
	sessPeerList = append(sessPeerList, sNode, vNode) // add nodes to list
	// create the message structure
	m := message.NewSessionMessage(DEVID, sessPeerList)
	// send the message over the wire
	message.SendMessage(message.RELAY, m, IP, message.SessionPL{})
	time.Sleep(time.Second * 4)
	// check for session count
	if session.SessionList().Count() == 0 {
		t.Fatalf("No sessions within list")
	}
	// check for the correct devID within the session
	if !session.SessionList().Contains(DEVID) {
		t.Fatalf("The session for " + DEVID + " doesn't exist")
	}
	// check for any peers within the peerlist
	if node.PeerList().Count() == 0 {
		t.Fatalf("No peers within peerlist")
	}
	// check for the correct peers within the peerlist
	if !node.PeerList().Contains(SNODE) || !node.PeerList().Contains(VNODE) {
		t.Fatalf(SNODE + " and " + VNODE + " do not exist within peerlist")
	}
	session1 = session.SessionList().Get(DEVID)
	session2 = session.SessionList().Get(DEVID2)
	if session1.PL.Count() == 0 {
		t.Fatalf("There are no peers within the session")
	}
	// check for proper sessionPeers
	if session1.PL.Get(SNODE).GID != SNODE || session1.PL.Get(VNODE).GID != VNODE {
		t.Fatalf(SNODE + " and " + VNODE + " do not exist within the session peer list")
	}
	if session2.PL.Count() != 0 {
		t.Fatalf("Too many peers in session2")
	}
}
