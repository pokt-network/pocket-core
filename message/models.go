// This package is all message related code
package message

import "github.com/pokt-network/pocket-core/session"

// "models.go" holds all of the structures for the message package
// TODO proper type indexing for message payload
/*
NOTE: 	The ideology here in design is to maintain a simple message structure while swapping out the payload.

		For example: A message would have a payload (int) ID that would identify the proper
		decoding structure for the payload

		This is a WIP design, but for MVP it seems like the way to go
 */

type Payload struct {
	ID 			int				`json:"id"`
	Data		interface{}		`json:"data"`
}

type CreateSessPL struct {
	Session		session.Session	`json:"session"`
}

type Message struct {
	Network			int			`json:"net"`
	Client			string		`json:"client"`
	Nonce			int64		`json:"nonce"`
	Payload			Payload		`json:"payload"`
}


