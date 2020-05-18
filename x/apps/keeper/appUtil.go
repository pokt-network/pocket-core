package keeper

import (
	"container/list"
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Cache the amino decoding of applications, as it can be the case that repeated slashing calls
// cause many calls to GetApplication, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedApplication struct {
	app     types.Application
	address sdk.Address
}

func newCachedApplication(val types.Application, address sdk.Address) cachedApplication {
	return cachedApplication{
		app:     val,
		address: address,
	}
}

// appCaching - Retrieve a cached application
func (k Keeper) appCaching(value []byte, addr sdk.Address) types.Application {
	application, _ := types.UnmarshalApplication(k.cdc, value) // TODO fix error this when cache is updated
	cachedVal := newCachedApplication(application, addr)
	k.applicationCache[addr.String()] = cachedVal
	k.applicationCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the prevState element from it
	if int64(k.applicationCacheList.Len()) > types.ApplicationCacheSize {
		appToRemove := k.applicationCacheList.Remove(k.applicationCacheList.Front()).(cachedApplication)
		delete(k.applicationCache, appToRemove.address.String())
	}
	return application
}
func (k Keeper) getApplicationFromCache(addr sdk.Address) (application types.Application, found bool) {
	if app, ok := k.applicationCache[addr.String()]; ok {
		appToReturn := app.app
		// Doesn't mutate the cache's value
		appToReturn.Address = addr
		return appToReturn, true
	} else {
		return types.Application{}, false
	}
}

func (k Keeper) searchCacheList(application types.Application) (e *list.Element, found bool) {
	for e := k.applicationCacheList.Back(); e != nil; e = e.Prev() {
		v := e.Value.(cachedApplication)
		if v.address.String() == application.Address.String() {
			return e, true
		}
	}
	return nil, false
}

func (k Keeper) setOrUpdateInApplicationCache(application types.Application) {

	e, found := k.searchCacheList(application)
	if found {
		appToRemove := k.applicationCacheList.Remove(e).(cachedApplication)
		delete(k.applicationCache, appToRemove.address.String())
	}

	cachedApp := newCachedApplication(application, application.Address)
	k.applicationCache[application.Address.String()] = cachedApp
	k.applicationCacheList.PushBack(cachedApp)
}

// Application - wrapper for GetApplication call
func (k Keeper) Application(ctx sdk.Ctx, address sdk.Address) exported.ApplicationI {
	app, found := k.GetApplication(ctx, address)
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
		app, err := types.UnmarshalApplication(k.cdc, iterator.Value())
		if err != nil {
			k.Logger(ctx).Error("couldn't unmarshal application in AllApplications call: " + string(iterator.Value()) + "\n" + err.Error())
			continue
		}
		apps = append(apps, app)
	}
	return apps
}
