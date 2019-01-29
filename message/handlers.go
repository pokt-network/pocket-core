package message

import (
	"net"

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
	sl := session.GetSessionList()
	// extract the SessionPL
	spl := message.Payload.Data.(SessionPL)
	// create a session using developerID from payload
	s := session.Session{DevID: spl.DevID,
		Peers: session.NewPeerList()}
	s.AddPeers(spl.Peers)
	// adds new session to sessionList and adds peers to the peerList
	sl.Add(s)
	session.AddPeer(spl.Peers)
}
