package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the pocketcore Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetStruct(ctx sdk.Context, name string) struct{} {
	store := ctx.KVStore(k.storeKey)
	if !k.Contains(ctx, name) {
		// return new struct{}
	}
	structEncoding := store.Get([]byte(name))
	var s struct{}
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &s)
	return s
}

// Sets the entire struct for a string
func (k Keeper) SetStruct(ctx sdk.Context, str string, s struct{}) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(str), k.cdc.MustMarshalBinaryBare(s))
}

// Deletes the entire struct metadata struct for a string
func (k Keeper) DeleteStruct(ctx sdk.Context, str string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(str))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetAllIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// Check if the name is present in the store or not
func (k Keeper) Contains(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}
