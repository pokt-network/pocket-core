package keeper

import (
	"fmt"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	pcTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// RewardForRelays - Award coins to an address using the default multiplier
func (k Keeper) RewardForRelays(ctx sdk.Ctx, relays sdk.BigInt, address sdk.Address) sdk.BigInt {
	return k.RewardForRelaysPerChain(ctx, "", relays, address)
}

// CalculateRelayReward - Calculate the amount of rewards based on the given
// number of relays and the staked tokens.  The returned reward amount includes
// DAO & Proposer portions.
func (k Keeper) CalculateRelayReward(
	ctx sdk.Ctx,
	chain string,
	relays sdk.BigInt,
	stake sdk.BigInt,
) (sdk.BigInt, sdk.BigInt) {
	// feature flags
	isAfterRSCAL := k.Cdc.IsAfterNamedFeatureActivationHeight(ctx.BlockHeight(), codec.RSCALKey)
	multiplier := k.GetChainSpecificMultiplier(ctx, chain)

	var coins sdk.BigInt
	//check if PIP22 is enabled, if so scale the rewards
	if isAfterRSCAL {
		//floorstake to the lowest bin multiple or take ceiling, whicherver is smaller
		flooredStake := sdk.MinInt(
			stake.Sub(stake.Mod(k.ServicerStakeFloorMultiplier(ctx))),
			k.ServicerStakeWeightCeiling(ctx).
				Sub(k.ServicerStakeWeightCeiling(ctx).Mod(k.ServicerStakeFloorMultiplier(ctx))),
		)
		//Convert from tokens to a BIN number
		bin := flooredStake.Quo(k.ServicerStakeFloorMultiplier(ctx))
		//calculate the weight value, weight will be a floatng point number so cast
		// to DEC here and then truncate back to big int
		weight := bin.ToDec().
			FracPow(k.ServicerStakeFloorMultiplierExponent(ctx), Pip22ExponentDenominator).
			Quo(k.ServicerStakeWeightMultiplier(ctx))
		coinsDecimal := multiplier.ToDec().Mul(relays.ToDec()).Mul(weight)
		//truncate back to int
		coins = coinsDecimal.TruncateInt()
	} else {
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

	if k.Cdc.IsAfterDelegatorUpgrade(ctx.BlockHeight()) {
		rewardCost := k.AccountKeeper.GetFee(ctx, pcTypes.MsgClaim{}).
			Add(k.AccountKeeper.GetFee(ctx, pcTypes.MsgProof{}))
		if toNode.LT(rewardCost) {
			rewardCost = toNode
		}
		k.mint(ctx, rewardCost, validator.Address)
		toNode = toNode.Sub(rewardCost)
	}

	SplitNodeRewards(
		toNode,
		address,
		validator.Delegators,
		func(recipient sdk.Address, share sdk.BigInt) {
			k.mint(ctx, share, recipient)
		},
	)
	if toFeeCollector.IsPositive() {
		k.mint(ctx, toFeeCollector, k.getFeePool(ctx).GetAddress())
	}
	return toNode
}

// Splits rewards into the primary recipient and delegator addresses and
// invokes a callback per share.
func SplitNodeRewards(
	rewards sdk.BigInt,
	primaryRecipient sdk.Address,
	delegators map[string]uint32,
	callback func(sdk.Address, sdk.BigInt),
) {
	if !rewards.IsPositive() {
		return
	}

	totalShare := int64(0)
	for _, share := range delegators {
		totalShare = totalShare + int64(share)
		if totalShare > 100 {
			// If the total shares for delegators exceeds 100,
			// all rewards go to the primary recipient.
			delegators = nil
			break
		}
	}

	remains := rewards
	for addrStr, share := range delegators {
		addr, err := sdk.AddressFromHex(addrStr)
		if err != nil {
			continue
		}
		percentage := sdk.NewDecWithPrec(int64(share), 2)
		allocation := rewards.ToDec().Mul(percentage).TruncateInt()
		if allocation.IsPositive() {
			callback(addr, allocation)
		}
		remains = remains.Sub(allocation)
	}

	if remains.IsPositive() {
		callback(primaryRecipient, remains)
	}
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
	// check for zero fees
	if feesCollected.IsZero() {
		return
	}
	// get the dao and proposer % ex DAO .1 or 10% Proposer .01 or 1%
	daoAllocation := sdk.NewDec(k.DAOAllocation(ctx))
	proposerAllocation := sdk.NewDec(k.ProposerAllocation(ctx))
	daoAndProposerAllocation := daoAllocation.Add(proposerAllocation)
	// get the new percentages based on the total. This is needed because the node (relayer) cut has already been allocated
	daoAllocation = daoAllocation.Quo(daoAndProposerAllocation)
	// dao cut calculation truncates int ex: 1.99uPOKT = 1uPOKT
	daoCut := feesCollected.ToDec().Mul(daoAllocation).TruncateInt()
	// proposer is whatever is left
	proposerCut := feesCollected.Sub(daoCut)
	// send to the two parties
	feeAddr := feesCollector.GetAddress()
	err := k.AccountKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, govTypes.DAOAccountName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, daoCut)))
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("unable to send %s cut of block reward to the dao: %s, at height %d", daoCut.String(), err.Error(), ctx.BlockHeight()))
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		outputAddress, found := k.GetValidatorOutputAddress(ctx, previousProposer)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("unable to send %s cut of block reward to the proposer: %s, with error %s, at height %d", proposerCut.String(), previousProposer, types.ErrNoValidatorForAddress(types.ModuleName), ctx.BlockHeight()))
			return
		}
		err = k.AccountKeeper.SendCoins(ctx, feeAddr, outputAddress, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, proposerCut)))
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("unable to send %s cut of block reward to the proposer: %s, with error %s, at height %d", proposerCut.String(), previousProposer, err.Error(), ctx.BlockHeight()))
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
