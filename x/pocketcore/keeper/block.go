package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetLatestSessionBlock(ctx sdk.Context) (hexBlockHash string) {
	sessionBlockHeight := (ctx.BlockHeight() % types.SESSIONBLOCKFREQUENCY) + 1
	return hex.EncodeToString(ctx.WithBlockHeight(sessionBlockHeight).BlockHeader().GetLastBlockId().Hash)
}
