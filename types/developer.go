package types

import (
	"github.com/pokt-network/pocket-core/legacy"
)

// "Developer" is the base structure for a Pocket Network Developer
type Developer struct {
	Account         `json:"routing"`
	RequestedChains []DeveloperRequestedBlockchain `json:"requestedChains"`
}

type Developers []Developer

func NewDeveloper(address Address, publicKey AccountPublicKey, balance POKT, stakeAmount POKT, requestedChains []DeveloperRequestedBlockchain) Developer {
	return Developer{
		Account: Account{
			Address:     address,
			PubKey:      publicKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		RequestedChains: requestedChains,
	}
}

type DeveloperRequestedBlockchain struct {
	legacy.Blockchain    `json:"blockchain"`
	AllocationPercentage uint8 `json:"allocationPercentage"`
}
