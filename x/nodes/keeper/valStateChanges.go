package keeper

import (
	"bytes"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Apply and return accumulated updates to the staked Validator set
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
func (k Keeper) UpdateTendermintValidators(ctx sdk.Context) (updates []abci.ValidatorUpdate) {
	// get the world state
	store := ctx.KVStore(k.storeKey)
	// allow all waiting to begin unstaking to begin unstaking
	if ctx.BlockHeight()%k.SessionBlockFrequency(ctx) == 0 { // one block before new session (mod 1 would be session block)
		k.ReleaseWaitingValidators(ctx)
	}
	maxValidators := k.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()
	// Retrieve the prevState Validator set addresses mapped to their respective staking power
	prevStatePowerMap := k.getPrevStatePowerMap(ctx)
	// Iterate over staked validators, highest power to lowest.
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// get the Validator address
		valAddr := sdk.Address(iterator.Value())
		// return the Validator from the current store
		validator := k.mustGetValidator(ctx, valAddr)
		// sanity check for no jailed validators
		if validator.Jailed {
			panic("should never retrieve a jailed Validator from the staked validators")
		}
		if validator.PotentialConsensusPower() == 0 {
			panic("should never have a zero consensus power Validator in the staked set")
		}
		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		// check the previous state: if found calculate current power...
		prevStatePowerBytes, found := prevStatePowerMap[valAddrBytes]
		curStatePower := validator.ConsensusPower()
		curStatePowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(curStatePower)
		// if not found or the power has changed -> add this Validator to the updated list
		if !found || !bytes.Equal(prevStatePowerBytes, curStatePowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate())
			// update the previous state as this will soon be the previous state
			k.SetPrevStateValPower(ctx, valAddr, curStatePower)
		}
		// remove the Validator from power map, this structure is used to keep track of who is no longer staked
		delete(prevStatePowerMap, valAddrBytes)
		// keep count of the number of validators to ensure we don't go over the maximum number of validators
		count++
		// update the total power
		totalPower = totalPower.Add(sdk.NewInt(curStatePower))
	}
	// sort the no-longer-staked validators
	noLongerStaked := sortNoLongerStakedValidators(prevStatePowerMap)
	// iterate through the sorted no-longer-staked validators
	for _, valAddrBytes := range noLongerStaked {
		validator := k.mustGetValidator(ctx, valAddrBytes)
		// delete from the stake Validator index
		k.DeletePrevStateValPower(ctx, validator.GetAddress())
		// add to one of the updates for tendermint
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}
	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetPrevStateValidatorsPower(ctx, totalPower)
	}
	return updates
}

// register the Validator in the necessary stores in the world state
func (k Keeper) RegisterValidator(ctx sdk.Context, validator types.Validator) {
	k.BeforeValidatorRegistered(ctx, validator.Address)
	k.SetValidator(ctx, validator)                     // store Validator here (master list)
	k.SetStakedValidator(ctx, validator)               // store Validator here too (curr staked)
	k.AddPubKeyRelation(ctx, validator.GetPublicKey()) // store relationshiop between consAddr and consPub key
	k.AfterValidatorRegistered(ctx, validator.Address) // call after hook
}

// validate check called before staking
func (k Keeper) ValidateValidatorStaking(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	if !validator.IsUnstaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.coinKeeper.HasCoins(ctx, sdk.Address(validator.Address), coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// store ops when a Validator stakes
func (k Keeper) StakeValidator(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	// call the before hook
	k.BeforeValidatorStaked(ctx, validator.GetAddress())
	// send the coins from address to staked module account
	k.coinsFromUnstakedToStaked(ctx, validator, amount)
	// add coins to the staked field
	validator.AddStakedTokens(amount)
	// set the status to staked
	validator = validator.UpdateStatus(sdk.Staked)
	// save in the Validator store
	k.SetValidator(ctx, validator)
	// save in the staked store
	k.SetStakedValidator(ctx, validator)
	// ensure there's a signing info entry for the Validator (used in slashing)
	_, found := k.GetValidatorSigningInfo(ctx, validator.GetAddress())
	if !found {
		signingInfo := types.ValidatorSigningInfo{
			Address:     validator.GetAddress(),
			StartHeight: ctx.BlockHeight(),
			JailedUntil: time.Unix(0, 0),
		}
		k.SetValidatorSigningInfo(ctx, validator.GetAddress(), signingInfo)
	}
	// call the after hook
	k.AfterValidatorStaked(ctx, validator.GetAddress())
	return nil
}

func (k Keeper) ValidateValidatorBeginUnstaking(ctx sdk.Context, validator types.Validator) sdk.Error {
	// must be staked to begin unstaking
	if !validator.IsStaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	// sanity check
	if validator.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: Validator trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

func (k Keeper) WaitToBeginUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	k.SetWaitingValidator(ctx, validator) // todo could add hooks
	return nil
}

func (k Keeper) ReleaseWaitingValidators(ctx sdk.Context) {
	vals := k.GetWaitingValidators(ctx)
	for _, val := range vals {
		if err := k.ValidateValidatorBeginUnstaking(ctx, val); err == nil {
			// if able to begin unstaking
			if err := k.BeginUnstakingValidator(ctx, val); err != nil {
				panic("error releasing waiting validators: " + err.Error())
			}
			// create the event
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeBeginUnstake,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, val.Address.String()),
				),
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, val.Address.String()),
				),
			})
		}
		k.DeleteWaitingValidator(ctx, val.Address)
	}
}

