package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
)

// this file contains getters for all pocket core params
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) SessionNodeCount(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeySessionNodeCount, &res)
	return
}

func (k Keeper) ClaimExpiration(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyClaimExpiration, &res)
	return
}

func (k Keeper) SessionFrequency(ctx sdk.Ctx) int64 {
	frequency := k.posKeeper.SessionBlockFrequency(ctx)
	return frequency
}

func (k Keeper) ClaimSubmissionWindow(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyClaimSubmissionWindow, &res)
	return
}

func (k Keeper) SupportedBlockchains(ctx sdk.Ctx) (res []string) {
	k.Paramstore.Get(ctx, types.KeySupportedBlockchains, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		SessionNodeCount:      k.SessionNodeCount(ctx),
		ClaimSubmissionWindow: k.ClaimSubmissionWindow(ctx),
		SupportedBlockchains:  k.SupportedBlockchains(ctx),
		ClaimExpiration:       k.ClaimExpiration(ctx),
	}
}

// set the params object all at once
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
