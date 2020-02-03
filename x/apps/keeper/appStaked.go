package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// set staked application
func (k Keeper) SetStakedApplication(ctx sdk.Context, application types.Application) {
	if application.Jailed {
		return // jailed applications are not kept in the power index
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForAppInStakingSet(application), application.Address)
}

func (k Keeper) StakeDenom(ctx sdk.Context) string {
	return k.posKeeper.StakeDenom(ctx)
}

// delete application from staked set
func (k Keeper) deleteApplicationFromStakingSet(ctx sdk.Context, application types.Application) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForAppInStakingSet(application))
}

// Update the staked tokens of an existing application, update the applications power index key
func (k Keeper) removeApplicationTokens(ctx sdk.Context, v types.Application, tokensToRemove sdk.Int) types.Application {
	k.deleteApplicationFromStakingSet(ctx, v)
	v = v.RemoveStakedTokens(tokensToRemove)
	k.SetApplication(ctx, v)
	k.SetStakedApplication(ctx, v)
	return v
}

// Update the staked tokens of an existing application, update the applications power index key
func (k Keeper) removeApplicationRelays(ctx sdk.Context, v types.Application, relaysToRemove sdk.Int) types.Application {
	k.deleteApplicationFromStakingSet(ctx, v)
	v.MaxRelays = v.MaxRelays.Sub(relaysToRemove)
	k.SetApplication(ctx, v)
	k.SetStakedApplication(ctx, v)
	return v
}

func (k Keeper) getStakedApplications(ctx sdk.Context) types.Applications {
	var applications = make(types.Applications, 0)
	iterator := k.stakedAppsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		application := k.mustGetApplication(ctx, address)
		if application.IsStaked() {
			applications = append(applications, application)
		}
	}
	return applications
}

// returns an iterator for the current staked applications
func (k Keeper) stakedAppsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedAppsKey)
}

// iterate through the staked application set and perform the provided function
func (k Keeper) IterateAndExecuteOverStakedApps(
	ctx sdk.Context, fn func(index int64, application exported.ApplicationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedAppsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		application := k.mustGetApplication(ctx, address)
		if application.IsStaked() {
			stop := fn(i, application) // XXX is this safe will the application unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}
