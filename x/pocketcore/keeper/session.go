package keeper

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) Dispatch(ctx sdk.Context, header types.SessionHeader) (*types.Session, sdk.Error) {
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	sessionCtx := ctx.MustGetPrevCtx(header.SessionBlockHeight)
	sessionBlkHeader := sessionCtx.BlockHeader()
	return types.NewSession(header.ApplicationPubKey, header.Chain, hex.EncodeToString(sessionBlkHeader.LastBlockId.Hash),
		header.SessionBlockHeight, k.GetAllNodes(ctx), int(k.SessionNodeCount(ctx)))
}

// is the context block a session block?
func (k Keeper) IsSessionBlock(ctx sdk.Context) bool {
	return ctx.BlockHeight()%k.posKeeper.SessionBlockFrequency(ctx) == 1
}

// get the most recent session block from the cont
func (k Keeper) GetLatestSessionBlock(ctx sdk.Context) sdk.Context {
	var sessionBlockHeight int64
	if ctx.BlockHeight()%k.posKeeper.SessionBlockFrequency(ctx) == 0 {
		sessionBlockHeight = ctx.BlockHeight() - k.posKeeper.SessionBlockFrequency(ctx) + 1
	} else {
		sessionBlockHeight = (ctx.BlockHeight()/k.posKeeper.SessionBlockFrequency(ctx))*k.posKeeper.SessionBlockFrequency(ctx) + 1
	}
	return ctx.MustGetPrevCtx(sessionBlockHeight)
}

// is the blockchain supported at this specific context?
func (k Keeper) IsPocketSupportedBlockchain(ctx sdk.Context, chain string) bool {
	for _, c := range k.SupportedBlockchains(ctx) {
		if c == chain {
			return true
		}
	}
	return false
}
