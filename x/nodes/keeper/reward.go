package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	pcTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/libs/log"
)

// GetRewardCost - The cost a servicer needs to pay to earn relay rewards
func (k Keeper) GetRewardCost(ctx sdk.Ctx) sdk.BigInt {
	return k.AccountKeeper.GetFee(ctx, pcTypes.MsgClaim{}).
		Add(k.AccountKeeper.GetFee(ctx, pcTypes.MsgProof{}))
}

// RewardForRelays - Award coins to an address using the default multiplier
func (k Keeper) RewardForRelays(ctx sdk.Ctx, relays sdk.BigInt, address sdk.Address) sdk.BigInt {
	return k.RewardForRelaysPerChain(ctx, "", relays, address)
}

func (k Keeper) calculateRewardRewardPip22(
	ctx sdk.Ctx,
	relays, stake, multiplier sdk.BigInt,
) sdk.BigInt {
	// floorstake to the lowest bin multiple or take ceiling, whichever is smaller
	flooredStake := sdk.MinInt(
		stake.Sub(stake.Mod(k.ServicerStakeFloorMultiplier(ctx))),
		k.ServicerStakeWeightCeiling(ctx).
			Sub(k.ServicerStakeWeightCeiling(ctx).
				Mod(k.ServicerStakeFloorMultiplier(ctx))),
	)
	// Convert from tokens to a BIN number
	bin := flooredStake.Quo(k.ServicerStakeFloorMultiplier(ctx))
	// calculate the weight value, weight will be a floatng point number so cast
	// to DEC here and then truncate back to big int
	weight := bin.ToDec().
		FracPow(
			k.ServicerStakeFloorMultiplierExponent(ctx),
			Pip22ExponentDenominator,
		).
		Quo(k.ServicerStakeWeightMultiplier(ctx))
	coinsDecimal := multiplier.ToDec().Mul(relays.ToDec()).Mul(weight)
	// truncate back to int
	return coinsDecimal.TruncateInt()
}

// CalculateRelayReward - Calculates the amount of rewards based on the given
// number of relays and the staked tokens, and splits it to the servicer's cut
// and the DAO & Proposer cut.
func (k Keeper) CalculateRelayReward(
	ctx sdk.Ctx,
	chain string,
	relays sdk.BigInt,
	stake sdk.BigInt,
) (nodeReward, feesCollected sdk.BigInt) {
	isAfterRSCAL := k.Cdc.IsAfterNamedFeatureActivationHeight(
		ctx.BlockHeight(),
		codec.RSCALKey,
	)
	multiplier := k.GetChainSpecificMultiplier(ctx, chain)

	var coins sdk.BigInt
	if isAfterRSCAL {
		// scale the rewards if PIP22 is enabled
		coins = k.calculateRewardRewardPip22(ctx, relays, stake, multiplier)
	} else {
		// otherwise just apply rttm
		coins = multiplier.Mul(relays)
	}

	return k.splitRewards(ctx, coins)
}

