package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

// SetStakedApplication - Store staked application
func (k Keeper) SetStakedApplication(ctx sdk.Ctx, application types.Application) {
	if application.Jailed {
		return // jailed applications are not kept in the staking set
	}
	store := ctx.KVStore(k.storeKey)
	_ = store.Set(types.KeyForAppInStakingSet(application), application.Address)
	ctx.Logger().Info("Setting App on Staking Set " + application.Address.String())
}

// StakeDenom - Retrieve the denomination of coins.
func (k Keeper) StakeDenom(ctx sdk.Ctx) string {
	return k.POSKeeper.StakeDenom(ctx)
}

// deleteApplicationFromStakingSet - Remove application from staked set
func (k Keeper) deleteApplicationFromStakingSet(ctx sdk.Ctx, application types.Application) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForAppInStakingSet(application))
	ctx.Logger().Info("Removing App From Staking Set " + application.Address.String())
}

// removeApplicationTokens - Update the staked tokens of an existing application, update the applications power index key
func (k Keeper) removeApplicationTokens(ctx sdk.Ctx, application types.Application, tokensToRemove sdk.BigInt) (types.Application, error) {
	ctx.Logger().Info("Removing Application Tokens, tokensToRemove: " + tokensToRemove.String() + " App Address: " + application.Address.String())
	k.deleteApplicationFromStakingSet(ctx, application)
	application, err := application.RemoveStakedTokens(tokensToRemove)
	if err != nil {
		return types.Application{}, err
	}
	k.SetApplication(ctx, application)
	return application, nil
}

// getStakedApplications - Retrieve the current staked applications sorted by power-rank
func (k Keeper) getStakedApplications(ctx sdk.Ctx) types.Applications {
	var applications = make(types.Applications, 0)
	iterator, _ := k.stakedAppsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		application, found := k.GetApplication(ctx, address)
		if !found {
			k.Logger(ctx).Error(fmt.Errorf("application %s in staking set but not found in all applications store", address).Error())
			continue
		}
		if application.IsStaked() {
			applications = append(applications, application)
		}
	}
	return applications
}

// getStakedApplicationsCount returns a count of the total staked applcations currently
func (k Keeper) getStakedApplicationsCount(ctx sdk.Ctx) (count int64) {
	iterator, _ := k.stakedAppsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		count++
	}
	return
}

// stakedAppsIterator - Retrieve an iterator for the current staked applications
func (k Keeper) stakedAppsIterator(ctx sdk.Ctx) (sdk.Iterator, error) {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedAppsKey)
}

// IterateAndExecuteOverStakedApps - Goes through the staked application set and execute handler
func (k Keeper) IterateAndExecuteOverStakedApps(
	ctx sdk.Ctx, fn func(index int64, application exported.ApplicationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStoreReversePrefixIterator(store, types.StakedAppsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		application, found := k.GetApplication(ctx, address)
		if !found {
			k.Logger(ctx).Error(fmt.Errorf("application %s in staking set but not found in all applications store", address).Error())
			continue
		}
		if application.IsStaked() {
			stop := fn(i, application) // XXX is this safe will the application unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}
