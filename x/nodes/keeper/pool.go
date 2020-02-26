package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/supply/exported"
)

// GetStakedTokens total staking tokens supply which is staked
func (k Keeper) GetStakedTokens(ctx sdk.Ctx) sdk.Int {
	stakedPool := k.GetStakedPool(ctx)
	return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx))
}

// GetUnstakedTokens returns the amount of not staked tokens
func (k Keeper) GetUnstakedTokens(ctx sdk.Ctx) (unstakedTokens sdk.Int) {
	return k.TotalTokens(ctx).Sub(k.GetStakedPool(ctx).GetCoins().AmountOf(k.StakeDenom(ctx)))
}

// TotalTokens staking tokens from the total supply
func (k Keeper) TotalTokens(ctx sdk.Ctx) sdk.Int {
	return k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(k.StakeDenom(ctx))
}

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetStakedPool(ctx sdk.Ctx) (stakedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.StakedPoolName)
}

// moves coins from the module account to the validator -> used in unstaking
func (k Keeper) coinsFromStakedToUnstaked(ctx sdk.Ctx, validator types.Validator) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), validator.StakedTokens))
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, sdk.Address(validator.Address), coins)
	if err != nil {
		panic(err)
	}
}

// moves coins from the module account to validator -> used in staking
func (k Keeper) coinsFromUnstakedToStaked(ctx sdk.Ctx, validator types.Validator, amount sdk.Int) sdk.Error {
	if amount.LT(sdk.ZeroInt()) {
		return sdk.ErrInternal("cannot send a negative")
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, validator.Address, types.StakedPoolName, coins)
	return err
}

// burnStakedTokens removes coins from the staked pool module account
func (k Keeper) burnStakedTokens(ctx sdk.Ctx, amt sdk.Int) sdk.Error {
	if !amt.IsPositive() {
		return nil
	}
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amt))
	return k.supplyKeeper.BurnCoins(ctx, types.StakedPoolName, coins)
}

func (k Keeper) getFeePool(ctx sdk.Ctx) (feePool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, auth.FeeCollectorName)
}
