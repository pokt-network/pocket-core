package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// burn any custom application slashes
	k.burnApplications(ctx)
}

// Called every block, update application set
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	matureApplications := k.getMatureApplications(ctx)
	// Unstake all mature applications from the unstakeing queue.
	k.unstakeAllMatureApplications(ctx)
	for _, valAddr := range matureApplications {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnstaking,
				sdk.NewAttribute(types.AttributeKeyApplication, valAddr.String()),
			),
		)
	}
	return []abci.ValidatorUpdate{}
}
