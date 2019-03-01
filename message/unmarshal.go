package message

import (
	"errors"

	"github.com/pokt-network/pocket-core/message/fbs"
	"github.com/pokt-network/pocket-core/service"
)

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

func UnmarshalRelay(flatbuffer []byte) (*service.Relay, error) {
	relay := fbs.GetRootAsRelay(flatbuffer, 0)
	return &service.Relay{Blockchain: string(relay.Blockchain()), NetworkID: string(relay.Netid()), Version: string(relay.Version()), Data: string(relay.Data()), DevID: string(relay.DevID())}, nil
}

func UnmarshalValidateMessage(flatbuffer []byte) (*ValidateMessage, error) {
	vm := fbs.GetRootAsValidateMessage(flatbuffer, 0)
	r := &fbs.Relay{}
	vm.Relay(r)
	return &ValidateMessage{
		Relay: service.Relay{
			Blockchain: string(r.Blockchain()),
			NetworkID:  string(r.Netid()),
			Version:    string(r.Version()),
			Data:       string(r.Data()),
			DevID:      string(r.DevID())},
		Hash: vm.HashBytes()}, nil
}
