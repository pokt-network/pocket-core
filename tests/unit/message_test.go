package unit

import (
	"log"
	"testing"
	"time"
	
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/message"
	"github.com/pokt-network/pocket-core/message/fbs"
)

const gid = "TESTGID"

func TestMessageSerialization(t *testing.T) {
	// use a flatbuffers builder
	builder := flatbuffers.NewBuilder(0)
	// serialize a hello message struct into a flat buffer byte array
	hmBytes := message.MarshalHelloMessage(builder, message.HelloMessage{Gid: gid})
	// serialize a message struct into a flat buffer byte array with the hmBytes as the payload
	szMessage := message.MarshalMessage(builder, message.Message{Type_: fbs.MessageTypeDISC_HELLO, Payload: hmBytes})
	// proceed to unmarshal the message into a struct
	m := message.UnmarshalMessage(szMessage)
	// unmarshal the payload
	helloMessage := message.UnmarshalHelloMessage(m.Payload)
	if m.Type_ != fbs.MessageTypeDISC_HELLO || helloMessage.Gid != gid {
		t.Fatalf("Incorrect response from the serialization")
	}
	t.Log("\nThe message received was of type:", m.Type_, "\nThe message payload was:", helloMessage, "\nThe message was timestamped at ", time.Unix(int64(m.Timestamp), 0).UTC())
}

func TestHelloMessageSerialization(t *testing.T) {
	// Create a fbs builder to create our hello message flatbuffer
	builder := flatbuffers.NewBuilder(0)
	serializedResult := message.MarshalHelloMessage(builder, message.HelloMessage{Gid: gid})
	deserializedResult := message.UnmarshalHelloMessage(serializedResult)
	if deserializedResult.Gid != gid {
		log.Fatalf("Incorrect response from the serialization")
	}
	t.Log(deserializedResult)
	t.Log(serializedResult)
}
