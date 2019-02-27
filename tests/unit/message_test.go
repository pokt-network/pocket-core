package unit

import (
	"fmt"
	"log"
	"testing"
	"time"
	
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/message/fbs"
	"github.com/pokt-network/pocket-core/service"
)

const gid = "TESTGID"

func TestMessageSerialization(t *testing.T) {
	// use a flatbuffers builder
	builder := flatbuffers.NewBuilder(0)
	// serialize a hello message struct into a flat buffer byte array
	hmBytes := message.MarshalHelloMessage(builder, message.HelloMessage{Gid: gid})
	// serialize a message struct into a flat buffer byte array with the hmBytes as the payload
	szMessage := message.MarshalMessage(builder, message.Message{Type_: fbs.MessageTypeDISCHELLO, Payload: hmBytes})
	// proceed to unmarshal the message into a struct
	m := message.UnmarshalMessage(szMessage)
	// unmarshal the payload
	helloMessage, err := message.UnmarshalHelloMessage(m.Payload)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if m.Type_ != fbs.MessageTypeDISCHELLO || helloMessage.Gid != gid {
		t.Fatalf("Incorrect response from the serialization")
	}
	t.Log("\nThe message received was of type:", m.Type_, "\nThe message payload was:", helloMessage, "\nThe message was timestamped at ", time.Unix(int64(m.Timestamp), 0).UTC())
}

func TestHelloMessageSerialization(t *testing.T) {
	// Create a fbs builder to create our hello message flatbuffer
	builder := flatbuffers.NewBuilder(0)
	serializedResult := message.MarshalHelloMessage(builder, message.HelloMessage{Gid: gid})
	deserializedResult, err := message.UnmarshalHelloMessage(serializedResult)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if deserializedResult.Gid != gid {
		log.Fatalf("Incorrect response from the serialization")
	}
	t.Log(deserializedResult)
	t.Log(serializedResult)
}

func TestValidateMessageSerialization(t *testing.T) {
	// define validate message and relay struct
	const (
		dummyRelayAnswer = "dummy"
		blockchain       = "ethereum"
		version          = "0"
		netid            = "0"
		data             = "dummy2"
		devid            = "devid1"
	)
	hash := service.ValidationHash(dummyRelayAnswer)
	r := service.Relay{Blockchain: blockchain, NetworkID: netid, Version: version, Data: data, DevID: devid}
	vm := message.ValidateMessage{Relay: r, Hash: hash}
	// create a fbs buffer to create our validate message from
	builder := flatbuffers.NewBuilder(0)
	// serialize
	vmBytes := message.MarshalValidateMessage(builder, vm)
	// deserialize
	validateMessage := message.UnmarshalValidateMessage(vmBytes)
	fmt.Println(validateMessage)
}
