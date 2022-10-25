package keeper

import (
	"fmt"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// RewardForRelays - Award coins to an address (will be called at the beginning of the next block)
func (k Keeper) RewardForRelays(ctx sdk.Ctx, relays sdk.BigInt, address sdk.Address) sdk.BigInt {
	// feature flags
	isAfterRSCAL := k.Cdc.IsAfterNamedFeatureActivationHeight(ctx.BlockHeight(), codec.RSCALKey)
	isNonCustodialActive := k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight())

	// The conditions of the original non-custodial issue:
	isDuringFirstNonCustodialIssue := isNonCustodialActive && isAfterRSCAL && ctx.BlockHeight() <= codec.NonCustodial1RollbackHeight

	// get validator
	validator, found := k.GetValidator(ctx, address)
	if !found {
		ctx.Logger().Error(fmt.Errorf("no validator found for address %s; at height %d\n", address.String(), ctx.BlockHeight()).Error())
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
			ctx.Logger().Error(fmt.Errorf("no validator found for address %s; at height %d\n", address.String(), ctx.BlockHeight()).Error())
			return sdk.ZeroInt()
		}
	}

	var coins sdk.BigInt

	//check if PIP22 is enabled, if so scale the rewards
	if isAfterRSCAL {
		stake := validator.GetTokens()
		//floorstake to the lowest bin multiple or take ceiling, whicherver is smaller
		flooredStake := sdk.MinInt(stake.Sub(stake.Mod(k.ServicerStakeFloorMultiplier(ctx))), k.ServicerStakeWeightCeiling(ctx).Sub(k.ServicerStakeWeightCeiling(ctx).Mod(k.ServicerStakeFloorMultiplier(ctx))))
		//Convert from tokens to a BIN number
		bin := flooredStake.Quo(k.ServicerStakeFloorMultiplier(ctx))
		//calculate the weight value, weight will be a floatng point number so cast to DEC here and then truncate back to big int
		weight := bin.ToDec().FracPow(k.ServicerStakeFloorMultiplierExponent(ctx), Pip22ExponentDenominator).Quo(k.ServicerStakeWeightMultiplier(ctx))
		coinsDecimal := k.RelaysToTokensMultiplier(ctx).ToDec().Mul(relays.ToDec()).Mul(weight)
		//truncate back to int
		coins = coinsDecimal.TruncateInt()
	} else {
		coins = k.RelaysToTokensMultiplier(ctx).Mul(relays)
	}

	toNode, toFeeCollector := k.NodeReward(ctx, coins)
	if toNode.IsPositive() {
		k.mint(ctx, toNode, address)
	}
	if toFeeCollector.IsPositive() {
		k.mint(ctx, toFeeCollector, k.getFeePool(ctx).GetAddress())
	}
	return toNode
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
