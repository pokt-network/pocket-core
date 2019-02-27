package message

import (
	"github.com/pokt-network/pocket-core/message/fbs"
	"github.com/pokt-network/pocket-core/service"
)

type Message struct {
	Type_     fbs.MessageType `json:"type"`
	Payload   []byte          `json:"payload"`
	Timestamp uint32          `json:"timestamp"`
}

type HelloMessage struct {
	Gid string `json:"gid"`
}

type ValidateMessage struct {
	Relay service.Relay `json:"relay"`
	Hash  []byte        `json:"hash"`
}
