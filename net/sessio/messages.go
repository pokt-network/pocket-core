// This package is all message related code
package sessio

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net/peers"
)
/***********************************************************************************************************************
Message Constructors
 */
func NewSessionMessage(session Session) message.Message{
	return message.NewMessage(message.Payload{ID: 1, Data: session})
}
/***********************************************************************************************************************
Message Handlers
 */

func HandleMessage(message *message.Message, sender *Connection){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		NewPrintMessage(message)
	case 1: // sessionPeers message
		NewSessionMessageHandler(message, sender)
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
func NewSessionMessageHandler(message *message.Message, peer *Connection){ //TODO confirm the nodes by referencing the blockchain
	sList := GetSessionList()
	pList := peers.GetPeerList()
	session:=message.Payload.Data.(Session)
	// register the session
	sList.AddSession(session)
	// extract the connection List
	connList :=session.ConnList
	// add sender to peerlist
	pList.AddPeer(peer.Node)
	// get each peer from the connectionList
	for _, connection := range connList {
		// add node to peerlist
		pList.AddPeer(connection.Node)
		// establish connection with each node
		session.Dial("3333",connection.RemoteIP, connection) // TODO allow for flexible/manual ports
	}
}
