package keeper

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// SetWaitingValidator - Store validator on WaitingToBeginUnstaking store
func (k Keeper) SetWaitingValidator(ctx sdk.Ctx, val types.Validator) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Set(types.KeyForValWaitingToBeginUnstaking(val.Address), val.Address)
}

func (k Keeper) SetWaitingValidators(ctx sdk.Ctx, vals types.Validators) {
	for _, val := range vals {
		k.SetWaitingValidator(ctx, val)
	}
}

// IsWaitingValidator - Check if validator is waiting
func (k Keeper) IsWaitingValidator(ctx sdk.Ctx, valAddr sdk.Address) bool {
	store := ctx.KVStore(k.storeKey)
	value, _ := store.Get(types.KeyForValWaitingToBeginUnstaking(valAddr))
	return !(value == nil)
}

// GetWaitingValidators - Retrieve waiting validators
func (k Keeper) GetWaitingValidators(ctx sdk.Ctx) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.WaitingToBeginUnstakingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		addr := iterator.Value()
		validator, found := k.GetValidator(ctx, addr)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("Could not find waiting validator: %s, at height %d", addr, ctx.BlockHeight()))
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
	_ = store.Delete(types.KeyForValWaitingToBeginUnstaking(valAddr))
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
	iterator, _ := sdk.KVStorePrefixIterator(store, types.UnstakingValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var addrs sdk.Addresses
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(iterator.Value(), &addrs, ctx.BlockHeight())
		for _, addr := range addrs {
			validator, found := k.GetValidator(ctx, addr)
			if !found {
				ctx.Logger().Error(fmt.Errorf("cannot find validator from unstaking set: %v, at height %d\n", addr, ctx.BlockHeight()).Error())
				continue
			}
			validators = append(validators, validator)
		}
	}
	return validators

}

// getUnstakingValidators - Retrieve all of the validators who will be unstaked at exactly this time
func (k Keeper) getUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time) (valAddrs sdk.Addresses) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := store.Get(types.KeyForUnstakingValidators(unstakingTime))
	if bz == nil {
		return
	}
	_ = k.Cdc.UnmarshalBinaryLengthPrefixed(bz, &valAddrs, ctx.BlockHeight())
	return valAddrs

}

// setUnstakingValidators - Store validators in unstaking queue at a certain unstaking time
func (k Keeper) setUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time, keys sdk.Addresses) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(&keys, ctx.BlockHeight())
	_ = store.Set(types.KeyForUnstakingValidators(unstakingTime), bz)

}

// deleteUnstakingValidators - Remove all the validators for a specific unstaking time
func (k Keeper) deleteUnstakingValidators(ctx sdk.Ctx, unstakingTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForUnstakingValidators(unstakingTime))
}

// unstakingValidatorsIterator - Retrieve an iterator for all unstaking validators up to a certain time
func (k Keeper) unstakingValidatorsIterator(ctx sdk.Ctx, endTime time.Time) (sdk.Iterator, error) {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.UnstakingValidatorsKey, sdk.InclusiveEndBytes(types.KeyForUnstakingValidators(endTime)))
}

// getMatureValidators - Retrieve a list of all the mature validators
func (k Keeper) getMatureValidators(ctx sdk.Ctx) (matureValsAddrs []sdk.Address) {
	matureValsAddrs = make([]sdk.Address, 0)
	unstakingValsIterator, _ := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValsIterator.Close()
	for ; unstakingValsIterator.Valid(); unstakingValsIterator.Next() {
		var validators sdk.Addresses
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(unstakingValsIterator.Value(), &validators, ctx.BlockHeight())
		matureValsAddrs = append(matureValsAddrs, validators...)

	}
	return matureValsAddrs
}

// unstakeAllMatureValidators -  Unstake all the unstaking validators that have finished their unstaking period
func (k Keeper) unstakeAllMatureValidators(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	unstakingValidatorsIterator, _ := k.unstakingValidatorsIterator(ctx, ctx.BlockHeader().Time)
	defer unstakingValidatorsIterator.Close()
	for ; unstakingValidatorsIterator.Valid(); unstakingValidatorsIterator.Next() {
		var unstakingVals sdk.Addresses
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(unstakingValidatorsIterator.Value(), &unstakingVals, ctx.BlockHeight())
		for _, valAddr := range unstakingVals {
			val, found := k.GetValidator(ctx, valAddr)
			if !found {
				ctx.Logger().Error("validator in the unstaking queue was not found, possible forced unstake? At height: ", ctx.BlockHeight())
				continue
			}
			err := k.ValidateValidatorFinishUnstaking(ctx, val)
			if err != nil {
				ctx.Logger().Error("Could not finish unstaking mature validator: "+err.Error(), "at height: ", ctx.BlockHeight())
				continue
			}
			k.FinishUnstakingValidator(ctx, val)
			k.DeleteValidator(ctx, valAddr)
		}
		_ = store.Delete(unstakingValidatorsIterator.Key())

	}
}
