package keeper

import (
	"bytes"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"time"
)

// Insert a application address to the appropriate position in the unstaking queue
func (k Keeper) SetUnstakingApplication(ctx sdk.Context, val types.Application) {
	applications := k.getUnstakingApplications(ctx, val.UnstakingCompletionTime)
	applications = append(applications, val.Address)
	k.setUnstakingApplications(ctx, val.UnstakingCompletionTime, applications)
}

// DeleteInvoice a application address from the unstaking queue
func (k Keeper) deleteUnstakingApplication(ctx sdk.Context, val types.Application) {
	applications := k.getUnstakingApplications(ctx, val.UnstakingCompletionTime)
	var newApplications []sdk.Address
	for _, addr := range applications {
		if !bytes.Equal(addr, val.Address) {
			newApplications = append(newApplications, addr)
		}
	}
	if len(newApplications) == 0 {
		k.deleteUnstakingApplications(ctx, val.UnstakingCompletionTime)
	} else {
		k.setUnstakingApplications(ctx, val.UnstakingCompletionTime, newApplications)
	}
}

// get the set of all unstaking applications with no limits
func (k Keeper) getAllUnstakingApplications(ctx sdk.Context) (applications []types.Application) {
	applications = make(types.Applications, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakingAppsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addrs []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &addrs)
		for _, addr := range addrs {
			validator := k.mustGetApplication(ctx, addr)
			applications = append(applications, validator)
		}
	}
	return applications
}

// gets all of the applications who will be unstaked at exactly this time
func (k Keeper) getUnstakingApplications(ctx sdk.Context, unstakingTime time.Time) (valAddrs []sdk.Address) {
	valAddrs = make([]sdk.Address, 0)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForUnstakingApps(unstakingTime))
	if bz == nil {
		return []sdk.Address{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return valAddrs
}

// Sets applications in unstaking queue at a certain unstaking time
func (k Keeper) setUnstakingApplications(ctx sdk.Context, unstakingTime time.Time, keys []sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.KeyForUnstakingApps(unstakingTime), bz)
}

// Deletes all the applications for a specific unstaking time
func (k Keeper) deleteUnstakingApplications(ctx sdk.Context, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForUnstakingApps(unstakingTime))
}

// iterator for all unstaking applications up to a certain time
func (k Keeper) unstakingApplicationsIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingAppsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingApps(endTime)))
}

// Returns a list of all the mature applications
func (k Keeper) getMatureApplications(ctx sdk.Context) (matureValsAddrs []sdk.Address) {
	matureValsAddrs = make([]sdk.Address, 0)
	unstakingValsIterator := k.unstakingApplicationsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()
	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var applications []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &applications)
		matureValsAddrs = append(matureValsAddrs, applications...)
	}
	return matureValsAddrs
}

// Unstakes all the unstaking applications that have finished their unstaking period
func (k Keeper) unstakeAllMatureApplications(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	unstakingApplicationsIterator := k.unstakingApplicationsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingApplicationsIterator.Close()
	for ; unstakingApplicationsIterator.Valid(); unstakingApplicationsIterator.Next() {
		var unstakingVals []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingApplicationsIterator.Value(), &unstakingVals)
		for _, valAddr := range unstakingVals {
			val, found := k.GetApplication(ctx, valAddr)
			if !found {
				panic("application in the unstaking queue was not found")
			}
			err := k.ValidateApplicationFinishUnstaking(ctx, val)
			if err != nil {
				panic(err)
			}
			err = k.FinishUnstakingApplication(ctx, val)
			if err != nil {
				panic(err)
			}
		}
		store.Delete(unstakingApplicationsIterator.Key())
	}
}
