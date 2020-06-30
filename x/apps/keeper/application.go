package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"math"
	"math/big"
	"os"
)

// GetApplication - Retrieve a single application from the main store
func (k Keeper) GetApplication(ctx sdk.Ctx, addr sdk.Address) (application types.Application, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForAppByAllApps(addr))
	if value == nil {
		return application, false
	}
	application, err := types.UnmarshalApplication(k.cdc, value)
	if err != nil {
		k.Logger(ctx).Error("could not unmarshal application from store")
		return application, false
	}
	return application, true
}

// SetApplication - Add a single application the main store
func (k Keeper) SetApplication(ctx sdk.Ctx, application types.Application) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(application)
	if err != nil {
		k.Logger(ctx).Error("could not marshal application object")
		os.Exit(1)
	}
	store.Set(types.KeyForAppByAllApps(application.Address), bz)
	ctx.Logger().Info("Setting App on Main Store " + application.Address.String())
}

// GetAllApplications - Retrieve the set of all applications with no limits from the main store
func (k Keeper) GetAllApplications(ctx sdk.Ctx) (applications types.Applications) {
	applications = make([]types.Application, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		application, err := types.UnmarshalApplication(k.cdc, iterator.Value())
		if err != nil {
			k.Logger(ctx).Error("couldn't unmarshal application in GetAllApplications call: " + string(iterator.Value()) + "\n" + err.Error())
			continue
		}
		applications = append(applications, application)
	}
	return applications
}

// GetAllApplicationsWithOpts - Retrieve the set of all applications with no limits from the main store
func (k Keeper) GetAllApplicationsWithOpts(ctx sdk.Ctx, opts types.QueryApplicationsWithOpts) (applications types.Applications) {
	applications = make([]types.Application, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		application, err := types.UnmarshalApplication(k.cdc, iterator.Value())
		if err != nil {
			k.Logger(ctx).Error("couldn't unmarshal application in GetAllApplicationsWithOpts call: " + string(iterator.Value()) + "\n" + err.Error())
			continue
		}
		if opts.IsValid(application) {
			applications = append(applications, application)
		}
	}
	return applications
}

// GetApplications - Retrieve a a given amount of all the applications
func (k Keeper) GetApplications(ctx sdk.Ctx, maxRetrieve uint16) (applications types.Applications) {
	store := ctx.KVStore(k.storeKey)
	applications = make([]types.Application, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		application, err := types.UnmarshalApplication(k.cdc, iterator.Value())
		if err != nil {
			k.Logger(ctx).Error("couldn't unmarshal application in GetApplications call: " + string(iterator.Value()) + "\n" + err.Error())
			continue
		}
		applications[i] = application
		i++
	}
	return applications[:i] // trim if the array length < maxRetrieve
}

// IterateAndExecuteOverApps - Goes through the application set and perform the provided function
func (k Keeper) IterateAndExecuteOverApps(
	ctx sdk.Ctx, fn func(index int64, application exported.ApplicationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllApplicationsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		application, err := types.UnmarshalApplication(k.cdc, iterator.Value())
		if err != nil {
			k.Logger(ctx).Error("couldn't unmarshal application in IterateAndExecuteOverApps call: " + string(iterator.Value()) + "\n" + err.Error())
			continue
		}
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
	baselineThroughput := basePercentage.Mul(application.StakedTokens.ToDec().Quo(sdk.NewDec(1000000)))
	result := participationRate.Mul(baselineThroughput).Add(stakingAdjustment).TruncateInt()

	// bounding Max Amount of relays Value to be 18,446,744,073,709,551,615
	maxRelays := sdk.NewIntFromBigInt(new(big.Int).SetUint64(math.MaxUint64))
	if result.GTE(maxRelays) {
		result = maxRelays
	}

	return result
}
