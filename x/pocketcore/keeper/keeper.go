package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/rpc/client"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	posKeeper         types.PosKeeper
	appKeeper         types.AppsKeeper
	TmNode            client.Client
	hostedBlockchains *types.HostedBlockchains
	Paramstore        sdk.Subspace
	storeKey          sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc               *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the pocketcore module Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, posKeeper types.PosKeeper, appKeeper types.AppsKeeper, hostedChains *types.HostedBlockchains, paramstore sdk.Subspace) Keeper {
	return Keeper{
		storeKey:          storeKey,
		cdc:               cdc,
		posKeeper:         posKeeper,
		appKeeper:         appKeeper,
		hostedBlockchains: hostedChains,
		Paramstore:        paramstore.WithKeyTable(ParamKeyTable()),
	}
}

// "GetBlock" returns the block from the tendermint node at a certain height
func (k Keeper) GetBlock(height int) (*core_types.ResultBlock, error) {
	h := int64(height)
	return k.TmNode.Block(&h)
}
