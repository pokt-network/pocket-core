package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
)

func (k Keeper) BurnApplication(ctx sdk.Context, address sdk.ValAddress, severityPercentage sdk.Dec) {
	curBurn, _ := k.getApplicationBurn(ctx, address)
	newSeverity := curBurn.Add(severityPercentage)
	k.setApplicationBurn(ctx, newSeverity, address)
}

// slash a application for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight, power int64, slashFactor sdk.Dec) {
	// error check slash
	application := k.validateSlash(ctx, consAddr, infractionHeight, power, slashFactor)
	if application.Address == nil {
		return // invalid slash
	}
	logger := k.Logger(ctx)
	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := sdk.TokensFromConsensusPower(power)
	slashAmount := amount.ToDec().Mul(slashFactor).TruncateInt()
	k.BeforeApplicationSlashed(ctx, application.Address, slashFactor)
	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(slashAmount, application.StakedTokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.
	// Deduct from application's staked tokens and update the application.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	application = k.removeApplicationTokens(ctx, application, tokensToBurn)
	err := k.burnStakedTokens(ctx, tokensToBurn)
	if err != nil {
		panic(err)
	}
	// if falls below minimum force burn all of the stake
	if application.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		err := k.ForceApplicationUnstake(ctx, application)
		if err != nil {
			panic(err)
		}
	}
	// Log that a slash occurred
	logger.Info(fmt.Sprintf("application %s slashed by slash factor of %s; burned %v tokens",
		application.GetAddress(), slashFactor.String(), tokensToBurn))
	k.AfterApplicationSlashed(ctx, application.Address, slashFactor)
}

func (k Keeper) validateSlash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) types.Application {
	logger := k.Logger(ctx)
	if slashFactor.LT(sdk.ZeroDec()) {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}
	if infractionHeight > ctx.BlockHeight() {
		panic(fmt.Errorf( // Can't slash infractions in the future
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))
	}
	// see if infraction height is outside of unstaking time
	blockTime := ctx.BlockTime()
	infractionTime := ctx.WithBlockHeight(infractionHeight).BlockTime()
	if blockTime.After(infractionTime.Add(k.UnStakingTime(ctx))) {
		logger.Info(fmt.Sprintf( // could've been overslashed and removed
			"INFO: tried to slash with expired evidence: %s %s", infractionTime, blockTime))
		return types.Application{}
	}
	application, found := k.GetAppByConsAddr(ctx, consAddr)
	if !found {
		logger.Error(fmt.Sprintf( // could've been overslashed and removed
			"WARNING: Ignored attempt to slash a nonexistent application with address %s, we recommend you investigate immediately",
			consAddr))
		return types.Application{}
	}
	// should not be slashing an unstaked application
	if application.IsUnstaked() {
		panic(fmt.Errorf("should not be slashing unstaked application: %s", application.GetAddress()))
	}
	return application
}

func (k Keeper) getBurnFromSeverity(ctx sdk.Context, address sdk.ValAddress, severityPercentage sdk.Dec) sdk.Int {
	app := k.mustGetApplication(ctx, address)
	amount := sdk.TokensFromConsensusPower(app.ConsensusPower())
	slashAmount := amount.ToDec().Mul(severityPercentage).TruncateInt()
	return slashAmount
}

// called on begin blocker
func (k Keeper) burnApplications(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BurnApplicationKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		severity := sdk.Dec{}
		address := sdk.ValAddress(types.AddressFromKey(iterator.Key()))
		amino.MustUnmarshalBinaryBare(iterator.Value(), &severity)
		val := k.mustGetApplication(ctx, address)
		k.slash(ctx, sdk.ConsAddress(address), ctx.BlockHeight(), val.ConsensusPower(), severity)
		// remove from the burn store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a application burn
func (k Keeper) setApplicationBurn(ctx sdk.Context, amount sdk.Dec, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForAppBurn(address), amino.MustMarshalBinaryBare(amount))
}

func (k Keeper) getApplicationBurn(ctx sdk.Context, address sdk.ValAddress) (coins sdk.Dec, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForAppBurn(address))
	if value == nil {
		return coins, false
	}
	found = true
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	return
}

func (k Keeper) deleteApplicationBurn(ctx sdk.Context, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForAppBurn(address))
}
