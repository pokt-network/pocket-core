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

func UnmarshalValidateMessage(flatBuffer []byte) ValidateMessage {
	validateMessage := fbs.GetRootAsValidateMessage(flatBuffer, 0)
	fbsRelay := &fbs.Relay{}
	validateMessage.Relay(fbsRelay)
	return ValidateMessage{ConvertFBSRelay(fbsRelay), validateMessage.HashBytes()}
}

func ConvertFBSRelay(fbsRelay *fbs.Relay) service.Relay {
	return service.Relay{
		Blockchain: string(fbsRelay.Blockchain()),
		NetworkID:  string(fbsRelay.NetworkID()),
		Version:    string(fbsRelay.Version()),
		Data:       string(fbsRelay.Data()),
		DevID:      string(fbsRelay.DevID())}
}
