package keeper

import (
	"errors"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/supply/exported"
)

// StakedRatio the fraction of the staking tokens which are currently staked
func (k Keeper) StakedRatio(ctx sdk.Context) sdk.Dec {
	stakedPool := k.GetStakedPool(ctx)

	stakeSupply := k.TotalTokens(ctx)
	if stakeSupply.IsPositive() {
		return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx)).ToDec().QuoInt(stakeSupply)
	}
	return sdk.ZeroDec()
}

// GetStakedTokens total staking tokens supply which is staked
func (k Keeper) GetStakedTokens(ctx sdk.Context) sdk.Int {
	stakedPool := k.GetStakedPool(ctx)
	return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx))
}

// GetUnstakedTokens returns the amount of not staked tokens
func (k Keeper) GetUnstakedTokens(ctx sdk.Context) (unstakedTokens sdk.Int) {
	return k.TotalTokens(ctx).Sub(k.GetStakedPool(ctx).GetCoins().AmountOf(k.StakeDenom(ctx)))
}

// TotalTokens staking tokens from the total supply
func (k Keeper) TotalTokens(ctx sdk.Context) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.StakeDenom(ctx))
}

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetStakedPool(ctx sdk.Context) (stakedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.StakedPoolName)
}

// moves coins from the module account to the application -> used in unstaking
func (k Keeper) coinsFromStakedToUnstaked(ctx sdk.Context, application types.Application) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), application.StakedTokens))
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, application.Address, coins)
	if err != nil {
		panic(err)
	}
}

// moves coins from the module account to application -> used in staking
func (k Keeper) coinsFromUnstakedToStaked(ctx sdk.Context, application types.Application, amount sdk.Int) error {
	if amount.LT(sdk.ZeroInt()) {
		return errors.New("cannot stake a negative amount of coins")
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sdk.Address(application.Address), types.StakedPoolName, coins)
	if err != nil {
		return err
	}
	return nil
}

// burnStakedTokens removes coins from the staked pool module account
func (k Keeper) burnStakedTokens(ctx sdk.Context, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		return nil
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, types.StakedPoolName, coins)
}

func (k Keeper) getFeePool(ctx sdk.Context) (feePool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
}
