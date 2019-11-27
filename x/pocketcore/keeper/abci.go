package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k Keeper) {
	// set new developer relays coefficient
	k.appKeeper.SetRelayCoefficient(ctx, int(k.GetStakedRatio(ctx).MulInt(sdk.NewInt(100)).TruncateInt().Int64()))
}
