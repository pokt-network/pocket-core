package keeper

import (
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "GetHostedBlockchains" returns the non native chains hosted locally on this node
func (k Keeper) GetHostedBlockchains() *pc.HostedBlockchains {
	return k.hostedBlockchains
}

func (k Keeper) SetHostedBlockchains(m map[string]pc.HostedBlockchain) *pc.HostedBlockchains {
	k.hostedBlockchains.L.Lock()
	k.hostedBlockchains.M = m
	k.hostedBlockchains.L.Unlock()
	return k.hostedBlockchains
}
