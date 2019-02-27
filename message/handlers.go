package message

import (
	"errors"
	
	"github.com/pokt-network/pocket-core/message/fbs"
)

func RouteMessageByPayload(m Message) error {
	switch m.Type_ {
	case fbs.MessageTypeUNDEFINED:
		return nil
	case fbs.MessageTypeDISCHELLO:
		return HandleDISCHELLOMessage(m)
	default:
		return errors.New("unsupported message type" + string(m.Type_))
	}
}

func HandleUndefinedMessage(m Message) string {
	return string(m.Payload)
}

func HandleDISCHELLOMessage(m Message) error { // TODO
	helloMessage, err := UnmarshalHelloMessage(m.Payload)
	helloMessage = helloMessage // arbitrary
	if err != nil {
		return err
	}
	return nil
}
