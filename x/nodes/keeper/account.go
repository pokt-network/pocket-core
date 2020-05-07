package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
)

// GetBalance - Retrieve balance for account
func (k Keeper) GetBalance(ctx sdk.Ctx, addr sdk.Address) sdk.Int {
	coins := k.AccountKeeper.GetCoins(ctx, addr)
	return coins.AmountOf(k.StakeDenom(ctx))
}

// GetAccount - Retrieve account info
func (k Keeper) GetAccount(ctx sdk.Ctx, addr sdk.Address) (acc *auth.BaseAccount) {
	a := k.AccountKeeper.GetAccount(ctx, addr)
	if a == nil {
		return &auth.BaseAccount{
			Address: sdk.Address{},
		}
	}
	return a.(*auth.BaseAccount)
}

// SendCoins - Deliver coins to account
func (k Keeper) SendCoins(ctx sdk.Ctx, fromAddress sdk.Address, toAddress sdk.Address, amount sdk.Int) sdk.Error {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.AccountKeeper.SendCoins(ctx, fromAddress, toAddress, coins)
	if err != nil {
		return err
	}
	return nil
}
