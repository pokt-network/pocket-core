package keeper

import (
	"bytes"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"time"

	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// Apply and return accumulated updates to the staked validator set
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
	// Retrieve the prevState validator set addresses mapped to their respective staking power
	prevStatePowerMap := k.getPrevStatePowerMap(ctx)
	// Iterate over staked validators, highest power to lowest.
	iterator := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// get the validator address
		valAddr := sdk.Address(iterator.Value())
		// return the validator from the current store
		validator := k.mustGetValidator(ctx, valAddr)
		// sanity check for no jailed validators
		if validator.Jailed {
			panic("should never retrieve a jailed validator from the staked validators")
		}
		if validator.ConsensusPower() == 0 {
			panic("should never have a zero consensus power validator in the staked set")
		}
		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		// check the previous state: if found calculate current power...
		prevStatePowerBytes, found := prevStatePowerMap[valAddrBytes]
		curStatePower := validator.ConsensusPower()
		curStatePowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(curStatePower)
		// if not found or the power has changed -> add this validator to the updated list
		if !found || !bytes.Equal(prevStatePowerBytes, curStatePowerBytes) {
			ctx.Logger().Info(fmt.Sprintf("Updating Validator-Set to Tendermint: %s power changed to %d", validator.Address, validator.ConsensusPower()))
			updates = append(updates, validator.ABCIValidatorUpdate())
			// update the previous state as this will soon be the previous state
			k.SetPrevStateValPower(ctx, valAddr, curStatePower)
		}
		// remove the validator from power map, this structure is used to keep track of who is no longer staked
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
		// delete from the stake validator index
		k.DeletePrevStateValPower(ctx, validator.GetAddress())
		// add to one of the updates for tendermint
		ctx.Logger().Info(fmt.Sprintf("Updating Validator-Set to Tendermint: %s is no longer staked", validator.Address))
		updates = append(updates, validator.ABCIValidatorUpdate())
	}
	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetPrevStateValidatorsPower(ctx, totalPower)
	}
	return updates
}

// validate check called before staking
func (k Keeper) ValidateValidatorStaking(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	// check to see if teh public key has already been register for that validator
	val, found := k.GetValidator(ctx, validator.Address)
	if found {
		if !val.IsUnstaked() {
			return types.ErrValidatorStatus(k.codespace)
		}
		if validator.IsJailed() {
			return types.ErrValidatorJailed(k.codespace)
		}
	} else {
		// check the consensus params
		if ctx.ConsensusParams() != nil {
			tmPubKey := tmTypes.TM2PB.PubKey(validator.PublicKey.PubKey())
			if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
				return types.ErrValidatorPubKeyTypeNotSupported(k.Codespace(),
					tmPubKey.Type,
					ctx.ConsensusParams().Validator.PubKeyTypes)
			}
		}
	}
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.coinKeeper.HasCoins(ctx, validator.Address, coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// store ops when a validator stakes
func (k Keeper) StakeValidator(ctx sdk.Context, validator types.Validator, amount sdk.Int) sdk.Error {
	// send the coins from address to staked module account
	err := k.coinsFromUnstakedToStaked(ctx, validator, amount)
	if err != nil {
		return err
	}
	// add coins to the staked field
	validator = validator.AddStakedTokens(amount)
	// set the status to staked
	validator = validator.UpdateStatus(sdk.Staked)
	// save in the validator store
	k.SetValidator(ctx, validator)
	// save in the staked store
	k.SetStakedValidator(ctx, validator)
	// ensure there's a signing info entry for the validator (used in slashing)
	_, found := k.GetValidatorSigningInfo(ctx, validator.GetAddress())
	if !found {
		signingInfo := types.ValidatorSigningInfo{
			Address:     validator.GetAddress(),
			StartHeight: ctx.BlockHeight(),
			JailedUntil: time.Unix(0, 0),
		}
		k.SetValidatorSigningInfo(ctx, validator.GetAddress(), signingInfo)
	}
	ctx.Logger().Info("Successfully staked validator: " + validator.Address.String())
	return nil
}

func (k Keeper) ValidateValidatorBeginUnstaking(ctx sdk.Context, validator types.Validator) sdk.Error {
	// must be staked to begin unstaking
	if !validator.IsStaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if validator.IsJailed() {
		return types.ErrValidatorJailed(k.codespace)
	}
	// sanity check
	if validator.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: validator trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

func (k Keeper) WaitToBeginUnstakingValidator(ctx sdk.Context, validator types.Validator) sdk.Error {
	k.SetWaitingValidator(ctx, validator)
	ctx.Logger().Info("Validator " + validator.Address.String() + " is waiting to begin unstaking until session is over")
	return nil
}

func (k Keeper) ReleaseWaitingValidators(ctx sdk.Context) {
	vals := k.GetWaitingValidators(ctx)
	for _, val := range vals {
		if err := k.ValidateValidatorBeginUnstaking(ctx, val); err == nil {
			// if able to begin unstaking
			k.BeginUnstakingValidator(ctx, val)
			// create the event
			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.EventTypeBeginUnstake,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, val.Address.String()),
				),
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.EventTypeBeginUnstake),
					sdk.NewAttribute(sdk.AttributeKeySender, val.Address.String()),
				),
			})
		} else {
			ctx.Logger().Error("Unable to begin unstaking validator " + val.Address.String() + ": " + err.Error())
		}
		k.DeleteWaitingValidator(ctx, val.Address)
	}
}

