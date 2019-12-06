package keeper

import sdk "github.com/pokt-network/posmint/types"

// get the stake denomination from the node module
func (k Keeper) StakeDenom(ctx sdk.Context) string {
	return k.posKeeper.StakeDenom(ctx)
}

// get the total number of staked tokens for applications
func (k Keeper) GetAppStakedTokens(ctx sdk.Context) sdk.Int {
	return k.appKeeper.GetStakedTokens(ctx)
}

// get the total number of staked tokens for nodes
func (k Keeper) GetNodesStakedTokens(ctx sdk.Context) sdk.Int {
	return k.posKeeper.GetStakedTokens(ctx)
}

// get total tokens in the world supply
func (k Keeper) GetTotalTokens(ctx sdk.Context) sdk.Int {
	return k.posKeeper.TotalTokens(ctx)
}

// get the total staked tokens for both nodes and apps
func (k Keeper) GetTotalStakedTokens(ctx sdk.Context) sdk.Int {
	return k.GetNodesStakedTokens(ctx).Add(k.GetAppStakedTokens(ctx))
}

// get the ratio of staked tokens to unstaked tokens
func (k Keeper) GetStakedRatio(ctx sdk.Context) sdk.Dec {
	totalStaked := k.GetTotalStakedTokens(ctx).ToDec()
	return totalStaked.QuoInt(k.GetTotalTokens(ctx))
}
