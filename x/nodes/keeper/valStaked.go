package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// SetStakedValidator - Store staked validator
func (k Keeper) SetStakedValidator(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorInStakingSet(validator), validator.Address)
}

// SetStakedValidatorByChains - Store staked validator using networkId
func (k Keeper) SetStakedValidatorByChains(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	for _, c := range validator.Chains {
		cBz, err := hex.DecodeString(c)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("could not hex decode chains for validator: %s with network ID: %s", validator.Address, c).Error())
			continue
		}
		store.Set(types.KeyForValidatorByNetworkID(validator.Address, cBz), []byte{}) // use empty byte slice to save space
	}
}

// GetValidatorByChains - Returns the validator staked by network identifier
func (k Keeper) GetValidatorsByChain(ctx sdk.Ctx, networkID string) (validators []exported.ValidatorI) {
	cBz, err := hex.DecodeString(networkID)
	if err != nil {
		ctx.Logger().Error(fmt.Errorf("could not hex decode chains when GetValidatorByChain: with network ID: %s", networkID).Error())
		return
	}
	iterator := k.validatorByChainsIterator(ctx, cBz)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		address := types.AddressForValidatorByNetworkIDKey(iterator.Key(), cBz)
		validator, found := k.GetValidator(ctx, address)
		if !found {
			ctx.Logger().Error(fmt.Errorf("valdiator: %s, not found in the world state for GetValidatorsByChain call", address).Error())
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

func (k Keeper) deleteValidatorForChains(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	for _, c := range validator.Chains {
		cBz, err := hex.DecodeString(c)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("could not hex decode chains for validator: %s with network ID: %s", validator.Address, c).Error())
			continue
		}
		store.Delete(types.KeyForValidatorByNetworkID(validator.Address, cBz))
	}
}

// validatorByChainsIterator - returns an iterator for the current staked validators
func (k Keeper) validatorByChainsIterator(ctx sdk.Ctx, networkIDBz []byte) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyForValidatorsByNetworkID(networkIDBz))
}

// deleteValidatorFromStakingSet - delete validator from staked set
func (k Keeper) deleteValidatorFromStakingSet(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorInStakingSet(validator))
}

// removeValidatorTokens - Update the staked tokens of an existing validator, update the validators power index key
func (k Keeper) removeValidatorTokens(ctx sdk.Ctx, v types.Validator, tokensToRemove sdk.Int) (types.Validator, error) {
	k.deleteValidatorFromStakingSet(ctx, v)
	v, err := v.RemoveStakedTokens(tokensToRemove)
	if err != nil {
		return v, err
	}
	k.SetValidator(ctx, v)
	return v, nil
}

// getStakedValidators - Retrieve the current staked validators sorted by power-rank
func (k Keeper) getStakedValidators(ctx sdk.Ctx) types.Validators {
	validators := make([]types.Validator, 0)
	iterator := k.stakedValsIterator(ctx)
	defer iterator.Close()
	i := 0
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		validator, found := k.GetValidator(ctx, address)
		if !found {
			ctx.Logger().Error("validator not found in the world state")
			continue
		}
		if validator.IsStaked() {
			validators = append(validators, validator)
			i++
		}
	}
	return validators
}

// stakedValsIterator - Retrieve an iterator for the current staked validators
func (k Keeper) stakedValsIterator(ctx sdk.Ctx) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
}

// IterateAndExecuteOverStakedVals - Goes through the staked validator set and execute handler
func (k Keeper) IterateAndExecuteOverStakedVals(
	ctx sdk.Ctx, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		address := iterator.Value()
		validator, found := k.GetValidator(ctx, address)
		if !found {
			k.Logger(ctx).Error(fmt.Errorf("%s is not found int the main validator state", validator.Address).Error())
			continue
		}
		if validator.IsStaked() {
			stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
			if stop {
				break
			}
			i++
		}
	}
}
