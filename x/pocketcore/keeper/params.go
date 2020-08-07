package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "ParamKeyTable" - Registers the paramset types in a keytable and returns the table
func ParamKeyTable() sdk.KeyTable {
	return sdk.NewKeyTable().RegisterParamSet(&types.Params{})
}

// "SessionNodeCount" - Returns the session node count parameter from the paramstore
// Number of nodes dispatched in a single session
func (k Keeper) SessionNodeCount(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeySessionNodeCount, &res)
	return
}

// "ClaimExpiration" - Returns the claim expiration parameter from the paramstore
// Number of sessions pass before claim is expired
func (k Keeper) ClaimExpiration(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyClaimExpiration, &res)
	return
}

// "ReplayAttackBurn" - Returns the replay attack burn parameter from the paramstore
// The multiplier for how heavily nodes are burned for replay attacks
func (k Keeper) ReplayAttackBurnMultiplier(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyReplayAttackBurnMultiplier, &res)
	return
}

// "BlocksPerSession" - Returns blocksPerSession parameter from the paramstore
// How many blocks per session
func (k Keeper) BlocksPerSession(ctx sdk.Ctx) int64 {
	frequency := k.posKeeper.BlocksPerSession(ctx)
	return frequency
}

// "ClaimSubmissionWindow" - Returns claimSubmissionWindow parameter from the paramstore
// How long do you have to submit a claim before the secret is revealed and it's invalid
func (k Keeper) ClaimSubmissionWindow(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyClaimSubmissionWindow, &res)
	return
}

// "SupportedBlockchains" - Returns a supported blockchain parameter from the paramstore
// What blockchains are supported in pocket network (list of network identifier hashes)
func (k Keeper) SupportedBlockchains(ctx sdk.Ctx) (res []string) {
	k.Paramstore.Get(ctx, types.KeySupportedBlockchains, &res)
	return
}

// "SupportedBlockchains" - Returns a supported blockchain parameter from the paramstore
// What blockchains are supported in pocket network (list of network identifier hashes)
func (k Keeper) MinimumNumberOfProofs(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyMinimumNumberOfProofs, &res)
	return
}

// "GetParams" - Returns all module parameters in a `Params` struct
func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		SessionNodeCount:           k.SessionNodeCount(ctx),
		ClaimSubmissionWindow:      k.ClaimSubmissionWindow(ctx),
		SupportedBlockchains:       k.SupportedBlockchains(ctx),
		ClaimExpiration:            k.ClaimExpiration(ctx),
		ReplayAttackBurnMultiplier: k.ReplayAttackBurnMultiplier(ctx),
		MinimumNumberOfProofs:      k.MinimumNumberOfProofs(ctx),
	}
}

// "SetParams" - Sets all of the parameters in the paramstore using the params structure
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
