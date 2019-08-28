package model

import (
	"github.com/pokt-network/pocket-core/legacy"
)

type Node struct {
	Account         `json:"routing"`
	URL             []byte               `json:"url"`
	SupportedChains []NodeSupportedChain `json:"supportedChains"`
	IsAlive         bool
}

type NodeSupportedChain struct {
	legacy.Blockchain `json:"blockchain"`
}