// RewardForRelaysPerChain - Award coins to an address for relays of a specific chain
func (k Keeper) RewardForRelaysPerChain(ctx sdk.Ctx, chain string, relays sdk.BigInt, address sdk.Address) sdk.BigInt {
	// feature flags
	isAfterRSCAL := k.Cdc.IsAfterNamedFeatureActivationHeight(ctx.BlockHeight(), codec.RSCALKey)
	isNonCustodialActive := k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight())

	// The conditions of the original non-custodial issue:
	isDuringFirstNonCustodialIssue := isNonCustodialActive && isAfterRSCAL && ctx.BlockHeight() <= codec.NonCustodial1RollbackHeight

	// get validator
	validator, found := k.GetValidator(ctx, address)

	//adding "&& (isAfterRSCAL || isAfterNonCustodial)" to sync from scratch as weighted stake and non-custodial introduced this requirement
	if !found && (isAfterRSCAL || isNonCustodialActive) {
		ctx.Logger().Error(
			"no validator found",
			"address", address,
			"height", ctx.BlockHeight(),
		)
		return sdk.ZeroInt()
	}

	if isNonCustodialActive {
		// if during non-custodial rollback issue on height 69583, use the output address
		// else only use the output address if after noncustodial allowance height 74622
		if isDuringFirstNonCustodialIssue || ctx.BlockHeight() >= codec.NonCustodial2AllowanceHeight {
			address = k.GetOutputAddressFromValidator(validator)
		}
	}

	// This simulates the issue that happened between the original non-custodial rollout and the rollback height
	// This code is needed to allow a 'history replay' of the issue during sync-from-scratch
	if isDuringFirstNonCustodialIssue {
		_, found := k.GetValidator(ctx, address)
		if !found {
			ctx.Logger().Error(
				"no validator found",
				"address", address,
				"height", ctx.BlockHeight(),
			)
			return sdk.ZeroInt()
		}
	}

	toNode, toFeeCollector :=
		k.CalculateRelayReward(ctx, chain, relays, validator.GetTokens())

	// After the delegator upgrade, we compensate a servicer's operator wallet
	// for the transaction fee of claim and proof.
	if k.Cdc.IsAfterRewardDelegatorUpgrade(ctx.BlockHeight()) {
		rewardCost := k.GetRewardCost(ctx)
		if toNode.LT(rewardCost) {
			// If the servicer's portion is less than the reward cost, we send
			// all of the servicer's portion to the servicer and no reward is sent
			// to the output address or delegators.  This case causes a net loss to
			// the servicer.   If prevent_negative_reward_claim is set to true,
			// a servicer will not claim for tiny evidences that cause a net loss.
			rewardCost = toNode
		}
		if rewardCost.IsPositive() {
			k.mint(ctx, rewardCost, validator.Address)
			toNode = toNode.Sub(rewardCost)
		}
	}

	err := SplitNodeRewards(
		ctx.Logger(),
		toNode,
		address,
		validator.RewardDelegators,
		func(recipient sdk.Address, share sdk.BigInt) {
			k.mint(ctx, share, recipient)
		},
	)
	if err != nil {
		ctx.Logger().Error("unable to split relay rewards",
			"height", ctx.BlockHeight(),
			"servicer", validator.Address,
			"err", err.Error(),
		)
	}

	if toFeeCollector.IsPositive() {
		k.mint(ctx, toFeeCollector, k.getFeePool(ctx).GetAddress())
	}
	return toNode
}

// Splits rewards into the primary recipient and delegator addresses and
// invokes a callback per share.
// delegators - a map from address to its share (< 100)
// shareRewardsCallback - a callback to send `coins` of total rewards to `addr`
func SplitNodeRewards(
	logger log.Logger,
	rewards sdk.BigInt,
	primaryRecipient sdk.Address,
	delegators map[string]uint32,
	shareRewardsCallback func(addr sdk.Address, coins sdk.BigInt),
) error {
	if !rewards.IsPositive() {
		return errors.New("non-positive rewards")
	}

	normalizedDelegators, err := types.NormalizeRewardDelegators(delegators)
	if err != nil {
		// If the delegators field is invalid, do nothing.
		return errors.New("invalid delegators")
	}

	remains := rewards
	for _, pair := range normalizedDelegators {
		percentage := sdk.NewDecWithPrec(int64(pair.RewardShare), 2)
		allocation := rewards.ToDec().Mul(percentage).TruncateInt()
		if allocation.IsPositive() {
			shareRewardsCallback(pair.Address, allocation)
		}
		remains = remains.Sub(allocation)
	}

	if remains.IsPositive() {
		shareRewardsCallback(primaryRecipient, remains)
	} else {
		delegatorsBytes, _ := json.Marshal(delegators)
		logger.Error(
			"over-distributed rewards to delegators",
			"rewards", rewards,
			"remains", remains,
			"delegators", string(delegatorsBytes),
		)
	}
	return nil
}

// Calculates a chain-specific Relays-To-Token-Multiplier.
// Returns the default multiplier if the feature is not activated or a given
// chain is not set in the parameter.
func (k Keeper) GetChainSpecificMultiplier(ctx sdk.Ctx, chain string) sdk.BigInt {
	if k.Cdc.IsAfterPerChainRTTMUpgrade(ctx.BlockHeight()) {
		multiplierMap := k.RelaysToTokensMultiplierMap(ctx)
		if multiplier, found := multiplierMap[chain]; found {
			return sdk.NewInt(multiplier)
		}
	}
	return k.RelaysToTokensMultiplier(ctx)
}

