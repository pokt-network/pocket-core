package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// get a single application from the main store
func (k Keeper) GetApplication(ctx sdk.Ctx, addr sdk.Address) (application types.Application, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForAppByAllApps(addr))
	if value == nil {
		return application, false
	}
	application = k.appCaching(value, addr)
	return application, true
}

// set a application in the main store
func (k Keeper) SetApplication(ctx sdk.Ctx, application types.Application) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalApplication(k.cdc, application)
	store.Set(types.KeyForAppByAllApps(application.Address), bz)
	ctx.Logger().Info("Setting App on Main Store " + application.Address.String())

}

// get the set of all applications with no limits from the main store
func (k Keeper) GetAllApplications(ctx sdk.Ctx) (applications types.Applications) {
	applications = make([]types.Application, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		application := types.MustUnmarshalApplication(k.cdc, iterator.Value())
		applications = append(applications, application)
	}
	return applications
}

// get the set of all applications with no limits from the main store
func (k Keeper) GetAllApplicationsWithOpts(ctx sdk.Ctx, opts types.QueryApplicationsWithOpts) (applications types.Applications) {
	applications = make([]types.Application, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		application := types.MustUnmarshalApplication(k.cdc, iterator.Value())

		applications = append(applications, application)
	}
	return applications
}

// return a given amount of all the applications
func (k Keeper) GetApplications(ctx sdk.Ctx, maxRetrieve uint16) (applications types.Applications) {
	store := ctx.KVStore(k.storeKey)
	applications = make([]types.Application, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		application := types.MustUnmarshalApplication(k.cdc, iterator.Value())
		applications[i] = application
		i++
	}
	return applications[:i] // trim if the array length < maxRetrieve
}

// iterate through the application set and perform the provided function
func (k Keeper) IterateAndExecuteOverApps(
	ctx sdk.Ctx, fn func(index int64, application exported.ApplicationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		application := types.MustUnmarshalApplication(k.cdc, iterator.Value())
		stop := fn(i, application) // XXX is this safe will the application unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

func (k Keeper) CalculateAppRelays(ctx sdk.Ctx, application types.Application) sdk.Int {
	stakingAdjustment := sdk.NewDec(k.StakingAdjustment(ctx))
	participationRate := sdk.NewDec(1)
	baseRate := sdk.NewInt(k.BaselineThroughputStakeRate(ctx))
	if k.ParticipationRateOn(ctx) {
		appStakedCoins := k.GetStakedTokens(ctx)
		nodeStakedCoins := k.POSKeeper.GetStakedTokens(ctx)
		totalTokens := k.TotalTokens(ctx)
		participationRate = appStakedCoins.Add(nodeStakedCoins).ToDec().Quo(totalTokens.ToDec())
	}
	basePercentage := baseRate.ToDec().Quo(sdk.NewDec(100))
	baselineThroughput := basePercentage.Mul(application.StakedTokens.ToDec())
	return participationRate.Mul(baselineThroughput).Add(stakingAdjustment).TruncateInt()
}
