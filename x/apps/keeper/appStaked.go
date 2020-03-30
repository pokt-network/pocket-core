package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Set staked application by address in the store.
func (k Keeper) SetStakedApplication(ctx sdk.Ctx, application types.Application) {
	if application.Jailed {
		return // jailed applications are not kept in the staking set
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForAppInStakingSet(application), application.Address)
	ctx.Logger().Info("Setting App on Staking Set " + application.Address.String())
}

// Get the denomination of coins.
func (k Keeper) StakeDenom(ctx sdk.Ctx) string {
	return k.POSKeeper.StakeDenom(ctx)
}

// delete application from staked set
func (k Keeper) deleteApplicationFromStakingSet(ctx sdk.Ctx, application types.Application) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForAppInStakingSet(application))
	ctx.Logger().Info("Removing App From Staking Set " + application.Address.String())
}

// Update the staked tokens of an existing application, update the applications power index key
func (k Keeper) removeApplicationTokens(ctx sdk.Ctx, application types.Application, tokensToRemove sdk.Int) types.Application {
	ctx.Logger().Info("Removing Application Tokens, tokensToRemove: " + tokensToRemove.String() + " App Address: " + application.Address.String())
	k.deleteApplicationFromStakingSet(ctx, application)
	application = application.RemoveStakedTokens(tokensToRemove)
	k.SetApplication(ctx, application)
	k.SetStakedApplication(ctx, application)
	return application
}

func (k Keeper) getStakedApplications(ctx sdk.Ctx) types.Applications {
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
func (k Keeper) stakedAppsIterator(ctx sdk.Ctx) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedAppsKey)
}

// iterate through the staked application set and perform the provided function
func (k Keeper) IterateAndExecuteOverStakedApps(
	ctx sdk.Ctx, fn func(index int64, application exported.ApplicationI) (stop bool)) {
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
