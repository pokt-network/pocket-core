package blockchain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/types"
)

// BlockKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type BlockKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the BlockKeeper
func NewBlockKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) BlockKeeper {
	return BlockKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (bk BlockKeeper) GetLatestBlockID(ctx sdk.Context) types.BlockID {
	// return fixtures.GenerateBlockHash()
	header := ctx.BlockHeader()
	return types.BlockID(header.GetLastBlockId())
}

func (bk BlockKeeper) GetLatestSessionBlockID(ctx sdk.Context) types.BlockID {
	//return fixtures.GenerateBlockHash()
	latestsessionBlockHeight := bk.GetLatestSessionBlockHeight(ctx)
	ctxAtHeight := ctx.WithBlockHeight(latestsessionBlockHeight)
	return bk.GetLatestBlockID(ctxAtHeight)
}

func (bk BlockKeeper) GetLatestSessionBlockHeight(ctx sdk.Context) int64 {
	//return fixtures.GenerateBlockHash()
	blkHeight := ctx.BlockHeight()
	return (blkHeight / SESSIONBLOCKFREQUENCY) * SESSIONBLOCKFREQUENCY
}
