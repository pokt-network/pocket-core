package keeper

import (
	"encoding/hex"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "HandleDispatch" - Handles a client request for their session information
func (k Keeper) HandleDispatch(ctx sdk.Ctx, header types.SessionHeader) (*types.DispatchResponse, sdk.Error) {
	// retrieve the latest session block height
	latestSessionBlockHeight := k.GetLatestSessionBlockHeight(ctx)
	// set the session block height
	header.SessionBlockHeight = latestSessionBlockHeight
	// validate the header
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	// get the session context
	sessionCtx, er := ctx.PrevCtx(latestSessionBlockHeight)
	if er != nil {
		return nil, sdk.ErrInternal(er.Error())
	}
	// check cache
	session, found := types.GetSession(header, types.GlobalSessionCache)
	// if not found generate the session
	if !found {
		var err sdk.Error
		blockHashBz, er := sessionCtx.BlockHash(k.Cdc, sessionCtx.BlockHeight())
		if er != nil {
			return nil, sdk.ErrInternal(er.Error())
		}
		session, err = types.NewSession(sessionCtx, ctx, k.posKeeper, header, hex.EncodeToString(blockHashBz), int(k.SessionNodeCount(sessionCtx)))
		if err != nil {
			return nil, err
		}
		// add to cache
		types.SetSession(session, types.GlobalSessionCache)
	}
	actualNodes := make([]exported.ValidatorI, len(session.SessionNodes))
	for i, addr := range session.SessionNodes {
		actualNodes[i], _ = k.GetNode(sessionCtx, addr)
	}
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

	// DISC #1: IsProofSessionHeightWithinTolerance pocket-core is different from the spec
	// relaySessionBlockHeight = 11
	// Assumption: portal is synched with the network (session = 11)
	// Assumption: geo-mesh is not synched with the network (session = 7)
	// Expectation: succeed
	// Actual: failure
	// Enable: enable sessions in the past
	// Enable: enable

	// Goal:
	// Portal is at session 7
	// Geo-mesh is at session 11
	// you can ONLY handle sessions in the present

	// But, what if I didn't send my claim yet?
	// AND, I'm (a node) am ahead of portal (session 11)
	// AND, portal (session 7) sends me a relay

	// Question: why should som

	// Portal -> session = 15
	// Geo-mesh -> session = 11
	// Portal sees latest of the network:
	// 	- Good devops
	//  - good network
	// - good QoS
	// Geomesh:
	// 	- bad devops
	// 	- bad internet
	// 	- bad Qos => no reward

https: //docs.pokt.network/learn/protocol/servicing/#:~:text=Zero%20Knowledge%20Range%20Proof,-In%20order%20to&text=Generate%20the%20Merkle%20Tree%20using,leafs%20possible%20to%20select%20from


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
