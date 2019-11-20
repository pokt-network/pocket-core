package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/go-amino"
)

// award coins to an address (will be called at the beginning of the next block)
func (k Keeper) AwardCoinsTo(ctx sdk.Context, amount sdk.Int, address sdk.ValAddress) {
	award, _ := k.getValidatorAward(ctx, address)
	k.setValidatorAward(ctx, award.Add(amount), address)
}

// rewardFromFees handles distribution of the collected fees
func (k Keeper) rewardFromFees(ctx sdk.Context, previousProposer sdk.ConsAddress) {
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
	// calculate the total reward by adding relays to the
	totalReward := feesCollected.AmountOf(k.StakeDenom(ctx))
	// calculate previous proposer reward
	baseProposerRewardPercentage := k.getProposerRewardPercentage(ctx)
	// divide up the reward from the proposer reward and the dao reward
	proposerReward := baseProposerRewardPercentage.Mul(totalReward)
	daoReward := totalReward.Sub(proposerReward)
	// get the validator structure
	proposerValidator := k.validatorByConsAddr(ctx, previousProposer)
	if proposerValidator != nil {
		propRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), proposerReward))
		daoRewardCoins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), daoReward))
		// send to validator
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName,
			sdk.AccAddress(proposerValidator.GetAddress()), propRewardCoins); err != nil {
			panic(err)
		}
		// send to rest dao
		if err := k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.DAOPoolName, daoRewardCoins); err != nil {
			panic(err)
		}
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

// called on begin blocker
func (k Keeper) mintValidatorAwards(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AwardValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		amount := sdk.Int{}
		address := sdk.ValAddress{}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &amount)
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Key(), address)
		k.mint(ctx, amount, address)
		// remove from the award store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a validator award
func (k Keeper) setValidatorAward(ctx sdk.Context, amount sdk.Int, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorAward(address), amino.MustMarshalBinaryBare(amount))
}

func (k Keeper) getValidatorAward(ctx sdk.Context, address sdk.ValAddress) (coins sdk.Int, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorAward(address))
	if value == nil {
		return coins, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(value, &coins)
	return
}

func (k Keeper) deleteValidatorAward(ctx sdk.Context, address sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorAward(address))
}

// Mints sdk.Coins
func (k Keeper) mint(ctx sdk.Context, amount sdk.Int, address sdk.ValAddress) sdk.Result {
	coins := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	mintErr := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins.Add(coins))
	if mintErr != nil {
		return mintErr.Result()
	}
	sendErr := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address), coins)
	if sendErr != nil {
		return sendErr.Result()
	}

	logString := amount.String() + " was successfully minted to " + address.String()
	return sdk.Result{
		Log: logString,
	}
}

// get the proposer public key for this block
func (k Keeper) GetPreviousProposer(ctx sdk.Context) (consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ProposerKey)
	if b == nil {
		panic("Previous proposer not set")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &consAddr)
	return
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposer(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(consAddr)
	store.Set(types.ProposerKey, b)
}

// returns the current BaseProposerReward rate from the global param store
// nolint: errcheck
func (k Keeper) getProposerRewardPercentage(ctx sdk.Context) sdk.Int {
	return sdk.NewInt(int64(k.ProposerRewardPercentage(ctx))).Quo(sdk.NewInt(100))
}
