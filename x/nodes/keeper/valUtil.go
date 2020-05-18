package keeper

import (
	"bytes"
	"container/list"
	"fmt"
	"sort"

	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedValidator struct {
	val     types.Validator
	address sdk.Address
}

func newCachedValidator(val types.Validator, address sdk.Address) cachedValidator {
	return cachedValidator{
		val:     val,
		address: address,
	}
}

// validatorCaching - Retrieve a cached validator
func (k Keeper) validatorCaching(value []byte, addr sdk.Address) types.Validator {
	validator := types.MustUnmarshalValidator(k.cdc, value)
	cachedVal := newCachedValidator(validator, addr)
	k.validatorCache[validator.Address.String()] = cachedVal
	k.validatorCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the prevState element from it
	if k.validatorCacheList.Len() > aminoCacheSize {
		valToRemove := k.validatorCacheList.Remove(k.validatorCacheList.Front()).(cachedValidator)
		delete(k.validatorCache, valToRemove.address.String())
	}
	return validator
}

func (k Keeper) deleteValidatorFromCache(addr sdk.Address) {
	delete(k.validatorCache, addr.String())
}

func (k Keeper) getValidatorFromCache(addr sdk.Address) (validator types.Validator, found bool) {
	if val, ok := k.validatorCache[addr.String()]; ok {
		valToReturn := val.val
		// Doesn't mutate the cache's value
		valToReturn.Address = addr
		return valToReturn, true
	} else {
		return types.Validator{}, false
	}
}

func (k Keeper) searchCacheList(validator types.Validator) (e *list.Element, found bool) {
	for e := k.validatorCacheList.Back(); e != nil; e = e.Prev() {
		v := e.Value.(cachedValidator)
		if v.address.String() == validator.Address.String() {
			return e, true
		}
	}
	return nil, false
}

func (k Keeper) setOrUpdateInValidatorCache(validator types.Validator) {

	e, found := k.searchCacheList(validator)
	if found {
		valToRemove := k.validatorCacheList.Remove(e).(cachedValidator)
		delete(k.validatorCache, valToRemove.address.String())
	}

	cachedVal := newCachedValidator(validator, validator.Address)
	k.validatorCache[validator.Address.String()] = cachedVal
	k.validatorCacheList.PushBack(cachedVal)
}

// Validator - wrapper for GetValidator call
func (k Keeper) Validator(ctx sdk.Ctx, address sdk.Address) exported.ValidatorI {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil
	}
	return val
}

// AllValidators - Retrieve a list of all validators
func (k Keeper) AllValidators(ctx sdk.Ctx) (validators []exported.ValidatorI) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// GetStakedValidators - Retreive StakedValidators
func (k Keeper) GetStakedValidators(ctx sdk.Ctx) (validators []exported.ValidatorI) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, found := k.GetValidator(ctx, iterator.Value())
		if !found {
			ctx.Logger().Error(fmt.Errorf("cannot find validator from staking set: %v\n", iterator.Value()).Error())
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

// map of validator addresses to serialized power
type valPowerMap map[[sdk.AddrLen]byte][]byte

// getPrevStatePowerMap - Retrieve the prevState validator set
func (k Keeper) getPrevStatePowerMap(ctx sdk.Ctx) valPowerMap {
	prevState := make(valPowerMap)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
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
	return prevState
}

// sortNoLongerStakedValidators - Given a map of remaining validators to previous staked power
// returns the list of validators to be unbstaked, sorted by operator address
func sortNoLongerStakedValidators(prevState valPowerMap) [][]byte {
	// sort the map keys for determinism
	noLongerStaked := make([][]byte, len(prevState))
	index := 0
	for valAddrBytes := range prevState {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr, valAddrBytes[:])
		noLongerStaked[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerStaked, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerStaked[i], noLongerStaked[j]) == -1
	})
	return noLongerStaked
}
