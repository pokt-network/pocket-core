package message

import (
	"github.com/pokt-network/pocket-core/session"
)

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

type SessionPL struct {
	DevID string         `json:"devid"` // the devID of the session
	Peers []session.Peer `json:"peers"` // the list of peers
}
