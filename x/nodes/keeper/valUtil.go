package keeper

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/pocket-core/types"
)

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
			ctx.Logger().Error(fmt.Errorf("cannot find validator from staking set: %v, at height %d\n", iterator.Value(), ctx.BlockHeight()).Error())
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
