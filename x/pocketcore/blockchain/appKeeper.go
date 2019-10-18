package blockchain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pokt-network/pocket-core/types"
)

// ApplicationKeeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type ApplicationKeeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the blockchain ApplicationKeeper
func NewApplicationKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) ApplicationKeeper {
	return ApplicationKeeper{
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Application metadata struct for an address
func (k ApplicationKeeper) GetApplication(ctx sdk.Context, address string) types.Application {
	store := ctx.KVStore(k.storeKey)
	if !k.ContainsApplication(ctx, address) {
		return types.Application{} // todo handle standard error logic
	}
	structEncoding := store.Get([]byte(address))
	var n types.Application
	k.cdc.MustUnmarshalBinaryBare(structEncoding, &n)
	return n
}

// Sets the entire Application for an address
func (k ApplicationKeeper) SetApplication(ctx sdk.Context, address string, n types.Application) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(address), k.cdc.MustMarshalBinaryBare(n))
}

// Deletes the entire metadata struct for an address
func (k ApplicationKeeper) DeleteApplication(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(address))
}

// Check if the application is present in the store or not
func (k ApplicationKeeper) ContainsApplication(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(address))
}

func (k ApplicationKeeper) GetAllApplications(ctx sdk.Context) ([]types.Application, sdk.Error) {
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
func (k ApplicationKeeper) GetApplicationsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
