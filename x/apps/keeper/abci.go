package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Ctx, _ abci.RequestBeginBlock, k Keeper) {}

// EndBlocker - Called at the end of every block, update validator set
func EndBlocker(ctx sdk.Ctx, k Keeper) []abci.ValidatorUpdate {
	// Unstake all mature applications from the unstakeing queue.
	k.unstakeAllMatureApplications(ctx)
	return []abci.ValidatorUpdate{}
}
