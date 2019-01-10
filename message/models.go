// This package is all message related code
package message

import (
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
)

// "models.go" holds all of the structures for the message package

/*
NOTE: 	The ideology here in design is to maintain a simple message structure while swapping out the payload.

		For example: A message would have a payload (int) ID that would identify the proper
		decoding structure for the payload

		This is a WIP design, but for MVP it seems like the way to go
*/

// Payload is the 'meat' of the message
type Payload struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}

// Generalized message structure that describes the network, client, nonce, and payload
type Message struct {
	Network int     `json:"net"`
	Client  string  `json:"client"`
	Nonce   int64   `json:"nonce"`
	Payload Payload `json:"payload"`
}

type NewSessionPayload struct {
	DevID string                `json:"devid"` // the devID of the session
	Peers []session.SessionPeer `json:"peers"` // the list of peers
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
type ExitNetworkPayload struct {
	node.Node `json:"node"`
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
type EnterNetworkPayload struct {
	node.Node `json:"node"`
}

