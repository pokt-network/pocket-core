package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/x/supply/exported"
)

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetDAOPool(ctx sdk.Context) (stakedPool exported.ModuleAccountI) {
	return k.supplyKeeper.GetModuleAccount(ctx, types.DAOPoolName)
}

// moves coins from the module account to the validator -> used in unstaking
func (k Keeper) coinsFromDAOToValidator(ctx sdk.Context, validator types.Validator, amount sdk.Int) {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.DAOPoolName, sdk.AccAddress(validator.Address), coins)
	if err != nil {
		panic(err)
	}
}

// GetStakedTokens total staking tokens supply which is staked
func (k Keeper) GetDAOTokens(ctx sdk.Context) sdk.Int {
	stakedPool := k.GetDAOPool(ctx)
	return stakedPool.GetCoins().AmountOf(k.StakeDenom(ctx))
}
