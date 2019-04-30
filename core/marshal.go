package core

import (
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/core/fbs"
	"strconv"
)

// "MarshalBlockchain" converts a blockchain object into a serialized flatbuffer
func MarshalBlockchain(builder *flatbuffers.Builder, blockchain Blockchain) ([]byte, error) {
	netid, err := strconv.ParseUint(blockchain.NetID, 10, 32)
	if err != nil {
		return nil, err
	}
	vers, err := strconv.ParseUint(blockchain.Version, 10, 32)
	// this line allows us to reuse the same builder
	builder.Reset()
	// Create a variable to hold the blockchain name
	blockchainNameVector := builder.CreateByteVector([]byte(blockchain.Name))
	// Create the blockchain
	fbs.BlockchainStart(builder)
	fbs.BlockchainAddName(builder, blockchainNameVector)
	fbs.BlockchainAddNetid(builder, uint32(netid))
	fbs.BlockchainAddVersion(builder, uint32(vers))
	bc := fbs.BlockchainEnd(builder)
	builder.Finish(bc)
	return builder.Bytes, nil
}

func marshalToken(builder *flatbuffers.Builder, token Token) (flatbuffers.UOffsetT, error) {
	expdateVector := builder.CreateByteVector(token.ExpDate)
	fbs.TokenStart(builder)
	fbs.TokenAddExpdate(builder, expdateVector)
	return fbs.TokenEnd(builder), nil
}

func marshalRelay(builder *flatbuffers.Builder, relay Relay) (flatbuffers.UOffsetT, error) {
	var urlVector flatbuffers.UOffsetT
	builder.Reset()
	blockchainVector := builder.CreateByteVector(relay.Blockchain)
	payloadVector := builder.CreateByteVector(relay.Payload)
	devidVector := builder.CreateByteVector(relay.DevID)
	methodVector := builder.CreateByteVector(relay.Method)
	if relay.Path != nil && len(relay.Path) != 0 {
		urlVector = builder.CreateByteVector(relay.Path)
	}
	t, err := marshalToken(builder, relay.Token)
	if err != nil {
		return t, err
	}
	// create the relay
	fbs.RelayStart(builder)
	fbs.RelayAddPayload(builder, payloadVector)
	fbs.RelayAddDevid(builder, devidVector)
	fbs.RelayAddMethod(builder, methodVector)
	fbs.RelayAddToken(builder, t)
	fbs.RelayAddBlockchain(builder, blockchainVector)
	if relay.Path != nil && len(relay.Path) != 0 {
		fbs.RelayAddUrl(builder, urlVector)
	}
	return fbs.RelayEnd(builder), nil
}

// "MarshalRelay" converts a relay ojbect into a serialized flatbuffer
func MarshalRelay(builder *flatbuffers.Builder, relay Relay) ([]byte, error) { // TODO set empty http method as POST
	r, err := marshalRelay(builder, relay)
	if err != nil {
		return nil, err
	}
	builder.Finish(r)
	return builder.FinishedBytes(), nil
}

// "MarshalRelayMessage" converts a relay message into a serialized flatbuffer
func MarshalRelayMessage(builder *flatbuffers.Builder, relayMessage RelayMessage) ([]byte, error) {
	builder.Reset()
	relay, err := marshalRelay(builder, relayMessage.Relay)
	if err != nil {
		return nil, err
	}
	signatureVector := builder.CreateByteVector(relayMessage.Signature)
	fbs.RelayMessageStart(builder)
	fbs.RelayMessageAddRelay(builder, relay)
	fbs.RelayMessageAddSignature(builder, signatureVector)
	rm := fbs.RelayMessageEnd(builder)
	builder.Finish(rm)
	return builder.FinishedBytes(), nil
}