// store ops when validator begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingValidator(ctx sdk.Context, validator types.Validator) {
	// get params
	params := k.GetParams(ctx)
	// delete the validator from the staking set, as it is technically staked but not going to participate
	k.deleteValidatorFromStakingSet(ctx, validator)
	// set the status
	validator = validator.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	validator.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	// save the now unstaked validator record and power index
	k.SetValidator(ctx, validator)
	// Adds to unstaking validator queue
	k.SetUnstakingValidator(ctx, validator)
	ctx.Logger().Info("Began unstaking validator " + validator.Address.String())
}

func (k Keeper) ValidateValidatorFinishUnstaking(ctx sdk.Context, validator types.Validator) sdk.Error {
	if !validator.IsUnstaking() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if validator.IsJailed() {
		return types.ErrValidatorJailed(k.codespace)
	}
	// sanity check
	if validator.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: validator trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops to unstake a validator -> called after unstaking time is up
func (k Keeper) FinishUnstakingValidator(ctx sdk.Context, validator types.Validator) {
	// delete the validator from the unstaking queue
	k.deleteUnstakingValidator(ctx, validator)
	// amount unstaked = stakedTokens
	amount := sdk.NewInt(validator.StakedTokens.Int64())
	// send the tokens from staking module account to validator account
	k.coinsFromStakedToUnstaked(ctx, validator)
	// removed the staked tokens field from validator structure
	validator = validator.RemoveStakedTokens(amount)
	// update the status to unstaked
	validator = validator.UpdateStatus(sdk.Unstaked)
	// update the validator in the main store
	k.SetValidator(ctx, validator)
	ctx.Logger().Info("Finished unstaking validator " + validator.Address.String())
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
}

// force unstake (called when slashed below the minimum)
func (k Keeper) ForceValidatorUnstake(ctx sdk.Context, validator types.Validator) sdk.Error {
	// delete the validator from staking set as they are unstaked
	switch validator.Status {
	case sdk.Staked:
		k.deleteValidatorFromStakingSet(ctx, validator)
	case sdk.Unstaking:
		k.deleteUnstakingValidator(ctx, validator)
	default:
		panic(sdk.ErrInternal("trying to force unstake an already unstaked validator"))
	}
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, validator.StakedTokens)
	if err != nil {
		return err
	}
	// remove their tokens from the field
	validator = validator.RemoveStakedTokens(validator.StakedTokens)
	// update their status to unstaked
	validator = validator.UpdateStatus(sdk.Unstaked)
	// set the validator in store
	k.SetValidator(ctx, validator)
	ctx.Logger().Info("Force Unstaked validator " + validator.Address.String())
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

// send a validator to jail
func (k Keeper) JailValidator(ctx sdk.Context, addr sdk.Address) {
	validator := k.mustGetValidator(ctx, addr)
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}
	k.deleteValidatorFromStakingSet(ctx, validator)
	validator.Jailed = true
	k.SetValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s jailed", addr))
}

func (k Keeper) ValidateUnjailMessage(ctx sdk.Context, msg types.MsgUnjail) (addr sdk.Address, err sdk.Error) {
	validator := k.Validator(ctx, msg.ValidatorAddr)
	if validator == nil {
		return nil, types.ErrNoValidatorForAddress(k.Codespace())
	}
	// cannot be unjailed if no self-delegation exists
	selfDel := validator.GetTokens()
	if selfDel == sdk.ZeroInt() {
		return nil, types.ErrMissingSelfDelegation(k.Codespace())
	}
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) { // TODO look into this state change (stuck in jail)
		return nil, types.ErrSelfDelegationTooLowToUnjail(k.Codespace())
	}
	// cannot be unjailed if not jailed
	if !validator.IsJailed() {
		return nil, types.ErrValidatorNotJailed(k.Codespace())
	}
	addr = validator.GetAddress()
	info, found := k.GetValidatorSigningInfo(ctx, addr)
	if !found {
		return nil, types.ErrNoValidatorForAddress(k.Codespace())
	}
	if info.JailedUntil.After(time.Now()) {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	// cannot be unjailed if tombstoned
	if info.Tombstoned {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	return
}

// remove a validator from jail
func (k Keeper) UnjailValidator(ctx sdk.Context, addr sdk.Address) {
	validator := k.mustGetValidator(ctx, addr)
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
	}
	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetStakedValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s unjailed", addr))
}
