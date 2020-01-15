package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.Address) sdk.Int {
	coins := k.coinKeeper.GetCoins(ctx, sdk.Address(addr))
	return coins.AmountOf(k.StakeDenom(ctx))
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddress sdk.Address, toAddress sdk.Address, amount sdk.Int) sdk.Error {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.coinKeeper.SendCoins(ctx, sdk.Address(fromAddress), sdk.Address(toAddress), coins)
	if err != nil {
		return err
	}
	return nil
}
