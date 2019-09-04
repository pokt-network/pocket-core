package types

import (
	"github.com/pokt-network/pocket-core/legacy"
)

// "Application" is the base structure for a Pocket Network Application
type Application struct {
	Account         `json:"routing"`
	RequestedChains []ApplicationRequestedBlockchain `json:"requestedChains"`
}

type Applications []Application

func NewApplication(address Address, publicKey AccountPublicKey, balance POKT, stakeAmount POKT, requestedChains []ApplicationRequestedBlockchain) Application {
	return Application{
		Account: Account{
			Address:     address,
			PubKey:      publicKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		RequestedChains: requestedChains,
	}
}

type ApplicationRequestedBlockchain struct {
	legacy.Blockchain    `json:"blockchain"`
	AllocationPercentage uint8 `json:"allocationPercentage"`
}
