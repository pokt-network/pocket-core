package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/rpc/client"
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
	cdc               *codec.Codec // The wire codec for binary encoding/decoding.
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
		cdc:               cdc,
	}
}
