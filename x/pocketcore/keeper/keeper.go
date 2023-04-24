package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	authKeeper        types.AuthKeeper
	posKeeper         types.PosKeeper
	appKeeper         types.AppsKeeper
	TmNode            client.Client
	hostedBlockchains *types.HostedBlockchains
	Paramstore        sdk.Subspace
	storeKey          sdk.StoreKey // Unexposed key to access store from sdk.Context
	Cdc               *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the pocketcore module Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, authKeeper types.AuthKeeper, posKeeper types.PosKeeper, appKeeper types.AppsKeeper, hostedChains *types.HostedBlockchains, paramstore sdk.Subspace) Keeper {
	return Keeper{
		authKeeper:        authKeeper,
		posKeeper:         posKeeper,
		appKeeper:         appKeeper,
		hostedBlockchains: hostedChains,
		Paramstore:        paramstore.WithKeyTable(ParamKeyTable()),
		storeKey:          storeKey,
		Cdc:               cdc,
	}
}

func (k Keeper) Codec() *codec.Codec {
	return k.Cdc
}

// "GetBlock" returns the block from the tendermint node at a certain height
func (k Keeper) GetBlock(height int) (*coretypes.ResultBlock, error) {
	h := int64(height)
	return k.TmNode.Block(&h)
}

func (k Keeper) UpgradeCodec(ctx sdk.Ctx) {
	if ctx.IsOnUpgradeHeight() {
		k.ConvertState(ctx)
	}
}

func (k Keeper) ConvertState(ctx sdk.Ctx) {
	k.Cdc.SetUpgradeOverride(false)
	params := k.GetParams(ctx)
	claims := k.GetAllClaims(ctx)
	k.Cdc.SetUpgradeOverride(true)
	k.SetParams(ctx, params)
	k.SetClaims(ctx, claims)
	k.Cdc.DisableUpgradeOverride()
}

func (k Keeper) ConsensusParamUpdate(ctx sdk.Ctx) *abci.ConsensusParams {
	return k.consensusBlockSizeParamUpdate(ctx)
}

func (k Keeper) consensusBlockSizeParamUpdate(ctx sdk.Ctx) *abci.ConsensusParams {
	previousBlockCtx, err := ctx.PrevCtx(ctx.BlockHeight() - 1)
	if err != nil {
		ctx.Logger().Error("failed to get previous block context")
		return &abci.ConsensusParams{}
	}

	currentHeightBlockSize := k.BlockByteSize(ctx)
	// If it's 0, we're using the default value from genesis (i.e. unset)
	if currentHeightBlockSize == 0 {
		return &abci.ConsensusParams{}
	}

	lastBlockSize := k.BlockByteSize(previousBlockCtx)
	// If the block size is unchanged, return empty params
	if lastBlockSize == currentHeightBlockSize {
		return &abci.ConsensusParams{}
	}

	if currentHeightBlockSize < types.DefaultBlockByteSize && codec.TestMode > -4 {
		ctx.Logger().Error("block size is less than default value, this should never happen")
		return &abci.ConsensusParams{}
	}

	// If the block size has changed
	return &abci.ConsensusParams{
		Block: &abci.BlockParams{
			MaxBytes: currentHeightBlockSize,
			MaxGas:   -1,
		},
		// INVESTIGATE: Looks like an extra measure to prevent the evidence pool from filling up during the upgrade
		Evidence: &abci.EvidenceParams{
			MaxAge: 50,
		},
	}
}
