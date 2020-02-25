package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/tendermint/go-amino"
)

// award coins to an address (will be called at the beginning of the next block)
func (k Keeper) AwardCoinsTo(ctx sdk.Context, relays sdk.Int, address sdk.Address) {
	award, _ := k.getValidatorAward(ctx, address)
	coins := k.RelaysToTokensMultiplier(ctx).Mul(relays)
	k.setValidatorAward(ctx, award.Add(coins), address)
	ctx.Logger().Info("Custom award of " + coins.String() + " set for " + address.String())
}

// blockReward handles distribution of the collected fees
func (k Keeper) blockReward(ctx sdk.Context, previousProposer sdk.Address) {
	logger := k.Logger(ctx)
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.getFeePool(ctx)
	feesCollected := feeCollector.GetCoins()
	// transfer collected fees to the pos module account
	err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, auth.FeeCollectorName, types.ModuleName, feesCollected)
	if err != nil {
		panic(err)
	}
	// get the reward from the total relays completed in the last block
	rewardForRelays := k.GetTotalCustomValidatorAwards(ctx)
	// calculate the total reward by adding relays to the fees
	totalReward := feesCollected.AmountOf(k.StakeDenom(ctx)).Add(rewardForRelays)
	// calculate previous proposer reward
	proposerAllocation := k.getProposerAllocaiton(ctx)
	daoAllocation := k.getDAOAllocation(ctx)
	// divide up the reward from the proposer reward and the dao reward
	proposerReward := proposerAllocation.Mul(totalReward).Quo(sdk.NewInt(100)) // truncates
	daoReward := daoAllocation.Mul(totalReward).Quo(sdk.NewInt(100))           // truncates
	// get the validator structure
	proposerValidator, found := k.GetValidator(ctx, previousProposer)
	if found {
		// create the proposer coins
		propRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), proposerReward))
		// create the dao coins
		daoRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), daoReward))
		// mint the coins for the proposer to the module account
		err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, propRewardCoins)
		if err != nil {
			panic(err)
		}
		// mint the coins for the dao to the module account
		err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, daoRewardCoins)
		if err != nil {
			panic(err)
		}
		// send to validator
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, proposerValidator.GetAddress(), propRewardCoins); err != nil {
			panic(err)
		}
		// send to rest dao
		if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.DAOPoolName, daoRewardCoins); err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("minted %s to block proposer: %s", propRewardCoins.String(), proposerValidator.GetAddress().String()))
		logger.Info(fmt.Sprintf("minted %s to DAO", daoRewardCoins.String()))
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposerReward,
				sdk.NewAttribute(sdk.AttributeKeyAmount, proposerReward.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, proposerValidator.GetAddress().String()),
			),
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDAOAllocation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, daoReward.String()),
			),
		)
	} else {
		logger.Error(fmt.Sprintf(
			"WARNING: Attempt to allocate proposer rewards to unknown proposer %s. "+
				"This should happen only if the proposer unstaked completely within a single block, "+
				"which generally should not happen except in exceptional circumstances (or fuzz testing). "+
				"We recommend you investigate immediately.",
			previousProposer.String()))
	}
}

func (k Keeper) GetTotalCustomValidatorAwards(ctx sdk.Context) sdk.Int {
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

// store functions used to keep track of a validator award
func (k Keeper) setValidatorAward(ctx sdk.Context, amount sdk.Int, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyForValidatorAward(address)
	val := amino.MustMarshalBinaryBare(amount)
	store.Set(key, val)
}

func (k Keeper) getValidatorAward(ctx sdk.Context, address sdk.Address) (coins sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorAward(address))
	if value == nil {
		return sdk.ZeroInt(), false
	}
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	found = true
	return coins, true
}

func (k Keeper) deleteValidatorAward(ctx sdk.Context, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorAward(address))
}

// called on begin blocker
func (k Keeper) mintNodeRelayRewards(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AwardValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.Int{}
		address := sdk.Address(types.AddressFromKey(iterator.Key()))
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &amount)
		amount = k.NodeCutOfReward(ctx).Mul(amount).Quo(sdk.NewInt(100)) // truncate
		k.mint(ctx, amount, address)
		// remove from the award store
		store.Delete(iterator.Key())
		ctx.Logger().Info("Relay reward of " + amount.String() + " minted to" + address.String())
	}
}

// Mints sdk.Coins and sends them to an address
func (k Keeper) mint(ctx sdk.Context, amount sdk.Int, address sdk.Address) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.supplyKeeper.MintCoins(ctx, types.StakedPoolName, coins)
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.StakedPoolName, address, coins)
	if sendErr != nil {
		return sendErr.Result()
	}
	logString := fmt.Sprintf("a custom reward of %s was minted to %s", amount.String(), address.String())
	k.Logger(ctx).Info(logString)
	return sdk.Result{
		Log: logString,
	}
}

// get the proposer public key for this block
func (k Keeper) GetPreviousProposer(ctx sdk.Context) (consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Context, consAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(types.ProposerKey, b)
}

func (k Keeper) getProposerAllocaiton(ctx sdk.Context) sdk.Int {
	return sdk.NewInt(k.ProposerAllocation(ctx))
}

func (k Keeper) getDAOAllocation(ctx sdk.Context) sdk.Int {
	return sdk.NewInt(k.DAOAllocation(ctx))
}
