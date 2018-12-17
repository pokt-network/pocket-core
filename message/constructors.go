// This package is all message related code
package message

import "github.com/pokt-network/pocket-core/const"

// "constructors.go" holds all of the message constructor code

func NewMessage(pLoad Payload) Message{
	message := Message{}
	message.Network=_const.NETID
	message.Client = _const.CLIENTID
	message.Nonce = 0				// TODO implement nonce
	message.Payload = pLoad
}

