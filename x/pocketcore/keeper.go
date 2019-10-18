package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/blockchain"
)

// PocketCoreKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type PocketCoreKeeper struct {
	ApplicationKeeper blockchain.ApplicationKeeper
	NodeKeeper        blockchain.NodeKeeper
	storeKey          sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc               *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewPocketCoreKeeper creates new instances of the pocketcore PocketCoreKeeper
func NewPocketCoreKeeper(appKeeper blockchain.ApplicationKeeper, nodeKeeper blockchain.NodeKeeper, storeKey sdk.StoreKey, cdc *codec.Codec) PocketCoreKeeper {
	return PocketCoreKeeper{
		ApplicationKeeper: appKeeper,
		NodeKeeper:        nodeKeeper,
		storeKey:          storeKey,
		cdc:               cdc,
	}
}
