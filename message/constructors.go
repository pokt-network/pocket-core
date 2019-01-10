// This package is all message related code
package message

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
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


func NewSessPayload(devID string, peers []session.SessionPeer) NewSessionPayload{
	return NewSessionPayload{devID,peers}
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func NewEnterNetPayload() EnterNetworkPayload{
	return EnterNetworkPayload{*node.GetSelf()}
}
//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func NewExitNetPayload() ExitNetworkPayload{
	return ExitNetworkPayload{*node.GetSelf()}
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func NewEnterNetMessage(payload EnterNetworkPayload) Message {
	return NewMessage(Payload{ID:2, Data: payload})
}

//NOTE: this is for pocket core mvp centralized dispatcher
// may remove for production
func NewExitNetMessage(payload ExitNetworkPayload) Message {
	return NewMessage(Payload{ID:3, Data: payload})
}
