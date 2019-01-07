// This package is all message related code
package message

import (
	"github.com/pokt-network/pocket-core/const"
)

// "constructors.go" holds all of the message and payload constructor code

func NewMessage(pLoad Payload) Message {
	return Message{_const.NETID, _const.CLIENTID, 0, pLoad}
}

/*
"NewSessionMessage" creates a new message whose payload is info that can be used to derive a session
 */
func NewSessionMessage(nSPL NewSessionPayload) Message {
	return NewMessage(Payload{ID: 1, Data: nSPL})	// return message with NewSessionPayload
}
