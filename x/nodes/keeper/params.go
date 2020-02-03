package keeper

import (
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/types"
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

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Context) (res uint64) {
	k.Paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// StakeDenom - Bondable coin denomination
func (k Keeper) StakeDenom(ctx sdk.Context) (res string) {
	k.Paramstore.Get(ctx, types.KeyStakeDenom, &res)
	return
}

func (k Keeper) MinimumStake(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyStakeMinimum, &res)
	return
}

func (k Keeper) ProposerAllocation(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyProposerAllocation, &res)
	return
}

// MaxEvidenceAge - max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

// Downtime slashing threshold
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) (res int64) {
	var minSignedPerWindow sdk.Dec
	k.Paramstore.Get(ctx, types.KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64() // todo may have to be int64 .RoundInt64()
}

// Downtime jail duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Context) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

func (k Keeper) RelaysToTokensMultiplier(ctx sdk.Context) sdk.Int {
	daoAllocation := k.DAOAllocation(ctx)
	proposerAllocation := k.ProposerAllocation(ctx)
	return sdk.NewInt(100 - daoAllocation + proposerAllocation)
}

func (k Keeper) DAOAllocation(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeyDAOAllocation, &res)
	return
}

func (k Keeper) SessionBlockFrequency(ctx sdk.Context) (res int64) {
	k.Paramstore.Get(ctx, types.KeySessionBlock, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.Params{
		UnstakingTime:           k.UnStakingTime(ctx),
		MaxValidators:           k.MaxValidators(ctx),
		StakeDenom:              k.StakeDenom(ctx),
		StakeMinimum:            k.MinimumStake(ctx),
		ProposerAllocation:      k.ProposerAllocation(ctx),
		SessionBlockFrequency:   k.SessionBlockFrequency(ctx),
		DAOAllocation:           k.DAOAllocation(ctx),
		MaxEvidenceAge:          k.MaxEvidenceAge(ctx),
		SignedBlocksWindow:      k.SignedBlocksWindow(ctx),
		MinSignedPerWindow:      sdk.NewDec(k.MinSignedPerWindow(ctx)),
		DowntimeJailDuration:    k.DowntimeJailDuration(ctx),
		SlashFractionDoubleSign: k.SlashFractionDoubleSign(ctx),
		SlashFractionDowntime:   k.SlashFractionDowntime(ctx),
	}
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
