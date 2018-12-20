// This package is all message related code
package sessio

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/net"
)

// holds message handlers

/*
"HandleMessage" is the parent function that calls upon specific message handler functions
				based on the index
 */
func HandleMessage(message *message.Message, sender *Connection){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		PrintPayload(message)
	case 1: // sessionPeers message
		NewSessionMessageHandler(message, sender)
	}
}
/*
Prints the payload of a message
 */
func PrintPayload(message *message.Message){
	fmt.Println(message.Payload.Data)
}

/*
Handles the session peers message
 */
func NewSessionMessageHandler(message *message.Message, peer *Connection){ //TODO confirm the nodes by referencing the blockchain
	session:=message.Payload.Data
	// register the session
	RegisterSession(session.(Session))
	// extract the connection list
	connList :=session.(Session).ConnList
	// get peerlist
	peerList:=net.GetPeerList()
	// add sender to peerlist
	peerList[peer.GID]=peer.Node
	// get each peer from the connectionList
	for _, connection := range connList {
		// add node to peerlist
		peerList[connection.GID]=connection.Node
		// establish connection with each node
		connection.CreateConnection("3333",connection.RemoteIP, session.(Session))	// TODO allow for flexible/manual ports
	}
}
