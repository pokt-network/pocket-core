package service

import "github.com/pokt-network/pocket-core/types"

type ServiceBlockchain types.Blockchain

type ServiceBlockchains types.Blockchains

func (sbc ServiceBlockchain) String() string {
	return types.Blockchain(sbc).String()
}

func (sbcs ServiceBlockchains) GetChainURL(blockchainHex string) (string, error) {
	return types.Blockchains(sbcs).GetChainURL(blockchainHex)
}

func (sbcs ServiceBlockchains) Contains(blockchainHex string) bool {
	return types.Blockchains(sbcs).Contains(blockchainHex)
}
