package service

import (
	"github.com/pokt-network/pocket-core/types"
)

type ServiceBlockchain types.Blockchain

type ServiceBlockchains types.HostedBlockchains

func (s ServiceBlockchains) Contains(hash string) bool {
	hbc := types.HostedBlockchains(s)
	return hbc.ContainsFromString(hash)
}

func (s ServiceBlockchain) GetHostedChainURL(hostChains ServiceBlockchains) (string, error) {
	hc := types.HostedBlockchains(hostChains)
	if hc.Len() == 0 {
		return "", EmptyHostedChainsError
	}
	return hc.GetChainFromBytes(s).URL, nil
}
