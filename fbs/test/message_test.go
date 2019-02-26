package test

import (
	"log"
	"testing"
	"time"
	
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/fbs/messages"
)

const gid = "TESTGID"

func TestMessageSerialization(t *testing.T) {
	builder := flatbuffers.NewBuilder(0)
	szHelloMessage := SerializeHelloMessage(builder, gid)
	szMessage := SerializeMessage(builder, messages.TypeDISC_HELLO, szHelloMessage)
	timestamp, tp, payload := DeSerialzeMessage(szMessage)
	dzHelloMessage := DeSerializeHelloMessage(payload)
	if tp != messages.TypeDISC_HELLO || dzHelloMessage != gid {
		t.Fatalf("Incorrect response from the serialization")
	}
	t.Log("\nThe message received was of type:", tp, "\nThe message payload was:", dzHelloMessage, "\nThe message was timestamped at ", time.Unix(int64(timestamp), 0).UTC())
}

func TestHelloMessageSerialization(t *testing.T) {
	// Create a fbs builder to create our hello message flatbuffer
	builder := flatbuffers.NewBuilder(0)
	serializedResult := SerializeHelloMessage(builder, gid)
	deserializedResult := DeSerializeHelloMessage(serializedResult)
	if deserializedResult != gid {
		log.Fatalf("Incorrect response from the serialization")
	}
	t.Log(deserializedResult)
	t.Log(serializedResult)
}

func SerializeHelloMessage(builder *flatbuffers.Builder, gid string) []byte {
	// this line allows us to reuse the same builder
	builder.Reset()
	// Create a variable to hold the flatbuffer byte vector
	gidVector := builder.CreateByteVector([]byte(gid))
	// Create the hello message
	messages.HelloMessageStart(builder)
	messages.HelloMessageAddGid(builder, gidVector)
	helloMessage := messages.HelloMessageEnd(builder)
	// since helloMessage is the root_object
	builder.Finish(helloMessage)
	return builder.FinishedBytes()
}

func DeSerializeHelloMessage(flatBuffer []byte) string {
	helloMessage := messages.GetRootAsHelloMessage(flatBuffer, 0)
	return string(helloMessage.GidBytes())
}

func SerializeMessage(builder *flatbuffers.Builder, t messages.Type, p []byte) []byte {
	// this line allows us to reuse the same builder
	builder.Reset()
	// Create a variable to hold the payload
	payloadVector := builder.CreateByteVector(p)
	// Create the message
	messages.MessageStart(builder)
	messages.MessageAddPayload(builder, payloadVector)
	messages.MessageAddType(builder, t)
	messages.MessageAddTimestamp(builder, uint32(time.Now().Unix()))
	message := messages.MessageEnd(builder)
	builder.Finish(message)
	return builder.FinishedBytes()
}

func DeSerialzeMessage(flatBuffer []byte) (timestamp uint32, t messages.Type, payload []byte) {
	message := messages.GetRootAsMessage(flatBuffer, 0)
	timestamp = message.Timestamp()
	t = message.Type()
	payload = message.PayloadBytes()
	return
}
