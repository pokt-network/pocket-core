package pocketcore

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/pocket-core/types"
)

// applicationKeeper handles application (type of account) access/modifiers
type applicationKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
}

// newApplicationKeeper creates new instances of the applicationKeeper
func newApplicationKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) applicationKeeper {
	return applicationKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Gets the entire Application metadata struct for an address
func (k applicationKeeper) GetApplication(ctx sdk.Context, address string) types.Application {
	store := ctx.KVStore(k.storeKey)
	if !k.ContainsApplication(ctx, address) {
		return types.Application{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Application
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Gets the entire Application metadata struct for an address at a certain block height
func (k applicationKeeper) GetApplicationAtHeight(ctx sdk.Context, address string, height int64) types.Application {
	atHeight := ctx.WithBlockHeight(height)
	store := atHeight.KVStore(k.storeKey)
	if !k.ContainsApplication(atHeight, address) {
		return types.Application{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Application
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Sets the entire Application for an address
func (k applicationKeeper) SetApplication(ctx sdk.Context, address string, n types.Application) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(n))
}

// Deletes the entire metadata struct for an address
func (k applicationKeeper) DeleteApplication(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(address))
}

// Check if the application is present in the store or not
func (k applicationKeeper) ContainsApplication(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

// Returns all applications within the world state
func (k applicationKeeper) GetAllApplications(ctx sdk.Context) ([]types.Application, sdk.Error) {
	var appList []types.Application
	app := new(types.Application)

	iterator := k.GetApplicationsIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		err := k.cdc.UnmarshalBinaryBare(iterator.Key(), app)
		if err != nil {
			return nil, types.NewError(types.CODEAPPUNMARSHALUNSUCCESSFUL, err.Error())
		}
		appList = append(appList, *app)
	}
	return appList, nil
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k applicationKeeper) GetApplicationsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
