// This package is all message related code
package sessio

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
	"sync"
)

type NewSessionPayload struct {
	DevID string			`json:"devid"`
	Peers []SessionPeer		`json:"peers"`
}

/***********************************************************************************************************************
Message Constructors
 */
func NewSessionMessage(nSPL NewSessionPayload) message.Message{
	return message.NewMessage(message.Payload{ID: 1, Data: nSPL})
}
/***********************************************************************************************************************
Message Handlers
 */

func HandleMessage(message *message.Message){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		NewPrintMessage(message)
	case 1: // sessionPeers message
		NewSessionMessageHandler(message)
	}
}

/*
Prints the payload of a message (index 0)
 */
func NewPrintMessage(message *message.Message){
	fmt.Println(message.Payload.Data)
}

/*
Handles the session peers message
 */
func NewSessionMessageHandler(message *message.Message){ //TODO confirm the nodes by referencing the blockchain
	sList := GetSessionList()
	// extract the NewSessionPayload
	nSPL:=message.Payload.Data.(NewSessionPayload)
	// create a session using developerID from payload
	session := Session{DevID:nSPL.DevID, ConnList: make(map[string]Connection), Mutex:sync.Mutex{}}
	// create new connections with each peer
	session.NewConnections(nSPL.Peers)
	// register the session
	sList.AddSession(session)
	// add peers to peerList
	AddSessionPeersToPeerlist(nSPL.Peers)
}

