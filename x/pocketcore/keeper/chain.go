package keeper

import (
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "GetHostedBlockchains" returns the non native chains hosted locally on this node
func (k Keeper) GetHostedBlockchains() *pc.HostedBlockchains {
	return k.hostedBlockchains
}