// store ops when Validator begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call before unstaking hook
	k.BeforeValidatorBeginUnstaking(ctx, validator.GetAddress())
	// get params
	params := k.GetParams(ctx)
	// delete the Validator from the staking set, as it is technically staked but not going to participate
	k.deleteValidatorFromStakingSet(ctx, validator)
	// set the status
	validator = validator.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	validator.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	// save the now unstaked Validator record and power index
	k.SetValidator(ctx, validator)
	// Adds to unstaking Validator queue
	k.SetUnstakingValidator(ctx, validator)
	// call after hook
	k.AfterValidatorBeginUnstaking(ctx, validator.GetAddress())
	return nil
}

func (k Keeper) ValidateValidatorFinishUnstaking(ctx sdk.Context, validator types.Validator) sdk.Error {
	if !validator.IsUnstaking() {
		return types.ErrValidatorStatus(k.codespace)
	}
	// sanity check
	if validator.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: Validator trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops to unstake a Validator -> called after unstaking time is up
func (k Keeper) FinishUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call the before hook
	k.BeforeValidatorUnstaked(ctx, validator.GetAddress(), validator.Address)
	// delete the Validator from the unstaking queue
	k.deleteUnstakingValidator(ctx, validator)
	// amount unstaked = stakedTokens
	amount := sdk.NewInt(validator.StakedTokens.Int64())
	// removed the staked tokens field from Validator structure
	validator = validator.RemoveStakedTokens(amount)
	// send the tokens from staking module account to Validator account
	k.coinsFromStakedToUnstaked(ctx, validator)
	// update the status to unstaked
	validator = validator.UpdateStatus(sdk.Unstaked)
	// update the Validator in the main store
	k.SetValidator(ctx, validator)
	// call the after hook
	k.AfterValidatorUnstaked(ctx, validator.GetAddress(), validator.Address)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
	})
	return nil
}

// force unstake (called when slashed below the minimum)
func (k Keeper) ForceValidatorUnstake(ctx sdk.Context, validator types.Validator) sdk.Error {
	// call the before unstaked hook
	k.BeforeValidatorUnstaked(ctx, validator.GetAddress(), validator.Address)
	// delete the Validator from staking set as they are unstaked
	k.deleteValidatorFromStakingSet(ctx, validator)
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, validator.StakedTokens)
	if err != nil {
		return err
	}
	// remove their tokens from the field
	validator = validator.RemoveStakedTokens(validator.StakedTokens)
	// update their status to unstaked
	validator = validator.UpdateStatus(sdk.Unstaked)
	// set the Validator in store
	k.SetValidator(ctx, validator)
	// call after hook
	k.AfterValidatorUnstaked(ctx, validator.GetAddress(), validator.Address)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, validator.Address.String()),
		),
	})
	return nil
}

// send a Validator to jail
func (k Keeper) JailValidator(ctx sdk.Context, addr sdk.Address) {
	validator := k.mustGetValidator(ctx, addr)
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed Validator, Validator: %v\n", validator))
	}
	validator.Jailed = true
	k.SetValidator(ctx, validator)
	k.deleteValidatorFromStakingSet(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("Validator %s jailed", addr))
}

// remove a Validator from jail
func (k Keeper) UnjailValidator(ctx sdk.Context, addr sdk.Address) {
	validator := k.mustGetValidator(ctx, addr)
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed Validator, Validator: %v\n", validator))
	}
	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetStakedValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("Validator %s unjailed", addr))
}
