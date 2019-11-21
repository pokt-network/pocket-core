package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	tndmt "github.com/tendermint/tendermint/abci/types"
)

func (k Keeper) GetLatestSessionBlock(ctx sdk.Context) (blockID tndmt.BlockID) {
	sessionBlockHeight := (ctx.BlockHeight() % int64(k.posKeeper.SessionBlock(ctx))) + 1
	return ctx.WithBlockHeight(sessionBlockHeight).BlockHeader().GetLastBlockId()
}
