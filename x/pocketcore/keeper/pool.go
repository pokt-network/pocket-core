package keeper

import sdk "github.com/pokt-network/posmint/types"

func (k Keeper) StakeDenom(ctx sdk.Context) string {
	return k.posKeeper.StakeDenom(ctx)
}

func (k Keeper) GetAppStakedTokens(ctx sdk.Context) sdk.Int {
	return k.appKeeper.GetStakedTokens(ctx)
}

func (k Keeper) GetNodesStakedTokens(ctx sdk.Context) sdk.Int {
	return k.posKeeper.GetStakedTokens(ctx)
}

func (k Keeper) GetTotalTokens(ctx sdk.Context) sdk.Int {
	return k.posKeeper.TotalTokens(ctx)
}

func (k Keeper) GetTotalStakedTokens(ctx sdk.Context) sdk.Int {
	return k.GetNodesStakedTokens(ctx).Add(k.GetAppStakedTokens(ctx))
}

func (k Keeper) GetStakedRatio(ctx sdk.Context) sdk.Dec {
	totalStaked := k.GetTotalStakedTokens(ctx).ToDec()
	return totalStaked.QuoInt(k.GetTotalTokens(ctx))
}
