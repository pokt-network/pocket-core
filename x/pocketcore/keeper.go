package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// PocketCoreKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type PocketCoreKeeper struct {
	AccountKeeper     accountKeeper
	ApplicationKeeper applicationKeeper
	BlockKeeper       blockKeeper
	NodeKeeper        nodeKeeper
	PoktKeeper        poktKeeper
	storeKey          sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc               *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewPocketCoreKeeper creates new instances of the pocketcore PocketCoreKeeper
func NewPocketCoreKeeper(coinKeeper bank.Keeper, supplyKeeper supply.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) PocketCoreKeeper {
	return PocketCoreKeeper{
		AccountKeeper:     newAccountKeeper(coinKeeper, storeKey, cdc),
		ApplicationKeeper: newApplicationKeeper(storeKey, cdc),
		BlockKeeper:       newBlockKeeper(storeKey, cdc),
		NodeKeeper:        newNodeKeeper(storeKey, cdc),
		PoktKeeper:        newPoktKeeper(supplyKeeper, storeKey, cdc),
		storeKey:          storeKey,
		cdc:               cdc,
	}
}
