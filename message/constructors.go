// This package is all message related code
package message

import (
	"github.com/pokt-network/pocket-core/const"
)

// "constructors.go" holds all of the message and payload constructor code

func NewMessage(pLoad Payload) Message {
	return Message{_const.NETID, _const.CLIENTID, 0, pLoad}
}

func NewSessionPeersMessage(sessionPeersJSON []byte) Message{
	// TODO
}
