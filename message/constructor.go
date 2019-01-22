package message

import (
	"github.com/pokt-network/pocket-core/const"
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