// blockReward - Handles distribution of the collected fees
func (k Keeper) blockReward(ctx sdk.Ctx, previousProposer sdk.Address) {
	feesCollector := k.getFeePool(ctx)
	feesCollected := feesCollector.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	if feesCollected.IsZero() {
		return
	}

	daoCut, proposerCut := k.splitFeesCollected(ctx, feesCollected)

	// send to the two parties
	feeAddr := feesCollector.GetAddress()
	err := k.AccountKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, govTypes.DAOAccountName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, daoCut)))
	if err != nil {
		ctx.Logger().Error("unable to send a DAO cut of block reward",
			"height", ctx.BlockHeight(),
			"cut", daoCut,
			"err", err.Error(),
		)
	}

	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		validator, found := k.GetValidator(ctx, previousProposer)
		if !found {
			ctx.Logger().Error("unable to find a validator to send a block reward to",
				"height", ctx.BlockHeight(),
				"addr", previousProposer,
			)
			return
		}

		if !k.Cdc.IsAfterRewardDelegatorUpgrade(ctx.BlockHeight()) {
			validator.RewardDelegators = nil
		}

		err := SplitNodeRewards(
			ctx.Logger(),
			proposerCut,
			k.GetOutputAddressFromValidator(validator),
			validator.RewardDelegators,
			func(recipient sdk.Address, share sdk.BigInt) {
				err = k.AccountKeeper.SendCoins(
					ctx,
					feeAddr,
					recipient,
					sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, share)),
				)
				if err != nil {
					ctx.Logger().Error("unable to send a cut of block reward",
						"height", ctx.BlockHeight(),
						"cut", share,
						"addr", recipient,
						"err", err.Error(),
					)
				}
			},
		)
		if err != nil {
			ctx.Logger().Error("unable to split block rewards",
				"height", ctx.BlockHeight(),
				"validator", validator.Address,
				"err", err.Error(),
			)
		}

		return
	}

	err = k.AccountKeeper.SendCoins(ctx, feeAddr, previousProposer, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, proposerCut)))
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("unable to send %s cut of block reward to the proposer: %s, with error %s, at height %d", proposerCut.String(), previousProposer, err.Error(), ctx.BlockHeight()))
	}
}

// "mint" - takes an amount and mints it to the node staking pool, then sends the coins to the address
func (k Keeper) mint(ctx sdk.Ctx, amount sdk.BigInt, address sdk.Address) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.AccountKeeper.MintCoins(ctx, types.StakedPoolName, coins)
	if mintErr != nil {
		ctx.Logger().Error(fmt.Sprintf("unable to mint tokens, at height %d: ", ctx.BlockHeight()) + mintErr.Error())
		return mintErr.Result()
	}
	sendErr := k.AccountKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, address, coins)
	if sendErr != nil {
		ctx.Logger().Error(fmt.Sprintf("unable to send tokens, at height %d: ", ctx.BlockHeight()) + sendErr.Error())
		return sendErr.Result()
	}
	logString := fmt.Sprintf("a reward of %s was minted to %s", amount.String(), address.String())
	k.Logger(ctx).Info(logString)
	return sdk.Result{
		Log: logString,
	}
}

// GetPreviousProposer - Retrieve the proposer public key for this block
func (k Keeper) GetPreviousProposer(ctx sdk.Ctx) (addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b, _ := store.Get(types.ProposerKey)
	if b == nil {
		k.Logger(ctx).Error("Previous proposer not set")
		return nil
		//os.Exit(1)
	}
	_ = k.Cdc.UnmarshalBinaryLengthPrefixed(b, &addr, ctx.BlockHeight())
	return addr

}

// SetPreviousProposer -  Store proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Ctx, consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b, err := k.Cdc.MarshalBinaryLengthPrefixed(&consAddr, ctx.BlockHeight())
	if err != nil {
		panic(err)
	}
	_ = store.Set(types.ProposerKey, b)
}
