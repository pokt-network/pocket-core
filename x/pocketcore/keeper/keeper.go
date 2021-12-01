package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
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
