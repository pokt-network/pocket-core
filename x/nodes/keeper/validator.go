package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// GetValidator - Retrieve validator with address from the main store
func (k Keeper) GetValidator(ctx sdk.Ctx, addr sdk.Address) (validator types.Validator, found bool) {
	val, found := k.validatorCache.GetWithCtx(ctx, addr.String())
	if found {
		return val.(types.Validator), found
	}
	store := ctx.KVStore(k.storeKey)
	value, _ := store.Get(types.KeyForValByAllVals(addr))
	if value == nil {
		return validator, false
	}
	validator, err := types.UnmarshalValidator(k.Cdc, value)
	if err != nil {
		ctx.Logger().Error("can't get validator: " + err.Error())
		return validator, false
	}
	_ = k.validatorCache.AddWithCtx(ctx, addr.String(), validator)
	return validator, true
}

// SetValidator - Store validator in the main store and state stores (stakingset/ unstakingset)
func (k Keeper) SetValidator(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz, err := types.MarshalValidator(k.Cdc, validator)
	if err != nil {
		ctx.Logger().Error("could not marshal validator: " + err.Error())
	}
	err = store.Set(types.KeyForValByAllVals(validator.Address), bz)
	if err != nil {
		ctx.Logger().Error("could not set validator: " + err.Error())
	}
	if validator.IsUnstaking() {
		// Adds to unstaking validator queue
		k.SetUnstakingValidator(ctx, validator)
	}
	if validator.IsStaked() {
		if !validator.IsJailed() {
			// save in the staked store
			k.SetStakedValidator(ctx, validator)
		}
	}
	_ = k.validatorCache.AddWithCtx(ctx, validator.Address.String(), validator)
}

// SetValidator - Store validator in the main store
func (k Keeper) DeleteValidator(ctx sdk.Ctx, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValByAllVals(addr))
	k.DeleteValidatorSigningInfo(ctx, addr)
	k.validatorCache.RemoveWithCtx(ctx, addr.String())
}

// GetAllValidators - Retrieve set of all validators with no limits from the main store
func (k Keeper) GetAllValidators(ctx sdk.Ctx) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.Cdc, iterator.Value())
		if err != nil {
			ctx.Logger().Error("can't get validator in GetAllValidators: " + err.Error())
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

// GetAllValidators - Retrieve set of all validators with no limits from the main store
func (k Keeper) GetAllValidatorsProto(ctx sdk.Ctx) (validators []*types.ValidatorProto) {
	validators = make([]*types.ValidatorProto, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, _ := types.UnmarshalProtoValidator(k.Cdc, iterator.Value())
		validators = append(validators, &validator)
	}
	return validators
}

// GetAllValidators - - Retrieve the set of all validators with no limits from the main store
func (k Keeper) GetAllValidatorsWithOpts(ctx sdk.Ctx, opts types.QueryValidatorsParams) (validators []types.Validator) {
	validators = make([]types.Validator, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.Cdc, iterator.Value())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator in GetAllValidatorsWithOpts: ", err.Error())
			continue
		}
		if opts.IsValid(validator) {
			validators = append(validators, validator)
		}
	}
	return validators
}

// GetValidators - Retrieve a given amount of all the validators
func (k Keeper) GetValidators(ctx sdk.Ctx, maxRetrieve uint16) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)
	validators = make([]types.Validator, maxRetrieve)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.Cdc, iterator.Value())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator in GetValidators: ", err.Error())
			continue
		}
		validators[i] = validator
		i++
	}
	return validators[:i] // trim if the array length < maxRetrieve
}

func (k Keeper) ClearSessionCache() {
	if k.PocketKeeper != nil {
		k.PocketKeeper.ClearSessionCache()
	}
}

// IterateAndExecuteOverVals - Goes through the validator set and executes handler
func (k Keeper) IterateAndExecuteOverVals(
	ctx sdk.Ctx, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		validator, err := types.UnmarshalValidator(k.Cdc, iterator.Value())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator in IterateAndExecuteOverVals: ", err.Error())
			continue
		}
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}
