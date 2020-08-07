package keeper

import (
	"fmt"
	"os"

	sdk "github.com/pokt-network/pocket-core/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// RewardForRelays - Award coins to an address (will be called at the beginning of the next block)
func (k Keeper) RewardForRelays(ctx sdk.Ctx, relays sdk.Int, address sdk.Address) sdk.Int {
	coins := k.RelaysToTokensMultiplier(ctx).Mul(relays)
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
	err = k.AccountKeeper.SendCoins(ctx, feeAddr, previousProposer, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, proposerCut)))
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("unable to send %s cut of block reward to the proposer: %s, at height %d", proposerCut.String(), err.Error(), ctx.BlockHeight()))
	}
}

// "mint" - takes an amount and mints it to the node staking pool, then sends the coins to the address
func (k Keeper) mint(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) sdk.Result {
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
	b := store.Get(types.ProposerKey)
	if b == nil {
		k.Logger(ctx).Error("Previous proposer not set")
		os.Exit(1)
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &addr)
	return
}

// SetPreviousProposer -  Store proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Ctx, consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(types.ProposerKey, b)
}
