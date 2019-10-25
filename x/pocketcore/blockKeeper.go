package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/types"
)

// blockKeeper handles access/modifiers of blocks
type blockKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// newBlockKeeper creates a new instance of block keeper
func newBlockKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) blockKeeper {
	return blockKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (bk blockKeeper) GetLatestBlockID(ctx sdk.Context) types.BlockID {
	// return fixtures.GenerateBlockHash()
	header := ctx.BlockHeader()
	return types.BlockID(header.GetLastBlockId())
}

func (bk blockKeeper) GetLatestSessionBlockID(ctx sdk.Context) types.BlockID {
	//return fixtures.GenerateBlockHash()
	latestsessionBlockHeight := bk.GetLatestSessionBlockHeight(ctx)
	ctxAtHeight := ctx.WithBlockHeight(latestsessionBlockHeight)
	return bk.GetLatestBlockID(ctxAtHeight)
}

func (bk blockKeeper) GetLatestSessionBlockHeight(ctx sdk.Context) int64 {
	//return fixtures.GenerateBlockHash()
	blkHeight := ctx.BlockHeight()
	return (blkHeight / SESSIONBLOCKFREQUENCY) * SESSIONBLOCKFREQUENCY
}
