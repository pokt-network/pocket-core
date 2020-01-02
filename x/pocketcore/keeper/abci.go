package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k Keeper) {
	// delete the proofs held within the world state for too long
	k.DeleteExpiredClaims(ctx)
}
