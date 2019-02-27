package message

import (
	"errors"
	
	"github.com/pokt-network/pocket-core/message/fbs"
)

func RouteMessageByPayload(m Message) (interface{}, error) {
	switch m.Type_ {
	case fbs.MessageTypeUNDEFINED:
		return m.Payload, errors.New("undefined message type")
	case fbs.MessageTypeDISC_HELLO:
		return UnmarshalHelloMessage(m.Payload)
	default:
		return m.Payload, errors.New("unsupported message type" + string(m.Type_))
	}
}

func UnmarshalMessage(flatBuffer []byte) Message {
	message := fbs.GetRootAsMessage(flatBuffer, 0)
	return Message{message.Type(), message.PayloadBytes(), message.Timestamp()}
}

func UnmarshalHelloMessage(flatBuffer []byte) (*HelloMessage, error) {
	helloMessage := fbs.GetRootAsHelloMessage(flatBuffer, 0)
	// TODO add more error checking on the GID
	if len(string(helloMessage.GidBytes())) == 0 {
		return nil, errors.New("unable to unmarshal to hello message")
	}
	return &HelloMessage{string(helloMessage.GidBytes())}, nil
}
