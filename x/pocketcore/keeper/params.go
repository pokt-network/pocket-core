package keeper

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) SessionNodeCount(ctx sdk.Context) (res uint) {
	k.Paramstore.Get(ctx, types.KeySessionNodeCount, &res)
	return
}

func (k Keeper) SessionFrequency(ctx sdk.Context) int64 {
	return k.posKeeper.SessionBlockFrequency(ctx)
}

func (k Keeper) ProofWaitingPeriod(ctx sdk.Context) (res uint) {
	k.Paramstore.Get(ctx, types.KeyProofWaitingPeriod, &res)
	return
}

func (k Keeper) SupportedBlockchains(ctx sdk.Context) (res []string) {
	k.Paramstore.Get(ctx, types.KeySupportedBlockchains, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		SessionNodeCount:     k.SessionNodeCount(ctx),
		ProofWaitingPeriod:   k.ProofWaitingPeriod(ctx),
		SupportedBlockchains: k.SupportedBlockchains(ctx),
	}
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
