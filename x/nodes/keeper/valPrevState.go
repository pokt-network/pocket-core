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
	k.setPrevStateValidatorsPowerMem(ctx, addr, bz)
}

// SetPrevStateValidatorsPower - Store the prevState total validator power (used in moving the curr to prev)
func (k Keeper) setPrevStateValidatorsPowerMem(ctx sdk.Ctx, addr sdk.Address, powerBz []byte) {
	var valAddr [sdk.AddrLen]byte
	copy(valAddr[:], addr)
	powerMap := k.getMemPrevStatePowerMap(ctx)
	powerMap[valAddr] = powerBz
	k.setMemPowerMap(ctx, powerMap)
}

// DeletePrevStateValPower - Remove the power of a SINGLE staked validator from the previous state
func (k Keeper) DeletePrevStateValPower(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValidatorPrevStateStateByPower(addr))
	k.deletePrevStateValPowerMem(ctx, addr)

}

// DeletePrevStateValPower - Remove the power of a SINGLE staked validator from the previous state
func (k Keeper) deletePrevStateValPowerMem(ctx sdk.Ctx, addr []byte) {
	var valAddr [sdk.AddrLen]byte
	copy(valAddr[:], addr)
	powerMap := k.getMemPrevStatePowerMap(ctx)
	delete(powerMap, valAddr)
	k.setMemPowerMap(ctx, powerMap)
}

// map of validator addresses to serialized power
type valPowerMap map[[sdk.AddrLen]byte][]byte

// getPrevStatePowerMap - Retrieve the prevState validator set
func (k Keeper) getPrevStatePowerMap(ctx sdk.Ctx) valPowerMap {
	prevState := make(valPowerMap)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	defer iterator.Close()
	// iterate over the prevState validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		prevState[valAddr] = make([]byte, len(powerBytes))
		copy(prevState[valAddr], powerBytes)
	}
	k.setMemPowerMap(ctx, prevState) // on the next height will turn previous, otherwise would be left out of the set
	return prevState
}

// getCurStatePowerMapCache - makes sure to retrieve current power map, will become prev on the next block
// CONTRACT: will become prev power map on next height
func (k Keeper) getMemPrevStatePowerMap(ctx sdk.Ctx) valPowerMap {
	v := k.getMemPowerMap(ctx)
	//  if doesn't exist get prev from cache
	if v == nil {
		return k.getPrevStatePowerMap(ctx)
	}
	powerM, ok := v.(valPowerMap)
	//  if corrupt get from prev cache
	if !ok {
		return k.getPrevStatePowerMap(ctx)
	}
	return powerM
}

func (k Keeper) getMemPowerMap(ctx sdk.Ctx) interface{} {
	return k.valPowerCache.Get(ctx)
}
func (k Keeper) setMemPowerMap(ctx sdk.Ctx, prevState valPowerMap) {
	k.valPowerCache.Set(ctx, prevState)
}
