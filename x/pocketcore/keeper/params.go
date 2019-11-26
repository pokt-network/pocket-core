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
