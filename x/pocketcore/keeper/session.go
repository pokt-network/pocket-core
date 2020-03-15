package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) Dispatch(ctx sdk.Ctx, header types.SessionHeader) (*types.DispatchResponse, sdk.Error) {
	latestSessionBlockHeight := k.GetLatestSessionBlockHeight(ctx)
	header.SessionBlockHeight = latestSessionBlockHeight
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	sessionCtx, er := ctx.PrevCtx(header.SessionBlockHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	// check cache
	session, found := types.GetSession(header)
	// if not found generate the session
	if !found {
		var err sdk.Error
		nodes := k.GetAllNodes(ctx)
		session, err = types.NewSession(header, types.BlockHash(sessionCtx), nodes, int(k.SessionNodeCount(sessionCtx)))
		if err != nil {
			return nil, err
		}
		// add to cache
		types.SetSession(session)
	}
	return &types.DispatchResponse{Session: session, BlockHeight: ctx.BlockHeight()}, nil
}

// is the context block a session block?
func (k Keeper) IsSessionBlock(ctx sdk.Ctx) bool {
	return ctx.BlockHeight()%k.posKeeper.SessionBlockFrequency(ctx) == 1
}

// get the most recent session block from the cont
func (k Keeper) GetLatestSessionBlockHeight(ctx sdk.Ctx) int64 {
	var sessionBlockHeight int64
	blockHeight := ctx.BlockHeight()
	frequency := k.posKeeper.SessionBlockFrequency(ctx)
	if blockHeight%frequency == 0 {
		sessionBlockHeight = ctx.BlockHeight() - k.posKeeper.SessionBlockFrequency(ctx) + 1
	} else {
		sessionBlockHeight = (ctx.BlockHeight()/k.posKeeper.SessionBlockFrequency(ctx))*k.posKeeper.SessionBlockFrequency(ctx) + 1
	}
	return sessionBlockHeight
}

// is the blockchain supported at this specific context?
func (k Keeper) IsPocketSupportedBlockchain(ctx sdk.Ctx, chain string) bool {
	for _, c := range k.SupportedBlockchains(ctx) {
		if c == chain {
			return true
		}
	}
	return false
}
