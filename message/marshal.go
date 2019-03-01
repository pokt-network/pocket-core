package message

import (
	"time"

	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/message/fbs"
	"github.com/pokt-network/pocket-core/service"
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
	gidVector := builder.CreateByteVector([]byte(helloMessage.Gid))
	// Create the hello message
	fbs.HelloMessageStart(builder)
	fbs.HelloMessageAddGid(builder, gidVector)
	hm := fbs.HelloMessageEnd(builder)
	// since helloMessage is the root_object
	builder.Finish(hm)
	return builder.FinishedBytes()
}

func MarshalValidateMessage(builder *flatbuffers.Builder, validateMessage ValidateMessage) []byte {
	builder.Reset()
	// create the fbs strings for the relay object
	bc := builder.CreateString(validateMessage.Relay.Blockchain)
	devid := builder.CreateString(validateMessage.Relay.DevID)
	data := builder.CreateString(validateMessage.Relay.Data)
	ver := builder.CreateString(validateMessage.Relay.Version)
	nid := builder.CreateString(validateMessage.Relay.NetworkID)
	// start the relay build process
	fbs.RelayStart(builder)
	fbs.RelayAddBlockchain(builder, bc)
	fbs.RelayAddDevID(builder, devid)
	fbs.RelayAddData(builder, data)
	fbs.RelayAddVersion(builder, ver)
	fbs.RelayAddNetid(builder, nid)
	r := fbs.RelayEnd(builder)
	// create the fbs byte vector for the hash
	h := builder.CreateByteVector(validateMessage.Hash)
	// start the validate message
	fbs.ValidateMessageStart(builder)
	fbs.ValidateMessageAddHash(builder, h)
	fbs.ValidateMessageAddRelay(builder, r)
	vm := fbs.ValidateMessageEnd(builder)
	// finish the validate message
	builder.Finish(vm)
	return builder.FinishedBytes()
}

func MarshalHelloSession(builder *flatbuffers.Builder, helloSessionMessage HelloSessionMessage) []byte {
	builder.Reset()
	gidVector := builder.CreateByteVector([]byte(helloSessionMessage.Gid))
	fbs.HelloSessionMessageStart(builder)
	fbs.HelloSessionMessageAddGid(builder, gidVector)
	fbs.HelloSessionMessageAddRole(builder, helloSessionMessage.Role)
	hsm := fbs.HelloSessionMessageEnd(builder)
	builder.Finish(hsm)
	return builder.FinishedBytes()
}

func MarshalRelay(builder *flatbuffers.Builder, relay service.Relay) []byte {
	builder.Reset()
	bc := builder.CreateString(relay.Blockchain)
	devid := builder.CreateString(relay.DevID)
	data := builder.CreateString(relay.Data)
	ver := builder.CreateString(relay.Version)
	nid := builder.CreateString(relay.NetworkID)
	fbs.RelayStart(builder)
	fbs.RelayAddBlockchain(builder, bc)
	fbs.RelayAddDevID(builder, devid)
	fbs.RelayAddData(builder, data)
	fbs.RelayAddVersion(builder, ver)
	fbs.RelayAddNetid(builder, nid)
	r := fbs.RelayEnd(builder)
	builder.Finish(r)
	return builder.FinishedBytes()
}
