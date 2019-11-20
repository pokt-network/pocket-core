package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"
)

func (k Keeper) GetLatestBlockID(ctx sdk.Context) abci.BlockID {
	header := ctx.BlockHeader()
	return header.GetLastBlockId()
}

func (k Keeper) GetBlockHeight(ctx sdk.Context) int64 {
	return ctx.BlockHeight()
}

func (k Keeper) GetBlockTime(ctx sdk.Context) time.Time {
	return ctx.BlockTime()
}
