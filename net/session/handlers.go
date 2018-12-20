// This package is all message related code
package session

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/message"
)

/*
"HandleMessage" is the parent function that calls upon specific message handler functions
				based on the index
 */
func HandleMessage(message *message.Message, sender *Peer){
	switch message.Payload.ID {
	case 0: // simple print message (testing purposes)
		PrintPayload(message)
	case 1: // sessionPeers message
		SessionPeers(message, sender)
	}
}

func PrintPayload(message *message.Message){
	fmt.Println(message.Payload.Data)
}

func SessionPeers(message *message.Message, peer *Peer){
	var data []Peer
	err := json.Unmarshal(message.Payload.Data.([]byte),&data)
	if err!=nil {
		logs.NewLog("Unable to unmarshal session peers to an array of peers " + err.Error(), logs.PanicLevel, logs.JSONLogFormat)
	}
	for _,n:= range data{
		RegisterSessionPeerConnection(n)
	}
	RegisterSessionPeerConnection(*peer)
}
