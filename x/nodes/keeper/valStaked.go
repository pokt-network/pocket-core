package keeper

import (
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// SetStakedValidator - Store staked validator
func (k Keeper) SetStakedValidator(ctx sdk.Ctx, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	_ = store.Set(types.KeyForValidatorInStakingSet(validator), validator.Address)
	// save in the network id stores for quick session generations
	//k.SetStakedValidatorByChains(ctx, validator)
	k.AddValAddrToCache(ctx, validator.GetAddress())
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
	return validators, count
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
	k.RemoveValAddrFromCache(ctx, validator.Address)
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

// GetStakedValidators - Retrieve StakedValidators
func (k Keeper) GetStakedValidatorsAddrs(ctx sdk.Ctx) (addrs []sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, found := k.GetValidator(ctx, iterator.Value())
		if !found {
			ctx.Logger().Error(fmt.Errorf("cannot find validator from staking set: %v, at height %d\n", iterator.Value(), ctx.BlockHeight()).Error())
			continue
		}
		if validator.IsStaked() {
			addrs = append(addrs, validator.Address)
		}
	}
	return addrs
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

func (k Keeper) GetMemValAddrs(ctx sdk.Ctx) []sdk.Address {
	v, ok := k.getMemValAddrs(ctx)
	//  if doesn't exist get prev from cache
	if !ok {
		return k.GetStakedValidatorsAddrs(ctx)
	}
	addrs, ok := v.([]sdk.Address)
	//  if corrupt get from prev cache
	if !ok {
		return k.GetStakedValidatorsAddrs(ctx)
	}
	return addrs
}

func (k Keeper) AddValAddrToCache(ctx sdk.Ctx, addr sdk.Address) {
	addrs := k.GetMemValAddrs(ctx)
	addrs = append(addrs, addr)
	addrs = k.sortValAddrsByPower(ctx, addrs)
	k.setMemValAddrs(ctx, addrs)
}
func (k Keeper) RemoveValAddrFromCache(ctx sdk.Ctx, addr sdk.Address) {
	addrs := k.GetMemValAddrs(ctx)
	for i := range addrs {
		if addrs[i].Equals(addr) {
			newAddrs := append(addrs[:i], addrs[i+1:]...)
			k.setMemValAddrs(ctx, newAddrs)
			return
		}
	}
}

func (k Keeper) getMemValAddrs(ctx sdk.Ctx) (interface{}, bool) {
	return k.stakedValAddrs.Get(ctx, "staked_val_addrs")
}
func (k Keeper) setMemValAddrs(ctx sdk.Ctx, addr []sdk.Address) bool {
	return k.valPowerCache.Add(ctx, "staked_val_addrs", addr)
}
func (k Keeper) sortValAddrsByPower(ctx sdk.Ctx, addrs []sdk.Address) []sdk.Address {
	sort.SliceStable(addrs, func(i, j int) bool {
		// -1 means strictly less than
		a, _ := k.GetValidator(ctx, addrs[i])
		b, _ := k.GetValidator(ctx, addrs[j])
		return a.ConsensusPower() > b.ConsensusPower()
	})
	return addrs
}
