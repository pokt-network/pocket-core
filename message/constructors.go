// This package is all message related code
package message

import (
	"github.com/pokt-network/pocket-core/const"
)

// "constructors.go" holds all of the message constructor code

func NewMessage(pLoad Payload) Message {
	return Message{_const.NETID, _const.CLIENTID, 0, pLoad}
}
// TODO all payload constructors go here
