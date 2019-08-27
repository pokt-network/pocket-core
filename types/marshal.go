package types

import (
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/types/fbs"
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
