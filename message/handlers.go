package message

import (
	"github.com/pokt-network/pocket-core/session"
	"net"
	"sync"
)

func HandleMessage(message *Message, addr *net.UDPAddr){
	switch message.Payload.ID {
	case 1: // sessionPeers message
		NewSessionMessageHandler(message)
	}
}

func NewSessionMessageHandler(message *Message) {
	sList := session.GetSessionList()
	nSPL := message.Payload.Data.(NewSessionPayload)				// extract the NewSessionPayload
	s := session.Session{DevID: nSPL.DevID, 						// create a session using developerID from payload
		Peers: session.SessionPeerList{List: make(map[string]session.SessionPeer)},
		Mutex: sync.Mutex{}}
	s.NewPeers(nSPL.Peers)                  						// create new connections with each peer
	sList.AddSession(s)                     						// register the session
	session.AddSessPeersToPL(nSPL.Peers) 							// add peers to peerList
}
