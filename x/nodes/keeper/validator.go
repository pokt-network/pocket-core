package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// get a single validator from the main store
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.Address) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValByAllVals(addr))
	if value == nil {
		return validator, false
	}
	validator = k.validatorCaching(value, addr)
	return validator, true
}

// set a validator in the main store
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidator(k.cdc, validator)
	store.Set(types.KeyForValByAllVals(validator.Address), bz)
}

// get the set of all validators with no limits from the main store
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// return a given amount of all the validators
func (k Keeper) GetValidators(ctx sdk.Context, maxRetrieve uint16) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	validators = make([]types.Validator, maxRetrieve)
	iterator := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators[i] = validator
		i++
	}
	return validators[:i] // trim if the array length < maxRetrieve
}

// iterate through the validator set and perform the provided function
func (k Keeper) IterateAndExecuteOverVals(
	ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}
