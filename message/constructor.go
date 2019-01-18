package message

import (
	"github.com/pokt-network/pocket-core/const"
	"github.com/pokt-network/pocket-core/node"
	"github.com/pokt-network/pocket-core/session"
)

// "NewMessage" creates a message structure from a payload structure
func NewMessage(pLoad Payload) Message {
	return Message{_const.NETID, _const.CLIENTID, 0, pLoad}
}

// "NewSessionMessage" creates a new message with a payload that can derive a session
func NewSessionMessage(devID string, peers []session.SessionPeer) Message {
	return NewMessage(Payload{ID: 1, Data: SessionPL{devID, peers}})
}

// DISCLAIMER: the code below is for pocket core mvp centralized dispatcher
// may remove for production

// "NewEnterMSG" creates a new message with a payload to enter the Pocket Network
func NewEnterMSG() Message {
	return NewMessage(Payload{ID: 2, Data: EnterPL{*node.GetSelf()}})
}

// "NewExitMSG" creates a new message with a payload to exit the Pocket Network.
func NewExitMSG() Message {
	return NewMessage(Payload{ID: 3, Data: ExitPL{*node.GetSelf()}})
}
