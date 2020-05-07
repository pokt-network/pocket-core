package keeper

import (
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamKeyTable for staking module
func ParamKeyTable() sdk.KeyTable {
	return sdk.NewKeyTable().RegisterParamSet(&types.Params{})
}

// UnStakingTime - Retrieve unstaking time param
func (k Keeper) UnStakingTime(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyUnstakingTime, &res)
	return
}

// MaxValidators - Retrieve maximum number of validators
func (k Keeper) MaxValidators(ctx sdk.Ctx) (res uint64) {
	k.Paramstore.Get(ctx, types.KeyMaxValidators, &res)
	return
}

// StakeDenom - Bondable coin denomination
func (k Keeper) StakeDenom(ctx sdk.Ctx) (res string) {
	k.Paramstore.Get(ctx, types.KeyStakeDenom, &res)
	return
}

// MinimumStake - Retrieve Minimum stake
func (k Keeper) MinimumStake(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyStakeMinimum, &res)
	return
}

// ProposerAllocation - Retrieve proposer allocation
func (k Keeper) ProposerAllocation(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyProposerAllocation, &res)
	return
}

// MaxEvidenceAge - Max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - Sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeySignedBlocksWindow, &res)
	return
}

// MinSignedPerWindow - Downtime slashing threshold
func (k Keeper) MinSignedPerWindow(ctx sdk.Ctx) (res int64) {
	var minSignedPerWindow sdk.Dec
	k.Paramstore.Get(ctx, types.KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64() // todo may have to be int64 .RoundInt64()
}

// DowntimeJailDuration - Downtime jail duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Ctx) (res time.Duration) {
	k.Paramstore.Get(ctx, types.KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign - Retrieve slash fraction for double signature
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Ctx) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime - Retrieve slash fraction time
func (k Keeper) SlashFractionDowntime(ctx sdk.Ctx) (res sdk.Dec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// RelaysToTokensMultiplier - Retrieve relay token multipler
func (k Keeper) RelaysToTokensMultiplier(ctx sdk.Ctx) sdk.Int {
	var multiplier int64
	k.Paramstore.Get(ctx, types.KeyRelaysToTokensMultiplier, &multiplier)
	return sdk.NewInt(multiplier)
}

// NodeCutOfReward - Retrieve node cut of reward
func (k Keeper) NodeCutOfReward(ctx sdk.Ctx) sdk.Int {
	daoAllocation := k.DAOAllocation(ctx)
	proposerAllocation := k.ProposerAllocation(ctx)
	return sdk.NewInt(100 - daoAllocation - proposerAllocation)
}

// DAOAllocation - Retrieve DAO allocation
func (k Keeper) DAOAllocation(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyDAOAllocation, &res)
	return
}

// BlocksPerSession - Retrieve blocks per session
func (k Keeper) BlocksPerSession(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeySessionBlock, &res)
	return
}

// GetParams - Retrieve all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		UnstakingTime:           k.UnStakingTime(ctx),
		MaxValidators:           k.MaxValidators(ctx),
		StakeDenom:              k.StakeDenom(ctx),
		StakeMinimum:            k.MinimumStake(ctx),
		ProposerAllocation:      k.ProposerAllocation(ctx),
		SessionBlockFrequency:   k.BlocksPerSession(ctx),
		DAOAllocation:           k.DAOAllocation(ctx),
		MaxEvidenceAge:          k.MaxEvidenceAge(ctx),
		SignedBlocksWindow:      k.SignedBlocksWindow(ctx),
		MinSignedPerWindow:      sdk.NewDec(k.MinSignedPerWindow(ctx)),
		DowntimeJailDuration:    k.DowntimeJailDuration(ctx),
		SlashFractionDoubleSign: k.SlashFractionDoubleSign(ctx),
		SlashFractionDowntime:   k.SlashFractionDowntime(ctx),
	}
}

// SetParams - Apply set of params
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
