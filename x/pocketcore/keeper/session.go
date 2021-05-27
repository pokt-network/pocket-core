package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"os"
	"sync"
	"time"
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
	session, found := types.GetSession(header)
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
		types.SetSession(session)
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
	types.ClearSessionCache()
}

func (k Keeper) InitSessionValidators(ctx sdk.Ctx) sdk.Error {
	ctx.Logger().Info("Started initializing session validators")
	start := time.Now()
	blocksBack := int64(25) // TODO get from config
	blocksPerSession := k.BlocksPerSession(ctx)
	sessionBlockHeight := k.GetLatestSessionBlockHeight(ctx)
	// get the session height from x blocks ago as the starting height
	height := sessionBlockHeight - blocksBack // (genesis block 0 is unable to be used)
	height = (height/blocksPerSession)*blocksPerSession + 1
	if height < 1 {
		height = 1
	}
	// create the global structure
	types.GlobalSessionVals = &types.SessionValidators{
		M: make(map[int64]map[string]exported.ValidatorI),
		S: make(map[int64]map[string][]types.SessionValidator),
		L: &sync.Mutex{},
	}
	for ; height <= sessionBlockHeight; height += blocksPerSession {
		var endHeight int64
		ctx.Logger().Info(fmt.Sprintf("Initializing session Validators for height: %d", height))
		sessionCtx, err := ctx.PrevCtx(height)
		if err != nil {
			return sdk.ErrInternal(errors.Wrap(err, "an error occurred getting prevctx in InitSessionValidators").Error())
		}
		sessionVals := k.posKeeper.AllValidators(sessionCtx)
		if ctx.BlockHeight()-sessionCtx.BlockHeight() < blocksPerSession {
			endHeight = ctx.BlockHeight()
		} else {
			endHeight = height + blocksPerSession - 1
		}
		endSessionCtx, err := ctx.PrevCtx(endHeight)
		if err != nil {
			ctx.Logger().Error("ERROR getting end session ctx from init session validators")
			os.Exit(1)
		}
		endSessionVals := k.posKeeper.AllValidators(endSessionCtx)
		types.InitSessionValidators(height, sessionVals, endSessionVals)
	}
	ctx.Logger().Debug("Finished initializing session validators: ", time.Since(start))
	return nil
}

func (k Keeper) UpdateSessionValidators(ctx sdk.Ctx) sdk.Error {
	ctx.Logger().Info("Started updating session validators")
	start := time.Now()
	blocksBack := int64(25) // TODO get from config
	blocksPerSession := k.BlocksPerSession(ctx)
	err := types.UpdateSessionValidators(ctx.BlockHeight()+1, blocksBack, blocksPerSession)
	if err != nil {
		return sdk.ErrInternal(errors.Wrap(err, "unable to update session validators").Error())
	}
	ctx.Logger().Debug(fmt.Sprintf("Updating Session Validators Took: %s", time.Since(start)))
	return nil
}

func (k Keeper) UpdateSessionValidator(ctx sdk.Ctx, val exported.ValidatorI) {
	start := time.Now()
	setS := false
	if k.IsSessionBlock(ctx) {
		setS = true
	}
	types.SetSessionValidator(k.GetLatestSessionBlockHeight(ctx), val, setS)
	ctx.Logger().Debug(fmt.Sprintf("Update Session Validator Took: %s", time.Since(start)))
}
