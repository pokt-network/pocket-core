// This package is all message related code
package sessio

import (
	"fmt"
	"github.com/pokt-network/pocket-core/message"
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
		SessionPeers(message, sender)
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
func SessionPeers(message *message.Message, peer *Connection){
	// TODO Fix this call
	// STEP 0 -> Add the sessionID to EmptySession Struct
	// STEP 1 -> Interpret the list of nodes and add to global peerlist
	// STEP 2 -> TODO *maybe* confirm the nodes by referencing the blockchain
	// STEP 2 -> Create individual connections with each corresponding node
	// STEP 3 -> Register the session within the sessionList
}
