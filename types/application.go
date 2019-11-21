package types

import (
	"github.com/pokt-network/pocket-core/legacy"
)

// "AppPubKey" is the base structure for a Pocket Network AppPubKey
type Application struct {
	Account               `json:"routing"`
	LegacyRequestedChains []LegacyApplicationRequestedBlockchain `json:"legacyRequestedChains"` // todo phase out
	RequestedBlockchains  ApplicationRequestedBlockchains        `json:"requestedChains"`
}

type Applications []Application

func NewApplication(address Address, publicKey AccountPublicKey, balance POKT, stakeAmount POKT, requestedChains ApplicationRequestedBlockchains) Application {
	return Application{
		Account: Account{
			Address:     address,
			PubKey:      publicKey,
			Balance:     balance,
			StakeAmount: stakeAmount,
		},
		RequestedBlockchains: requestedChains,
	}
}

type LegacyApplicationRequestedBlockchain struct {
	legacy.Blockchain    `json:"blockchain"`
	AllocationPercentage uint8 `json:"allocationPercentage"`
}

type ApplicationRequestedBlockchains []ApplicationRequestedBlockchain

type ApplicationRequestedBlockchain struct {
	Blockchain           `json:"blockchain"`
	AllocationPercentage uint8 `json:"allocationPercentage"`
}

func (apr ApplicationRequestedBlockchains) Contains(blockchain Blockchain) bool {
	for _, val := range apr {
		if val.Blockchain == blockchain {
			return true
		}
	}
	return false
}
