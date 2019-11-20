package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	posKeeper types.PosKeeper
	appKeeper types.AppsKeeper
	storeKey  sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc       *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewPocketCoreKeeper creates new instances of the pocketcore Keeper
func NewPocketCoreKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, posKeeper types.PosKeeper, appKeeper types.AppsKeeper) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		posKeeper: posKeeper,
		appKeeper: appKeeper,
	}
}
