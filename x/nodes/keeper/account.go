package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
)

func (k Keeper) GetBalance(ctx sdk.Ctx, addr sdk.Address) sdk.Int {
	coins := k.coinKeeper.GetCoins(ctx, addr)
	return coins.AmountOf(k.StakeDenom(ctx))
}

func (k Keeper) GetAccount(ctx sdk.Ctx, addr sdk.Address) (acc *auth.BaseAccount) {
	account := k.accountKeeper.GetAccount(ctx, addr)
	if account == nil {
		return &auth.BaseAccount{
			Address: sdk.Address{},
		}
	}
	acc = account.(*auth.BaseAccount)
	return
}

func (k Keeper) SendCoins(ctx sdk.Ctx, fromAddress sdk.Address, toAddress sdk.Address, amount sdk.Int) sdk.Error {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.coinKeeper.SendCoins(ctx, fromAddress, toAddress, coins)
	if err != nil {
		return err
	}
	return nil
}
