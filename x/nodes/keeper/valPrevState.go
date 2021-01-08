package keeper

import (
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// PrevStateValidatorsPower - Load the prevState total validator power.
func (k Keeper) PrevStateValidatorsPower(ctx sdk.Ctx) (power sdk.BigInt) {
	store := ctx.KVStore(k.storeKey)
	b, _ := store.Get(types.PrevStateTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	_ = k.Cdc.UnmarshalBinaryLengthPrefixed(b, &power, ctx.BlockHeight())
	return power
}

// SetPrevStateValidatorsPower - Store the prevState total validator power (used in moving the curr to prev)
func (k Keeper) SetPrevStateValidatorsPower(ctx sdk.Ctx, power sdk.BigInt) {
	store := ctx.KVStore(k.storeKey)
	b, _ := k.Cdc.MarshalBinaryLengthPrefixed(&power, ctx.BlockHeight())
	_ = store.Set(types.PrevStateTotalPowerKey, b)
}

// prevStateValidatorIterator - Retrieve an iterator for the consensus validators in the prevState block
func (k Keeper) prevStateValidatorsIterator(ctx sdk.Ctx) (iterator sdk.Iterator) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ = sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	return iterator
}

// IterateAndExecuteOverPrevStateValsByPower - Goes over prevState validator powers and perform a function on each validator.
func (k Keeper) IterateAndExecuteOverPrevStateValsByPower(
	ctx sdk.Ctx, handler func(address sdk.Address, power int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter, _ := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.Address(iter.Key()[len(types.PrevStateValidatorsPowerKey):])
		var power sdk.Int64
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &power, ctx.BlockHeight())
		if handler(addr, int64(power)) {
			break
		}
	}
}

// IterateAndExecuteOverPrevStateVals - Goes through the active validator set and perform the provided function
func (k Keeper) IterateAndExecuteOverPrevStateVals(
	ctx sdk.Ctx, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	iterator := k.prevStateValidatorsIterator(ctx)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := types.AddressFromKey(iterator.Key())
		validator, found := k.GetValidator(ctx, address)
		if !found {
			ctx.Logger().Error(fmt.Errorf("validator record not found for address: %v, at height %d\n", address, ctx.BlockHeight()).Error())
			continue
		}
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

// SetPrevStateValPower - Store the power of a SINGLE staked validator from the previous state
func (k Keeper) SetPrevStateValPower(ctx sdk.Ctx, addr sdk.Address, power int64) {
	store := ctx.KVStore(k.storeKey)
	a := sdk.Int64(power)
	bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(&a, ctx.BlockHeight())
	_ = store.Set(types.KeyForValidatorPrevStateStateByPower(addr), bz)
}

// DeletePrevStateValPower - Remove the power of a SINGLE staked validator from the previous state
func (k Keeper) DeletePrevStateValPower(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValidatorPrevStateStateByPower(addr))
}
