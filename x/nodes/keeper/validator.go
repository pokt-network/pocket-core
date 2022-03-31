package keeper

import (
	"bytes"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"sort"
)

func (k Keeper) MarshalValidator(ctx sdk.Ctx, validator types.Validator) ([]byte, error) {
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		bz, err := k.Cdc.MarshalBinaryLengthPrefixed(&validator, ctx.BlockHeight())
		if err != nil {
			ctx.Logger().Error("could not marshal validator: " + err.Error())
		}
		return bz, err
	}
	v := validator.ToLegacy()
	bz, err := k.Cdc.MarshalBinaryLengthPrefixed(&v, ctx.BlockHeight())
	if err != nil {
		ctx.Logger().Error("could not marshal validator: " + err.Error())
	}
	return bz, err
}

func (k Keeper) UnmarshalValidator(ctx sdk.Ctx, valBytes []byte) (val types.Validator, err error) {
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		err = k.Cdc.UnmarshalBinaryLengthPrefixed(valBytes, &val, ctx.BlockHeight())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator: " + err.Error())
		}
		return val, err
	}
	v := types.LegacyValidator{}
	err = k.Cdc.UnmarshalBinaryLengthPrefixed(valBytes, &v, ctx.BlockHeight())
	if err != nil {
		ctx.Logger().Error("could not unmarshal validator: " + err.Error())
	}
	return v.ToValidator(), err
}

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
	validator, err := k.UnmarshalValidator(ctx, value)
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
	bz, err := k.MarshalValidator(ctx, validator)
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

func (k Keeper) SetValidators(ctx sdk.Ctx, validators types.Validators) {
	for _, val := range validators {
		k.SetValidator(ctx, val)
	}
}

func (k Keeper) GetValidatorOutputAddress(ctx sdk.Ctx, operatorAddress sdk.Address) (sdk.Address, bool) {
	val, found := k.GetValidator(ctx, operatorAddress)
	if val.OutputAddress == nil {
		return val.Address, found
	}
	return val.OutputAddress, found
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
		validator, err := k.UnmarshalValidator(ctx, iterator.Value())
		if err != nil {
			ctx.Logger().Error("can't get validator in GetAllValidators: " + err.Error())
			continue
		}
		validators = append(validators, validator)
	}
	return validators
}

// GetAllValidators - Retrieve set of all validators with no limits from the main store
func (k Keeper) GetAllValidatorsAddrs(ctx sdk.Ctx) (validators []sdk.Address) {
	validators = make([]sdk.Address, 0)
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validators = append(validators, iterator.Key())
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
		validator, err := k.UnmarshalValidator(ctx, iterator.Value())
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
		validator, err := k.UnmarshalValidator(ctx, iterator.Value())
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
		validator, err := k.UnmarshalValidator(ctx, iterator.Value())
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
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator, err := k.UnmarshalValidator(ctx, iterator.Value())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator in AllValidators: ", err.Error())
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
