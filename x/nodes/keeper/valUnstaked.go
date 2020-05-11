package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// SetWaitingValidator - Store validator on WaitingToBeginUnstaking store
func (k Keeper) SetWaitingValidator(ctx sdk.Ctx, val types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValWaitingToBeginUnstaking(val.Address), val.Address)
}

// IsWaitingValidator - Check if validator is waiting
func (k Keeper) IsWaitingValidator(ctx sdk.Ctx, valAddr sdk.Address) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValWaitingToBeginUnstaking(valAddr))
	if value == nil {
		return false
	}
	return true
}

// GetWaitingValidators - Retrieve waiting validators
func (k Keeper) GetWaitingValidators(ctx sdk.Ctx) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WaitingToBeginUnstakingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		addr := iterator.Value()
		validator, found := k.GetValidator(ctx, addr)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("Could not find waiting validator: %s", addr))
			k.DeleteWaitingValidator(ctx, addr)
			return
		}
		validators = append(validators, validator)
	}
	return validators
}

// DeleteWaitingValidator - Remove waiting validators
func (k Keeper) DeleteWaitingValidator(ctx sdk.Ctx, valAddr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValWaitingToBeginUnstaking(valAddr))
}

// SetUnstakingValidator - Store a validator address to the appropriate position in the unstaking queue
func (k Keeper) SetUnstakingValidator(ctx sdk.Ctx, val types.Validator) {
	validators := k.getUnstakingValidators(ctx, val.UnstakingCompletionTime)
	validators = append(validators, val.Address)
	k.setUnstakingValidators(ctx, val.UnstakingCompletionTime, validators)
}

// deleteUnstakingValidator - DeleteInvoice a validator address from the unstaking queue
func (k Keeper) deleteUnstakingValidator(ctx sdk.Ctx, val types.Validator) {
	validators := k.getUnstakingValidators(ctx, val.UnstakingCompletionTime)
	var newValidators []sdk.Address
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

// getAllUnstakingValidators - Retrieve the set of all unstaking validators with no limits
func (k Keeper) getAllUnstakingValidators(ctx sdk.Ctx) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.UnstakingValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addrs []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &addrs)
		for _, addr := range addrs {
			validator, found := k.GetValidator(ctx, addr)
			if !found {
				ctx.Logger().Error(fmt.Errorf("cannot find validator from unstaking set: %v\n", addr).Error())
				continue
			}
			validators = append(validators, validator)
		}
	}
	return validators
}

// getUnstakingValidators - Retrieve all of the validators who will be unstaked at exactly this time
func (k Keeper) getUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time) (valAddrs []sdk.Address) {
	valAddrs = make([]sdk.Address, 0)
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyForUnstakingValidators(unstakingTime))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &valAddrs)
	return
}

// setUnstakingValidators - Store validators in unstaking queue at a certain unstaking time
func (k Keeper) setUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time, keys []sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(keys)
	store.Set(types.KeyForUnstakingValidators(unstakingTime), bz)
}

// deleteUnstakingValidators - Remove all the validators for a specific unstaking time
func (k Keeper) deleteUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForUnstakingValidators(unstakingTime))
}

// unstakingValidatorsIterator - Retrieve an iterator for all unstaking validators up to a certain time
func (k Keeper) unstakingValidatorsIterator(ctx sdk.Ctx, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingValidatorsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingValidators(endTime)))
}

// getMatureValidators - Retrieve a list of all the mature validators
func (k Keeper) getMatureValidators(ctx sdk.Ctx) (matureValsAddrs []sdk.Address) {
	matureValsAddrs = make([]sdk.Address, 0)
	unstakingValsIterator := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()
	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var validators []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &validators)
		matureValsAddrs = append(matureValsAddrs, validators...)
	}
	return matureValsAddrs
}

// unstakeAllMatureValidators -  Unstake all the unstaking validators that have finished their unstaking period
func (k Keeper) unstakeAllMatureValidators(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	unstakingValidatorsIterator := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValidatorsIterator.Close()
	for ; unstakingValidatorsIterator.Valid(); unstakingValidatorsIterator.Next() {
		var unstakingVals []sdk.Address
		k.cdc.MustUnmarshalBinaryLengthPrefixed(unstakingValidatorsIterator.Value(), &unstakingVals)
		for _, valAddr := range unstakingVals {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				ctx.Logger().Error("validator in the unstaking queue was not found, possible forced unstake?")
			}
			err := k.ValidateValidatorFinishUnstaking(ctx, val)
			if err != nil {
				ctx.Logger().Error("Could not finish unstaking mature validator: " + err.Error())
				continue
			}
			k.FinishUnstakingValidator(ctx, val)
		}
		store.Delete(unstakingValidatorsIterator.Key())
	}
}
