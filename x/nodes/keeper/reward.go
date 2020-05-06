package keeper

import (
	"fmt"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
)

// RewardForRelays - Award coins to an address (will be called at the beginning of the next block)
func (k Keeper) RewardForRelays(ctx sdk.Ctx, relays sdk.Int, address sdk.Address) {
	coins := k.RelaysToTokensMultiplier(ctx).Mul(relays)
	toNode, toFeeCollector := k.NodeCutOfReward(ctx, coins)
	k.mint(ctx, toNode, address)
	k.mint(ctx, toFeeCollector, k.getFeePool(ctx).GetAddress())
}

// blockReward - Handles distribution of the collected fees
func (k Keeper) blockReward(ctx sdk.Ctx, previousProposer sdk.Address) {
	feesCollector := k.getFeePool(ctx)
	feesCollected := feesCollector.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	// get the dao and proposer % ex DAO .1 or 10% Proposer .01 or 1%
	daoAllocation := sdk.NewDec(k.DAOAllocation(ctx))
	proposerAllocation := sdk.NewDec(k.ProposerAllocation(ctx))
	totalAllocationWithoutNodeCut := daoAllocation.Add(proposerAllocation)
	// get the new percentages based on the total. This is needed because the node (relayer) cut has already been allocated
	daoAllocation = daoAllocation.Quo(totalAllocationWithoutNodeCut)
	// dao cut calculation truncates int ex: 1.99uPOKT = 1uPOKT
	daoCut := feesCollected.ToDec().Mul(daoAllocation).TruncateInt()
	// proposer is whatever is left
	proposerCut := feesCollected.Sub(daoCut)
	// send to the two parties
	feeAddr := feesCollector.GetAddress()
	k.AccountKeeper.SendCoinsFromAccountToModule(ctx, feeAddr, "", sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, proposerCut)))
	k.AccountKeeper.SendCoins(ctx, feeAddr, previousProposer, sdk.NewCoins(sdk.NewCoin(proposerCut, sdk.DefaultStakeDenom)))
}

// GetTotalCustomValidatorAwards - Retrieve Custom Validator awards
func (k Keeper) GetTotalCustomValidatorAwards(ctx sdk.Ctx) sdk.Int {
	total := sdk.ZeroInt()
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AwardValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.Int{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &amount)
		total = total.Add(amount)
	}
	return total
}

// setValidatorAward - Store functions used to keep track of a validator award
func (k Keeper) setValidatorAward(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyForValidatorAward(address)
	val := amino.MustMarshalBinaryBare(amount)
	store.Set(key, val)
}

// getValidatorAward - Retrieve validator award
func (k Keeper) getValidatorAward(ctx sdk.Ctx, address sdk.Address) (coins sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorAward(address))
	if value == nil {
		return sdk.ZeroInt(), false
	}
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	found = true
	return coins, true
}

// deleteValidatorAward - Remove vallidaor award
func (k Keeper) deleteValidatorAward(ctx sdk.Ctx, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorAward(address))
}

// "mint" - takes an amount and mints it to the node staking pool, then sends the coins to the address
func (k Keeper) mint(ctx sdk.Ctx, amount sdk.Int, address sdk.Address) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.AccountKeeper.MintCoins(ctx, types.StakedPoolName, coins)
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := k.AccountKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, address, coins)
	if sendErr != nil {
		return sendErr.Result()
	}
	logString := fmt.Sprintf("a reward of %s was minted to %s", amount.String(), address.String())
	k.Logger(ctx).Info(logString)
	return sdk.Result{
		Log: logString,
	}
}

// GetPreviousProposer - Retrieve the proposer public key for this block
func (k Keeper) GetPreviousProposer(ctx sdk.Ctx) (consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return
}

// SetPreviousProposer -  Store proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Ctx, consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(types.ProposerKey, b)
}

// getProposerAllocation - Retrieve proposer allocation
func (k Keeper) getProposerAllocaiton(ctx sdk.Ctx) sdk.Int {
	return sdk.NewInt(k.ProposerAllocation(ctx))
}

// getDAOAllocation - Retrieve DAO allocation
func (k Keeper) getDAOAllocation(ctx sdk.Ctx) sdk.Int {
	return sdk.NewInt(k.DAOAllocation(ctx))
}
