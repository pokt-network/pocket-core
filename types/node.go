package types

import (
	"github.com/pokt-network/pocket-core/legacy"
)

// "Node" is the base structure for a Pocket Network Node"
type Node struct {
	Account         `json:"routing"`
	URL             []byte              `json:"url"`
	SupportedChains NodeSupportedChains `json:"supportedChains"`
	IsAlive         bool
}

type Nodes []Node

type NodeSupportedChain struct {
	legacy.Blockchain `json:"blockchain"`
}

type NodeSupportedChains map[string]NodeSupportedChain // [hex]->Blockchain

func (nsc NodeSupportedChains) Add(hexBlockchainhash string, blockchain NodeSupportedChain) {
	nsc[hexBlockchainhash] = blockchain
}

func (nsc NodeSupportedChains) Contains(hexBlockchainHash string) bool {
	_, contains := nsc[hexBlockchainHash]
	return contains
}
