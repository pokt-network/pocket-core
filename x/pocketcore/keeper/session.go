package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "HandleDispatch" - Handles a client request for their session information
func (k Keeper) HandleDispatch(ctx sdk.Ctx, header types.SessionHeader) (*types.DispatchResponse, sdk.Error) {
	fmt.Println("OLSH5")
	// retrieve the latest session block height
	latestSessionBlockHeight := k.GetLatestSessionBlockHeight(ctx)
	// set the session block height
	header.SessionBlockHeight = latestSessionBlockHeight
	// validate the header
	fmt.Println("OLSH6")
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	// get the session context
	fmt.Println("OLSH7")
	sessionCtx, er := ctx.PrevCtx(latestSessionBlockHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	fmt.Println("OLSH8")
	// check cache
	session, found := types.GetSession(header, types.GlobalSessionCache)
	// if not found generate the session
	if !found {
		fmt.Println("OLSH9")
		var err sdk.Error
		blockHashBz, er := sessionCtx.BlockHash(k.Cdc, sessionCtx.BlockHeight())
		if er != nil {
			return nil, sdk.ErrInternal(er.Error())
		}
		fmt.Println("OLSH10")
		session, err = types.NewSession(sessionCtx, ctx, k.posKeeper, header, hex.EncodeToString(blockHashBz), int(k.SessionNodeCount(sessionCtx)))
		if err != nil {
			return nil, err
		}
		fmt.Println("OLSH11")
		// add to cache
		types.SetSession(session, types.GlobalSessionCache)
	}
	fmt.Println("OLSH12")
	actualNodes := make([]exported.ValidatorI, len(session.SessionNodes))
	for i, addr := range session.SessionNodes {
		actualNodes[i], _ = k.GetNode(sessionCtx, addr)
	}
	fmt.Println("OLSH13")
	return &types.DispatchResponse{Session: types.DispatchSession{
		SessionHeader: session.SessionHeader,
		SessionKey:    session.SessionKey,
		SessionNodes:  actualNodes,
	}, BlockHeight: ctx.BlockHeight()}, nil
}

// "IsSessionBlock" - Returns true if current block, is a session block (beginning of a session)
func (k Keeper) IsSessionBlock(ctx sdk.Ctx) bool {
	return ctx.BlockHeight()%k.posKeeper.BlocksPerSession(ctx) == 1
}

// IsProofSessionHeightWithinTolerance checks if the relaySessionBlockHeight is bounded by (latestSessionBlockHeight - tolerance ) <= x <= latestSessionHeight
func (k Keeper) IsProofSessionHeightWithinTolerance(ctx sdk.Ctx, relaySessionBlockHeight int64) bool {
	// Session block height can never be zero.
	if relaySessionBlockHeight <= 0 {
		return false
	}
	latestSessionHeight := k.GetLatestSessionBlockHeight(ctx)
	tolerance := types.GlobalPocketConfig.ClientSessionSyncAllowance * k.posKeeper.BlocksPerSession(ctx)
	minHeight := latestSessionHeight - tolerance
	return sdk.IsBetween(relaySessionBlockHeight, minHeight, latestSessionHeight)
}

// "GetLatestSessionBlockHeight" - Returns the latest session block height (first block of the session, (see blocksPerSession))
func (k Keeper) GetLatestSessionBlockHeight(ctx sdk.Ctx) (sessionBlockHeight int64) {
	// get the latest block height
	blockHeight := ctx.BlockHeight()
	// get the blocks per session
	blocksPerSession := k.posKeeper.BlocksPerSession(ctx)
	// if block height / blocks per session remainder is zero, just subtract blocks per session and add 1
	if blockHeight%blocksPerSession == 0 {
		sessionBlockHeight = blockHeight - k.posKeeper.BlocksPerSession(ctx) + 1
	} else {
		// calculate the latest session block height by diving the current block height by the blocksPerSession
		sessionBlockHeight = (blockHeight/blocksPerSession)*blocksPerSession + 1
	}
	return
}

// "IsPocketSupportedBlockchain" - Returns true if network identifier param is supported by pocket
func (k Keeper) IsPocketSupportedBlockchain(ctx sdk.Ctx, chain string) bool {
	// loop through supported blockchains (network identifiers)
	for _, c := range k.SupportedBlockchains(ctx) {
		// if contains chain return true
		if c == chain {
			return true
		}
	}
	// else return false
	return false
}

func (Keeper) ClearSessionCache() {
	types.ClearSessionCache(types.GlobalSessionCache)
}
