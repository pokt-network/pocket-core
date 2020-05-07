package keeper

import (
	"fmt"

	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Cache the amino decoding of applications, as it can be the case that repeated slashing calls
// cause many calls to GetApplication, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedApplication struct {
	val        types.Application
	marshalled string // marshalled amino bytes for the application object (not operator address)
}

func newCachedApplication(val types.Application, marshalled string) cachedApplication {
	return cachedApplication{
		val:        val,
		marshalled: marshalled,
	}
}

// appCaching - Retrieve a cached application
func (k Keeper) appCaching(value []byte, addr sdk.Address) types.Application {
	// If these amino encoded bytes are in the cache, return the cached application
	strValue := string(value)
	if val, ok := k.applicationCache[strValue]; ok {
		valToReturn := val.val
		// Doesn't mutate the cache's value
		valToReturn.Address = addr
		return valToReturn
	}
	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	application := types.MustUnmarshalApplication(k.cdc, value)
	cachedVal := newCachedApplication(application, strValue)
	k.applicationCache[strValue] = newCachedApplication(application, strValue)
	k.applicationCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the prevState element from it
	if k.applicationCacheList.Len() > aminoCacheSize {
		valToRemove := k.applicationCacheList.Remove(k.applicationCacheList.Front()).(cachedApplication)
		delete(k.applicationCache, valToRemove.marshalled)
	}
	return application
}

// mustGetApplication - Retrieve application, panics if no application is found
func (k Keeper) mustGetApplication(ctx sdk.Ctx, addr sdk.Address) types.Application {
	application, found := k.GetApplication(ctx, addr)
	if !found {
		panic(fmt.Sprintf("application record not found for address: %X\n", addr))
	}
	return application
}

// mustGetApplicationByConsAddr - Retrieve application using consensus address, panics if no application is found
func (k Keeper) mustGetApplicationByConsAddr(ctx sdk.Ctx, consAddr sdk.Address) types.Application {
	application, found := k.GetApplication(ctx, consAddr)
	if !found {
		panic(fmt.Errorf("application with consensus-Address %s not found", consAddr))
	}
	return application
}

// Application - wrapper for GetApplication call
func (k Keeper) Application(ctx sdk.Ctx, address sdk.Address) exported.ApplicationI {
	app, found := k.GetApplication(ctx, address)
	if !found {
		return nil
	}
	return app
}

// applicationByConsAddr - wrapper for GetApplicationByConsAddress call
func (k Keeper) applicationByConsAddr(ctx sdk.Ctx, addr sdk.Address) exported.ApplicationI {
	app, found := k.GetApplication(ctx, addr)
	if !found {
		return nil
	}
	return app
}

// AllApplications - Retrieve a list of all applications
func (k Keeper) AllApplications(ctx sdk.Ctx) (apps []exported.ApplicationI) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		app := types.MustUnmarshalApplication(k.cdc, iterator.Value())
		apps = append(apps, app)
	}
	return apps
}
