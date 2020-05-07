package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/exported"
)

// StakedRatio - Retrieve the fraction of the staking tokens which are currently staked
func (k Keeper) StakedRatio(ctx sdk.Ctx) sdk.Dec {
	stakedPool := k.GetStakedPool(ctx)

	stakeSupply := k.TotalTokens(ctx)
	if stakeSupply.IsPositive() {
		return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx)).ToDec().QuoInt(stakeSupply)
	}
	return sdk.ZeroDec()
}

// GetStakedTokens - Retrieve total staking tokens supply which is staked
func (k Keeper) GetStakedTokens(ctx sdk.Ctx) sdk.Int {
	stakedPool := k.GetStakedPool(ctx)
	return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx))
}

// TotalTokens - Retrieve total staking tokens from the total supply
func (k Keeper) TotalTokens(ctx sdk.Ctx) sdk.Int {
	return k.AccountsKeeper.GetSupply(ctx).GetTotal().AmountOf(k.StakeDenom(ctx))
}

// GetStakedPool - Retrieve the staked tokens pool's module account
func (k Keeper) GetStakedPool(ctx sdk.Ctx) (stakedPool exported.ModuleAccountI) {
	return k.AccountsKeeper.GetModuleAccount(ctx, types.StakedPoolName)
}

// coinsFromStakedToUnstkaed - Transfer coins from the module account to the application -> used in unstaking
func (k Keeper) coinsFromStakedToUnstaked(ctx sdk.Ctx, application types.Application) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), application.StakedTokens))
	err := k.AccountsKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, application.Address, coins)
	if err != nil {
		panic(err)
	}
}

// coinsFromUnstakedToStaked - Transfer coins from the module account to application -> used in staking
func (k Keeper) coinsFromUnstakedToStaked(ctx sdk.Ctx, application types.Application, amount sdk.Int) sdk.Error {
	if amount.LT(sdk.ZeroInt()) {
		return sdk.ErrInternal("cannot stake a negative amount of coins")
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.AccountsKeeper.SendCoinsFromAccountToModule(ctx, sdk.Address(application.Address), types.StakedPoolName, coins)
	if err != nil {
		return err
	}
	return nil
}

// burnStakedTokens - Remove coins from the staked pool module account
func (k Keeper) burnStakedTokens(ctx sdk.Ctx, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		return nil
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amt))
	return k.AccountsKeeper.BurnCoins(ctx, types.StakedPoolName, coins)
}

// getFeePool - Retrieve fee pool
func (k Keeper) getFeePool(ctx sdk.Ctx) (feePool exported.ModuleAccountI) {
	return k.AccountsKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
}
