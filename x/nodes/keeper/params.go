package keeper

import (
	"time"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
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
func (k Keeper) MaxValidators(ctx sdk.Ctx) (res int64) {
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
	var minSignedPerWindow sdk.BigDec
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
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Ctx) (res sdk.BigDec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime - Retrieve slash fraction time
func (k Keeper) SlashFractionDowntime(ctx sdk.Ctx) (res sdk.BigDec) {
	k.Paramstore.Get(ctx, types.KeySlashFractionDowntime, &res)
	return
}

// RelaysToTokensMultiplier - Retrieve relay token multipler
func (k Keeper) RelaysToTokensMultiplier(ctx sdk.Ctx) sdk.BigInt {
	var multiplier int64
	k.Paramstore.Get(ctx, types.KeyRelaysToTokensMultiplier, &multiplier)
	return sdk.NewInt(multiplier)
}

func (k Keeper) NodeReward(ctx sdk.Ctx, reward sdk.BigInt) (nodeReward sdk.BigInt, feesCollected sdk.BigInt) {
	// convert reward to dec
	r := reward.ToDec()
	// get the dao and proposer % ex DAO .1 or 10% Proposer .01 or 1%
	daoAllocationPercentage := sdk.NewDec(k.DAOAllocation(ctx)).QuoInt64(int64(100))           // dec percentage
	proposerAllocationPercentage := sdk.NewDec(k.ProposerAllocation(ctx)).QuoInt64(int64(100)) // dec percentage
	// the dao and proposer allocations go to the fee collector
	daoAllocation := r.Mul(daoAllocationPercentage)
	proposerAllocation := r.Mul(proposerAllocationPercentage)
	// truncate int ex 1.99 uPOKT goes to 1 uPOKT
	feesCollected = daoAllocation.Add(proposerAllocation).TruncateInt()
	// the rest goes to the node
	nodeReward = reward.Sub(feesCollected)
	return
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
func (k Keeper) MaxChains(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyMaxChains, &res)
	return
}
func (k Keeper) MaxJailedBlocks(ctx sdk.Ctx) (res int64) {
	k.Paramstore.Get(ctx, types.KeyMaxJailedBlocks, &res)
	return
}

// GetParams - Retrieve all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Ctx) types.Params {
	return types.Params{
		RelaysToTokensMultiplier: k.RelaysToTokensMultiplier(ctx).Int64(),
		UnstakingTime:            k.UnStakingTime(ctx),
		MaxValidators:            k.MaxValidators(ctx),
		StakeDenom:               k.StakeDenom(ctx),
		StakeMinimum:             k.MinimumStake(ctx),
		SessionBlockFrequency:    k.BlocksPerSession(ctx),
		DAOAllocation:            k.DAOAllocation(ctx),
		ProposerAllocation:       k.ProposerAllocation(ctx),
		MaximumChains:            k.MaxChains(ctx),
		MaxJailedBlocks:          k.MaxJailedBlocks(ctx),
		MaxEvidenceAge:           k.MaxEvidenceAge(ctx),
		SignedBlocksWindow:       k.SignedBlocksWindow(ctx),
		MinSignedPerWindow:       sdk.NewDec(k.MinSignedPerWindow(ctx)),
		DowntimeJailDuration:     k.DowntimeJailDuration(ctx),
		SlashFractionDoubleSign:  k.SlashFractionDoubleSign(ctx),
		SlashFractionDowntime:    k.SlashFractionDowntime(ctx),
	}
}

// SetParams - Apply set of params
func (k Keeper) SetParams(ctx sdk.Ctx, params types.Params) {
	k.Paramstore.SetParamSet(ctx, &params)
}
