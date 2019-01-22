package message

import (
	"net"
	"sync"
	
	"github.com/pokt-network/pocket-core/session"
)

// "HandleMSG" handles all incoming messages and dispatches the information to the appropriate sub-handler.
func HandleMSG(message *Message, addr *net.UDPAddr) {
	switch message.Payload.ID {
	case 1: // sessionPeers message
		sessionMSG(message)
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
