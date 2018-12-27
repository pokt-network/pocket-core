// This package is network code relating to pocket 'sessions'
package sessio

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
	"sync"
)

// "sessionMessage.go" specifies the session related message payloads, constructors, and handlers

type NewSessionPayload struct {
	DevID string        `json:"devid"`								// the devID of the session
	Peers []SessionPeer `json:"peers"`								// the list of peers
}

/***********************************************************************************************************************
Message Constructors
*/

/*
"NewSessionMessage" creates a new message whose payload is info that can be used to derive a session
 */
func NewSessionMessage(nSPL NewSessionPayload) message.Message {
	return message.NewMessage(message.Payload{ID: 1, Data: nSPL})	// return message with NewSessionPayload
}

/***********************************************************************************************************************
Message Handlers
*/

/*
"HandleMessage" is a switch based on payloadID to direct to the proper handler functions
 */
 // TODO try to make this global
func HandleMessage(message *message.Message) {
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		NewPrintMessage(message)
	case 1: // sessionPeers message
		NewSessionMessageHandler(message)
	}
}

/*
Prints the payload of a message to the CLI (payload index 0)
*/
func NewPrintMessage(message *message.Message) {
	fmt.Println(message.Payload.Data)								// prints the payload to the CLI
}

/*
Handles the session peers message (payload index 1)
*/
//TODO confirm the nodes by referencing the blockchain
func NewSessionMessageHandler(message *message.Message) {
	sList := GetSessionList()
	nSPL := message.Payload.Data.(NewSessionPayload)				// extract the NewSessionPayload
	session := Session{DevID: nSPL.DevID, 							// create a session using developerID from payload
	ConnList: make(map[string]Connection), Mutex: sync.Mutex{}}
	session.NewConnections(nSPL.Peers)								// create new connections with each peer
	sList.AddSession(session)										// register the session
	AddSessionPeersToPeerlist(nSPL.Peers)							// add peers to peerList
}
