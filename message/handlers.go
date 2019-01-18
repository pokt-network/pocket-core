package message

import (
	"net"
	"sync"

	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
)

// "HandleMSG" handles all incoming messages and dispatches the information to the appropriate sub-handler.
func HandleMSG(message *Message, addr *net.UDPAddr) {
	switch message.Payload.ID {
	case 1: // sessionPeers message
		sessionMSG(message)
	case 2: // enter network
		enterMSG(message)
	case 3: // exit network
		exitNetworkMessage(message)
	}
}

// "sessionMSG" handles an incoming message by deriving a new session from the payload.
func sessionMSG(message *Message) {
	sList := session.GetSessionList()
	// extract the SessionPL
	nSPL := message.Payload.Data.(SessionPL)
	// create a session using developerID from payload
	s := session.Session{DevID: nSPL.DevID,
		Peers: session.PeerList{List: make(map[string]session.SessionPeer)},
		Mutex: sync.Mutex{}}
	// TODO create new connections with each peer
	s.NewPeers(nSPL.Peers)
	// adds new session to sessionList and adds peers to the peerList
	sList.AddSession(s)
	session.AddSPeers(nSPL.Peers)
}

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

// "enterMSG" handles an incoming message by adding a new service node to the centralized dispatcher's list(s)
func enterMSG(message *Message) {
	// get node from payload
	n := message.Payload.Data.(EnterPL)
	// cross check whitelist
	if node.EnsureWL(node.GetSWL(), n.GID) {
		// add to peerlist
		node.GetPeerList().Add(n.Node)
		// add to dispatch peers
		node.GetDispatchPeers().Add(n.Node)
	}
}

// "exitMSG" handles an incoming message by removing a service node from the centralized dispatcher's list(s)
func exitNetworkMessage(message *Message) {
	// get node from payload
	n := message.Payload.Data.(ExitPL)
	// add to peerlist
	node.GetPeerList().Remove(n.Node)
	// add to dispatch peers
	node.GetDispatchPeers().Delete(n.Node)
}
