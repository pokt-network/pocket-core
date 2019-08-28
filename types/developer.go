package model

import (
	"github.com/pokt-network/pocket-core/legacy"
)

type Developer struct {
	Account         `json:"routing"`
	RequestedChains []DeveloperRequestedBlockchain `json:"requestedChains"`
}

type DeveloperRequestedBlockchain struct {
	legacy.Blockchain    `json:"blockchain"`
	AllocationPercentage uint8 `json:"allocationPercentage"`
}
