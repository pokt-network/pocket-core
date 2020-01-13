package keeper

import (
	"bytes"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"time"
)

func (k Keeper) SetWaitingValidator(ctx sdk.Context, val types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidator(k.cdc, val)
	store.Set(types.KeyForValWaitingToBeginUnstaking(val.Address), bz)
}

func (k Keeper) IsWaitingValidator(ctx sdk.Context, valAddr sdk.ValAddress) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValWaitingToBeginUnstaking(valAddr))
	if value == nil {
		return false
	}
	return true
}

func (k Keeper) GetWaitingValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WaitingToBeginUnstakingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

func (k Keeper) DeleteWaitingValidator(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValWaitingToBeginUnstaking(valAddr))
}

// Insert a validator address to the appropriate position in the unstaking queue
func (k Keeper) SetUnstakingValidator(ctx sdk.Context, val types.Validator) {
	validators := k.getUnstakingValidators(ctx, val.UnstakingCompletionTime)
	validators = append(validators, val.Address)
	k.setUnstakingValidators(ctx, val.UnstakingCompletionTime, validators)
}

// DeleteInvoice a validator address from the unstaking queue
func (k Keeper) deleteUnstakingValidator(ctx sdk.Context, val types.Validator) {
	validators := k.getUnstakingValidators(ctx, val.UnstakingCompletionTime)
	var newValidators []sdk.ValAddress
	for _, addr := range validators {
		if !bytes.Equal(addr, val.Address) {
			newValidators = append(newValidators, addr)
		}
	}
	if len(newValidators) == 0 {
		k.deleteUnstakingValidators(ctx, val.UnstakingCompletionTime)
	} else {
		k.setUnstakingValidators(ctx, val.UnstakingCompletionTime, newValidators)
	}
}

// get the set of all unstaking validators with no limits
func (k Keeper) getAllUnstakingValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakingValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// retrieve all unstaked validators with no limits
func (k Keeper) getAllUnstakedValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakedValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// gets all of the validators who will be unstaked at exactly this time
func (k Keeper) getUnstakingValidators(ctx sdk.Context, unstakingTime time.Time) (valAddrs []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForUnstakingValidators(unstakingTime))
	if bz == nil {
		return []sdk.ValAddress{}
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return valAddrs
}

// Sets validators in unstaking queue at a certain unstaking time
func (k Keeper) setUnstakingValidators(ctx sdk.Context, unstakingTime time.Time, keys []sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.KeyForUnstakingValidators(unstakingTime), bz)
}

// Deletes all the validators for a specific unstaking time
func (k Keeper) deleteUnstakingValidators(ctx sdk.Context, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForUnstakingValidators(unstakingTime))
}

// iterator for all unstaking validators up to a certain time
func (k Keeper) unstakingValidatorsIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingValidatorsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingValidators(endTime)))
}

// Returns a list of all the mature validators
func (k Keeper) getMatureValidators(ctx sdk.Context) (matureValsAddrs []sdk.ValAddress) {
	unstakingValsIterator := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()
	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var validators []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &validators)
		matureValsAddrs = append(matureValsAddrs, validators...)
	}
	return matureValsAddrs
}

// Unstakes all the unstaking validators that have finished their unstaking period
func (k Keeper) unstakeAllMatureValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	unstakingValidatorsIterator := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValidatorsIterator.Close()
	for ; unstakingValidatorsIterator.Valid(); unstakingValidatorsIterator.Next() {
		var unstakingVals []sdk.ValAddress
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValidatorsIterator.Value(), &unstakingVals)
		for _, valAddr := range unstakingVals {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				panic("validator in the unstaking queue was not found")
			}
			err := k.ValidateValidatorFinishUnstaking(ctx, val)
			if err != nil {
				panic(err)
			}
			err = k.FinishUnstakingValidator(ctx, val)
			if err != nil {
				panic(err)
			}
		}
		store.Delete(unstakingValidatorsIterator.Key())
	}
}
