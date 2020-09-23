package keeper

import (
	"fmt"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// PrevStateValidatorsPower - Load the prevState total validator power.
func (k Keeper) PrevStateValidatorsPower(ctx sdk.Ctx) (power sdk.Int) {
	var p = sdk.IntProto{}
	store := ctx.KVStore(k.storeKey)
	b, _ := store.Get(types.PrevStateTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	if ctx.IsAfterUpgradeHeight() {
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(b, &p)
		return p.Int
	} else {
		_ = k.Cdc.UnmarshalBinaryLengthPrefixed(b, &power)
		return power
	}
}

// SetPrevStateValidatorsPower - Store the prevState total validator power (used in moving the curr to prev)
func (k Keeper) SetPrevStateValidatorsPower(ctx sdk.Ctx, power sdk.Int) {
	var p = sdk.IntProto{Int: power}
	store := ctx.KVStore(k.storeKey)
	if ctx.IsAfterUpgradeHeight() {
		b, _ := k.Cdc.MarshalBinaryLengthPrefixed(&p)
		_ = store.Set(types.PrevStateTotalPowerKey, b)
	} else {
		b, _ := k.Cdc.MarshalBinaryLengthPrefixed(&power)
		_ = store.Set(types.PrevStateTotalPowerKey, b)
	}
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
		if ctx.IsAfterUpgradeHeight() {
			var power types.Power
			_ = k.Cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &power)
			if handler(addr, power.Value) {
				break
			}
		} else {
			var power int64
			_ = k.Cdc.UnmarshalBinaryLengthPrefixed(iter.Value(), &power)
			if handler(addr, power) {
				break
			}
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
	p := types.Power{Value: power}
	store := ctx.KVStore(k.storeKey)
	if ctx.IsAfterUpgradeHeight() {
		bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(&p)
		_ = store.Set(types.KeyForValidatorPrevStateStateByPower(addr), bz)
	} else {
		bz, _ := k.Cdc.MarshalBinaryLengthPrefixed(power)
		_ = store.Set(types.KeyForValidatorPrevStateStateByPower(addr), bz)
	}
}

// DeletePrevStateValPower - Remove the power of a SINGLE staked validator from the previous state
func (k Keeper) DeletePrevStateValPower(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValidatorPrevStateStateByPower(addr))
}
