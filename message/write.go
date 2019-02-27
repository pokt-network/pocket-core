package message

import (
	"time"

	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/message/fbs"
)

func MarshalMessage(builder *flatbuffers.Builder, message Message) []byte {
	// this line allows us to reuse the same builder
	builder.Reset()
	// Create a variable to hold the payload
	payloadVector := builder.CreateByteVector(message.Payload)
	// Create the message
	fbs.MessageStart(builder)
	fbs.MessageAddPayload(builder, payloadVector)
	fbs.MessageAddType(builder, message.Type_)
	fbs.MessageAddTimestamp(builder, uint32(time.Now().Unix()))
	m := fbs.MessageEnd(builder)
	builder.Finish(m)
	return builder.FinishedBytes()
}

func MarshalHelloMessage(builder *flatbuffers.Builder, helloMessage HelloMessage) []byte {
	// this line allows us to reuse the same builder
	builder.Reset()
	// Create a variable to hold the flatbuffer byte vector
	gidVector := builder.CreateByteVector([]byte(helloMessage.Gid))
	// Create the hello message
	fbs.HelloMessageStart(builder)
	fbs.HelloMessageAddGid(builder, gidVector)
	hm := fbs.HelloMessageEnd(builder)
	// since helloMessage is the root_object
	builder.Finish(hm)
	return builder.FinishedBytes()
}
