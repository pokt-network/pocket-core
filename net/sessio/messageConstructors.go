package sessio

import "github.com/pokt-network/pocket-core/message"

func NewSessionMessage(session Session) message.Message{
	return message.NewMessage(message.Payload{ID: 1, Data: session})
}
