package keeper

import (
	"time"

	"github.com/pokt-network/pocket-core/x/apps/types"
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

// UnstakingTime
func (k Keeper) UnStakingTime(ctx sdk.Context) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyUnstakingTime, &res)
	return
}

func (k Keeper) RelayCoefficient(ctx sdk.Context) (coefficient uint8) {
	k.Paramstore.Get(ctx, types.KeyRelayCoefficientPercentage, &coefficient)
	return
}

// MaxApplications - Maximum number of applications
func (k Keeper) MaxApplications(ctx sdk.Context) (res uint64) {
	k.Paramstore.Get(ctx, types.KeyMaxApplications, &res)
	return
}

func (k Keeper) MinimumStake(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyApplicationMinStake, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		UnstakingTime:              k.UnStakingTime(ctx),
		MaxApplications:            k.MaxApplications(ctx),
		AppStakeMin:                k.MinimumStake(ctx),
		RelayCoefficientPercentage: k.RelayCoefficient(ctx),
	}
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
