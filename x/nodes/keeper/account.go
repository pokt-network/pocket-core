package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.ValAddress) sdk.Int {
	coins := k.coinKeeper.GetCoins(ctx, sdk.AccAddress(addr))
	return coins.AmountOf(k.StakeDenom(ctx))
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddress sdk.ValAddress, toAddress sdk.ValAddress, amount sdk.Int) sdk.Error {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	err := k.coinKeeper.SendCoins(ctx, sdk.AccAddress(fromAddress), sdk.AccAddress(toAddress), coins)
	if err != nil {
		return err
	}
	return nil
}
