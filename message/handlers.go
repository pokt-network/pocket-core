package message

import (
	"errors"
	
	"github.com/pokt-network/pocket-core/message/fbs"
	"github.com/pokt-network/pocket-core/service"
)

func RouteMessageByPayload(m Message) error {
	switch m.Type_ {
	case fbs.MessageTypeUNDEFINED:
		return nil
	case fbs.MessageTypeDISCHELLO:
		return HandleDISCHELLOMessage(m)
	case fbs.MessageTypeVALIDATE:
		return HandleValidateMessage(m)
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

func HandleValidateMessage(m Message) error { // TODO
	vm, err := UnmarshalValidateMessage(m.Payload)
	if err != nil {
		return err
	}
	b, err := service.Validate(vm.Relay, vm.Hash)
	if err != nil {
		return err
	}
	b = b // arbitrary
	return nil
}
