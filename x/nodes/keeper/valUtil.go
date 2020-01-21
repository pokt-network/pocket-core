package keeper

import (
	"bytes"
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"sort"
)

// Cache the amino decoding of validators, as it can be the case that repeated slashing calls
// cause many calls to GetValidator, which were shown to throttle the state machine in our
// simulation. Note this is quite biased though, as the simulator does more slashes than a
// live chain should, however we require the slashing to be fast as noone pays gas for it.
type cachedValidator struct {
	val        types.Validator
	marshalled string // marshalled amino bytes for the Validator object (not operator address)
}

func newCachedValidator(val types.Validator, marshalled string) cachedValidator {
	return cachedValidator{
		val:        val,
		marshalled: marshalled,
	}
}

func (k Keeper) validatorCaching(value []byte, addr sdk.Address) types.Validator {
	// If these amino encoded bytes are in the cache, return the cached Validator
	strValue := string(value)
	if val, ok := k.validatorCache[strValue]; ok {
		valToReturn := val.val
		// Doesn't mutate the cache's value
		valToReturn.Address = addr
		return valToReturn
	}
	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	validator := types.MustUnmarshalValidator(k.cdc, value)
	cachedVal := newCachedValidator(validator, strValue)
	k.validatorCache[strValue] = newCachedValidator(validator, strValue)
	k.validatorCacheList.PushBack(cachedVal)

	// if the cache is too big, pop off the prevState element from it
	if k.validatorCacheList.Len() > aminoCacheSize {
		valToRemove := k.validatorCacheList.Remove(k.validatorCacheList.Front()).(cachedValidator)
		delete(k.validatorCache, valToRemove.marshalled)
	}
	return validator
}

func (k Keeper) mustGetValidator(ctx sdk.Context, addr sdk.Address) types.Validator {
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		panic(fmt.Sprintf("Validator record not found for address: %X\n", addr))
	}
	return validator
}

// wrapper for GetValidator call
func (k Keeper) Validator(ctx sdk.Context, address sdk.Address) exported.ValidatorI {
	val, found := k.GetValidator(ctx, address)
	if !found {
		return nil
	}
	return val
}

func (k Keeper) AllValidators(ctx sdk.Context) (validators []exported.ValidatorI) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

// map of Validator addresses to serialized power
type valPowerMap map[[sdk.AddrLen]byte][]byte

// get the prevState Validator set
func (k Keeper) getPrevStatePowerMap(ctx sdk.Context) valPowerMap {
	prevState := make(valPowerMap)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PrevStateValidatorsPowerKey)
	defer iterator.Close()
	// iterate over the prevState Validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the Validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		prevState[valAddr] = make([]byte, len(powerBytes))
		copy(prevState[valAddr], powerBytes)
	}
	return prevState
}

// given a map of remaining validators to previous staked power
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
