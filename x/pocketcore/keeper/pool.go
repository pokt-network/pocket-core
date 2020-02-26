package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
)

// get the stake denomination from the node module
func (k Keeper) StakeDenom(ctx sdk.Ctx) (res string) {
	res = k.posKeeper.StakeDenom(ctx)
	return
}

// get the total number of staked tokens for applications
func (k Keeper) GetAppStakedTokens(ctx sdk.Ctx) (res sdk.Int) {
	res = k.appKeeper.GetStakedTokens(ctx)
	return
}

// get the total number of staked tokens for nodes
func (k Keeper) GetNodesStakedTokens(ctx sdk.Ctx) (res sdk.Int) {
	res = k.posKeeper.GetStakedTokens(ctx)
	return
}

// get total tokens in the world supply
func (k Keeper) GetTotalTokens(ctx sdk.Ctx) (res sdk.Int) {
	res = k.posKeeper.TotalTokens(ctx)
	return
}

// get the total staked tokens for both nodes and apps
func (k Keeper) GetTotalStakedTokens(ctx sdk.Ctx) (res sdk.Int) {
	res = k.GetNodesStakedTokens(ctx).Add(k.GetAppStakedTokens(ctx))
	return
}

// get the ratio of staked tokens to unstaked tokens
func (k Keeper) GetStakedRatio(ctx sdk.Ctx) sdk.Dec {
	totalStaked := k.GetTotalStakedTokens(ctx).ToDec()
	totalStaked = totalStaked.QuoInt(k.GetTotalTokens(ctx))
	return totalStaked
}
