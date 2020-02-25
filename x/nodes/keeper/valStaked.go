package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// set staked validator
func (k Keeper) SetStakedValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorInStakingSet(validator), validator.Address)
}

// delete validator from staked set
func (k Keeper) deleteValidatorFromStakingSet(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorInStakingSet(validator))
}

// Update the staked tokens of an existing validator, update the validators power index key
func (k Keeper) removeValidatorTokens(ctx sdk.Context, v types.Validator, tokensToRemove sdk.Int) types.Validator {
	k.deleteValidatorFromStakingSet(ctx, v)
	v = v.RemoveStakedTokens(tokensToRemove)
	k.SetValidator(ctx, v)
	k.SetStakedValidator(ctx, v)
	return v
}

// get the current staked validators sorted by power-rank
func (k Keeper) getStakedValidators(ctx sdk.Context) types.Validators {
	validators := make([]types.Validator, 0)
	iterator := k.stakedValsIterator(ctx)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		validator := k.mustGetValidator(ctx, address)
		if validator.IsStaked() {
			validators = append(validators, validator)
			i++
		}
	}
	return validators
}

// returns an iterator for the current staked validators
func (k Keeper) stakedValsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
}

// iterate through the staked validator set and perform the provided function
func (k Keeper) IterateAndExecuteOverStakedVals(
	ctx sdk.Context, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		validator := k.mustGetValidator(ctx, address)
		if validator.IsStaked() {
			stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}
