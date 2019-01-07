package message

import (
	"fmt"
	"github.com/pokt-network/pocket-core/session"
	"net"
	"sync"
)

func HandleMessage(message *Message, addr *net.UDPAddr){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		PrintMessage(message, addr)
	case 1: // sessionPeers message
		NewSessionMessageHandler(message)
	}
}

/*
Prints the payload of a message to the CLI (payload index 0)
*/
func PrintMessage(message *Message, addr *net.UDPAddr) {
	fmt.Println(message.Payload.Data, "from " + addr.IP.String() + ":", addr.Port)	// prints the payload to the CLI
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
