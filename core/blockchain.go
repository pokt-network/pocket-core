// The core package implements the 'core' functionality of the Pocket Network (routing, servicing and validating relays)
package core

import "github.com/google/flatbuffers/go"

// TODO may need to add other parameters like chainID
// see ethereum vs eth classic
type Blockchain struct {
	Name    string `json:"name"`
	NetID   string `json:"netid"`
	Version string `json:"string"`
}

// "GenerateChainHash" takes a blockchain object and converts it to the protocol specific
// chain hash (encoded by google's flatbuffers)
func GenerateChainHash(bchain Blockchain) ([]byte, error) {
	b, err := MarshalBlockchain(flatbuffers.NewBuilder(0), bchain)
	if err != nil {
		return nil, err
	}
	return SHA3FromBytes(b), nil
}
