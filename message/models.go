// This package is all message related code
package message

// "models.go" holds all of the structures for the message package
// TODO proper type indexing for message payload
/*
NOTE: 	The ideology here in design is to maintain a simple message structure while swapping out the payload.

		For example: A dispatchSession Message would have a payload (int) ID that would identify the proper
		decoding structure for the payload (in this case is a JSON []byte that can be unmarshalled into SessionPeers //TODO create sessionPeers

		This is a WIP design, but for MVP it seems like the way to go
 */

type Payload struct {
	ID 			int				`json:"id"`
	Data			[]byte		`json:"data"`
}

type DispatchSession struct {
	SessionPeersJSON	[]byte	`json:"sessionPeers"`
}

type Message struct {
	Network			int			`json:"net"`
	Client			string		`json:"client"`
	Nonce			int64		`json:"nonce"`
	Payload			Payload		`json:"payload"`
}

