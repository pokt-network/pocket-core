package message

import (
	"errors"
	
	"github.com/pokt-network/pocket-core/message/fbs"
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
