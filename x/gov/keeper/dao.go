package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	exported2 "github.com/pokt-network/pocket-core/x/auth/exported"
	"github.com/pokt-network/pocket-core/x/gov/types"
)

func (k Keeper) DAOTransferFrom(ctx sdk.Ctx, owner, to sdk.Address, amount sdk.BigInt) sdk.Result {
	if !k.GetDAOOwner(ctx).Equals(owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("non dao owner is trying to transfer from the dao %s", owner.String())).Result()
	}
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, amount))
	err := k.AuthKeeper.SendCoinsFromModuleToAccount(ctx, types.DAOAccountName, to, coins)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventDAOTransfer,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) DAOBurn(ctx sdk.Ctx, owner sdk.Address, amount sdk.BigInt) sdk.Result {
	if !k.GetDAOOwner(ctx).Equals(owner) {
		return sdk.ErrUnauthorized(fmt.Sprintf("non dao owner is trying to burn from the dao %s", owner.String())).Result()
	}
	coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, amount))
	err := k.AuthKeeper.BurnCoins(ctx, types.DAOAccountName, coins)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventDAOBurn,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, owner.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func (k Keeper) GetDAOTokens(ctx sdk.Ctx) sdk.BigInt {
	return k.GetDAOAccount(ctx).GetCoins().AmountOf(sdk.DefaultStakeDenom)
}

// GetStakedPool returns the staked tokens pool's module account
func (k Keeper) GetDAOAccount(ctx sdk.Ctx) (stakedPool exported2.ModuleAccountI) {
	return k.AuthKeeper.GetModuleAccount(ctx, types.DAOAccountName)
}
