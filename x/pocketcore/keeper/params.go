package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
)

// this file contains getters for all pocket core params
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) SessionNodeCount(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeySessionNodeCount, &res)
	return
}

func (k Keeper) ClaimExpiration(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyClaimExpiration, &res)
	return
}

func (k Keeper) SessionFrequency(ctx sdk.Context) int64 {
	frequency := k.posKeeper.SessionBlockFrequency(ctx)
	return frequency
}

func (k Keeper) ProofWaitingPeriod(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyProofWaitingPeriod, &res)
	return
}

func (k Keeper) SupportedBlockchains(ctx sdk.Context) (res []string) {
	k.Paramstore.Get(ctx, types.KeySupportedBlockchains, &res)
	return
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	params := types.Params{
		SessionNodeCount:     k.SessionNodeCount(ctx),
		ProofWaitingPeriod:   k.ProofWaitingPeriod(ctx),
		SupportedBlockchains: k.SupportedBlockchains(ctx),
		ClaimExpiration:      k.ClaimExpiration(ctx),
	}
	return params
}

// set the params object all at once
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	ctx.Logger().Info(fmt.Sprintf("SetParams(Params = %v) \n", params))
	k.Paramstore.SetParamSet(ctx, &params)
}
