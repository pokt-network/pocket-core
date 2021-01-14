package keeper

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

// SetUnstakingApplication - Store an application address to the appropriate position in the unstaking queue
func (k Keeper) SetUnstakingApplication(ctx sdk.Ctx, val types.Application) {
	applications := k.getUnstakingApplications(ctx, val.UnstakingCompletionTime)
	applications = append(applications, val.Address)
	k.setUnstakingApplications(ctx, val.UnstakingCompletionTime, applications)
}

// deleteUnstakingApplicaiton - DeleteEvidence an application address from the unstaking queue
func (k Keeper) deleteUnstakingApplication(ctx sdk.Ctx, val types.Application) {
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

// getAllUnstakingApplications - Retrieve the set of all unstaking applications with no limits
func (k Keeper) getAllUnstakingApplications(ctx sdk.Ctx) (applications []types.Application) {
	applications = make(types.Applications, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.UnstakingAppsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addrs sdk.Addresses
		err := k.Cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &addrs, ctx.BlockHeight())
		if err != nil {
			k.Logger(ctx).Error(fmt.Errorf("could not unmarshal unstakingApplications in getAllUnstakingApplications call: %s", string(iterator.Value())).Error())
			return
		}
		for _, addr := range addrs {
			app, found := k.GetApplication(ctx, addr)
			if !found {
				k.Logger(ctx).Error(fmt.Errorf("application %s in unstakingSet but not found in all applications store", app.Address).Error())
				continue
			}
			applications = append(applications, app)
		}

	}
	return applications
}

// getUnstakingApplications - Retrieve all of the applications who will be unstaked at exactly this time
func (k Keeper) getUnstakingApplications(ctx sdk.Ctx, unstakingTime time.Time) (valAddrs sdk.Addresses) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := store.Get(types.KeyForUnstakingApps(unstakingTime))
	if bz == nil {
		return []sdk.Address{}
	}
	err := k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &valAddrs, ctx.BlockHeight())
	if err != nil {
		panic(err)
	}
	return valAddrs

}

// setUnstakingApplications - Store applications in unstaking queue at a certain unstaking time
func (k Keeper) setUnstakingApplications(ctx sdk.Ctx, unstakingTime time.Time, keys sdk.Addresses) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.Cdc.MarshalBinaryLengthPrefixed(&keys, ctx.BlockHeight())
	if err != nil {
		panic(err)
	}
	_ = store.Set(types.KeyForUnstakingApps(unstakingTime), bz)
}

// delteUnstakingApplications - Remove all the applications for a specific unstaking time
func (k Keeper) deleteUnstakingApplications(ctx sdk.Ctx, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForUnstakingApps(unstakingTime))
}

// unstakingApplicationsIterator - Retrieve an iterator for all unstaking applications up to a certain time
func (k Keeper) unstakingApplicationsIterator(ctx sdk.Ctx, endTime time.Time) (sdk.Iterator, error) {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingAppsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingApps(endTime)))
}

// getMatureApplications - Retrieve a list of all the mature validators
func (k Keeper) getMatureApplications(ctx sdk.Ctx) (matureValsAddrs sdk.Addresses) {
	matureValsAddrs = make([]sdk.Address, 0)
	unstakingValsIterator, _ := k.unstakingApplicationsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()
	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var applications sdk.Addresses
		err := k.Cdc.UnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &applications, ctx.BlockHeight())
		if err != nil {
			panic(err)
		}
		matureValsAddrs = append(matureValsAddrs, applications...)

	}
	return matureValsAddrs
}

// unstakeAllMatureValidators - Unstake all the unstaking applications that have finished their unstaking period
func (k Keeper) unstakeAllMatureApplications(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	unstakingApplicationsIterator, _ := k.unstakingApplicationsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingApplicationsIterator.Close()
	for ; unstakingApplicationsIterator.Valid(); unstakingApplicationsIterator.Next() {
		var unstakingVals sdk.Addresses
		err := k.Cdc.UnmarshalBinaryLengthPrefixed(unstakingApplicationsIterator.Value(), &unstakingVals, ctx.BlockHeight())
		if err != nil {
			panic(err)
		}
		for _, valAddr := range unstakingVals {
			val, found := k.GetApplication(ctx, valAddr)
			if !found {
				k.Logger(ctx).Error(fmt.Errorf("application %s, in the unstaking queue was not found", valAddr).Error())
				continue
			}
			err := k.ValidateApplicationFinishUnstaking(ctx, val)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("Could not finish unstaking mature application at height %d: ", ctx.BlockHeight()) + err.Error())
				continue
			}
			k.FinishUnstakingApplication(ctx, val)
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeCompleteUnstaking,
					sdk.NewAttribute(types.AttributeKeyApplication, valAddr.String()),
				),
			)
			if ctx.IsAfterUpgradeHeight() {
				k.DeleteApplication(ctx, valAddr)
			}
		}
		_ = store.Delete(unstakingApplicationsIterator.Key())
	}
}
