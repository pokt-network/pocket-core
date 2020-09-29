package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// "StakeDenom" - Returns the stake coin denomination from the node module
func (k Keeper) StakeDenom(ctx sdk.Ctx) (res string) {
	res = k.posKeeper.StakeDenom(ctx)
	return
}

// "GetAppStakedTokens" - Returns the total number of staked tokens in the apps module
func (k Keeper) GetAppStakedTokens(ctx sdk.Ctx) (res sdk.BigInt) {
	res = k.appKeeper.GetStakedTokens(ctx)
	return
}

// "GetNodeStakedTokens" - Returns the total number of staked tokens in the nodes module
func (k Keeper) GetNodesStakedTokens(ctx sdk.Ctx) (res sdk.BigInt) {
	res = k.posKeeper.GetStakedTokens(ctx)
	return
}

// "GetTotalTokens" - Returns the total number of tokens kept in any/all modules
func (k Keeper) GetTotalTokens(ctx sdk.Ctx) (res sdk.BigInt) {
	res = k.posKeeper.TotalTokens(ctx)
	return
}

// "GetTotalStakedTokens" - Returns the summation of app staked tokens and node staked tokens
func (k Keeper) GetTotalStakedTokens(ctx sdk.Ctx) (res sdk.BigInt) {
	res = k.GetNodesStakedTokens(ctx).Add(k.GetAppStakedTokens(ctx))
	return
}
