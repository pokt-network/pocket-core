package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"time"
)

// SetStakedValidator - Store staked validator
func (k Keeper) SetStakedValidator(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Set(types.KeyForValidatorInStakingSet(validator), validator.Address)
	// save in the network id stores for quick session generations
	//k.SetStakedValidatorByChains(ctx, validator)
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
		_ = store.Set(types.KeyForValidatorByNetworkID(validator.Address, cBz), []byte{}) // use empty byte slice to save space
	}
}

// GetValidatorByChains - Returns the validator staked by network identifier
func (k Keeper) GetValidatorsByChain(ctx sdk.Ctx, networkID string) (validators []sdk.Address, count int) {
	defer sdk.TimeTrack(time.Now())
	l, exist := sdk.VbCCache.Get(sdk.GetCacheKey(int(ctx.BlockHeight()), networkID))

	if !exist {
		cBz, err := hex.DecodeString(networkID)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("could not hex decode chains when GetValidatorByChain: with network ID: %s, at height: %d", networkID, ctx.BlockHeight()).Error())
			return
		}
		iterator, _ := k.validatorByChainsIterator(ctx, cBz)
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			address := types.AddressForValidatorByNetworkIDKey(iterator.Key(), cBz)
			count++
			validators = append(validators, address)
		}
		if sdk.VbCCache.Cap() > 1 {
			_ = sdk.VbCCache.Add(sdk.GetCacheKey(int(ctx.BlockHeight()), networkID), validators)
		}

		return validators, count
	}

	validators = l.([]sdk.Address)
	return validators, len(validators)
}

func (k Keeper) deleteValidatorForChains(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	for _, c := range validator.Chains {
		cBz, err := hex.DecodeString(c)
		if err != nil {
			ctx.Logger().Error(fmt.Errorf("could not hex decode chains for validator: %s with network ID: %s, at height %d", validator.Address, c, ctx.BlockHeight()).Error())
			continue
		}
		_ = store.Delete(types.KeyForValidatorByNetworkID(validator.Address, cBz))
	}
}

// validatorByChainsIterator - returns an iterator for the current staked validators
func (k Keeper) validatorByChainsIterator(ctx sdk.Ctx, networkIDBz []byte) (sdk.Iterator, error) {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyForValidatorsByNetworkID(networkIDBz))
}

// deleteValidatorFromStakingSet - delete validator from staked set
func (k Keeper) deleteValidatorFromStakingSet(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Delete(types.KeyForValidatorInStakingSet(validator))
}

// removeValidatorTokens - Update the staked tokens of an existing validator, update the validators power index key
func (k Keeper) removeValidatorTokens(ctx sdk.Ctx, v types.Validator, tokensToRemove sdk.BigInt) (types.Validator, error) {
	k.deleteValidatorFromStakingSet(ctx, v)
	v, err := v.RemoveStakedTokens(tokensToRemove)
	if err != nil {
		return v, err
	}
	k.SetValidator(ctx, v)
	return v, nil
}

// GetStakedValidators - Retrieve StakedValidators
func (k Keeper) GetStakedValidators(ctx sdk.Ctx) (validators []exported.ValidatorI) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, found := k.GetValidator(ctx, iterator.Value())
		if !found {
			ctx.Logger().Error(fmt.Errorf("cannot find validator from staking set: %v, at height %d\n", iterator.Value(), ctx.BlockHeight()).Error())
			continue
		}
		if validator.IsStaked() {
			validators = append(validators, validator)
		}
	}
	return validators
}

// stakedValsIterator - Retrieve an iterator for the current staked validators
func (k Keeper) stakedValsIterator(ctx sdk.Ctx) (sdk.Iterator, error) {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
}

// IterateAndExecuteOverStakedVals - Goes through the staked validator set and execute handler
func (k Keeper) IterateAndExecuteOverStakedVals(
	ctx sdk.Ctx, fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
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
