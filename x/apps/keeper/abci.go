package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker - Called at the beggining of every block
func BeginBlocker(ctx sdk.Ctx, _ abci.RequestBeginBlock, k Keeper) {
	// burn applications triggered by the custom burning interface
	k.burnApplications(ctx)
}

// EndBlocker - Called at the end of every block, update validator set
func EndBlocker(ctx sdk.Ctx, k Keeper) []abci.ValidatorUpdate {
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
