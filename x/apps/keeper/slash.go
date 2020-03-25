package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
)

func (k Keeper) BurnApplication(ctx sdk.Ctx, address sdk.Address, amount sdk.Int) {
	curBurn, _ := k.getApplicationBurn(ctx, address)
	newSeverity := curBurn.Add(amount)
	k.setApplicationBurn(ctx, newSeverity, address)
}

// simpleSlash a application for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) simpleSlash(ctx sdk.Ctx, consAddr sdk.Address, amount sdk.Int) {
	// error check simpleSlash
	application := k.validateSimpleSlash(ctx, consAddr, amount)
	if application.Address == nil {
		return // invalid simpleSlash
	}
	logger := k.Logger(ctx)
	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(amount, application.StakedTokens)
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
	// Log that a simpleSlash occurred
	logger.Info(fmt.Sprintf("application %s simple slashed: burned %v tokens",
		application.GetAddress(), amount.String()))
}

func (k Keeper) validateSimpleSlash(ctx sdk.Ctx, addr sdk.Address, amount sdk.Int) types.Application {
	logger := k.Logger(ctx)
	if amount.LTE(sdk.ZeroInt()) {
		panic(fmt.Errorf("attempted to simpleSlash with a negative simpleSlash factor: %v", amount))
	}
	application, found := k.GetApplication(ctx, addr)
	if !found {
		logger.Error(fmt.Sprintf( // could've been overslashed and removed
			"WARNING: Ignored attempt to simpleSlash a nonexistent application with address %s, we recommend you investigate immediately",
			addr))
		return types.Application{}
	}
	// should not be slashing an unstaked application
	if application.IsUnstaked() {
		panic(fmt.Errorf("should not be slashing unstaked application: %s", application.GetAddress()))
	}
	return application
}

func (k Keeper) getBurnFromSeverity(ctx sdk.Ctx, address sdk.Address, severityPercentage sdk.Dec) sdk.Int {
	app := k.mustGetApplication(ctx, address)
	amount := sdk.TokensFromConsensusPower(app.ConsensusPower())
	slashAmount := amount.ToDec().Mul(severityPercentage).TruncateInt()
	return slashAmount
}

// called on begin blocker
func (k Keeper) burnApplications(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BurnApplicationKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		severity := sdk.ZeroInt()
		address := sdk.Address(types.AddressFromKey(iterator.Key()))
		amino.MustUnmarshalBinaryBare(iterator.Value(), &severity)
		k.simpleSlash(ctx, address, severity)
		// remove from the burn store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a application burn
func (k Keeper) setApplicationBurn(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForAppBurn(address), amino.MustMarshalBinaryBare(amount))
}

func (k Keeper) getApplicationBurn(ctx sdk.Ctx, address sdk.Address) (coins sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForAppBurn(address))
	if value == nil {
		return sdk.ZeroInt(), false
	}
	found = true
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	return
}

func (k Keeper) deleteApplicationBurn(ctx sdk.Ctx, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForAppBurn(address))
}
