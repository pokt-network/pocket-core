package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/libs/strings"
	"time"

	"github.com/pokt-network/pocket-core/crypto"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
)

// UpdateTendermintValidators - Apply and return accumulated updates to the staked validator set
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
func (k Keeper) UpdateTendermintValidators(ctx sdk.Ctx) (updates []abci.ValidatorUpdate) {
	// get the world state
	store := ctx.KVStore(k.storeKey)
	// allow all waiting to begin unstaking to begin unstaking
	if ctx.BlockHeight()%k.BlocksPerSession(ctx) == 0 { // one block before new session (mod 1 would be session block)
		k.ReleaseWaitingValidators(ctx)
	}
	maxValidators := k.MaxValidators(ctx)
	totalPower := sdk.ZeroInt()
	// Retrieve the prevState validator set addresses mapped to their respective staking power
	prevStatePowerMap := k.getPrevStatePowerMap(ctx)
	// Iterate over staked validators, highest power to lowest.
	iterator, _ := sdk.KVStoreReversePrefixIterator(store, types.StakedValidatorsKey)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// get the validator address
		valAddr := sdk.Address(iterator.Value())
		// return the validator from the current store
		validator, found := k.GetValidator(ctx, valAddr)
		if !found {
			k.Logger(ctx).Error("staking validator not found in UpdateTendermintValidators: " + valAddr.String())
			// Cant delete w/out validator object due to ordering
			continue
		}
		// sanity check for no jailed validators
		if validator.Jailed {
			k.Logger(ctx).Error("should never retrieve a jailed validator from the staked validators:" + validator.Address.String())
			continue
		}
		if validator.ConsensusPower() == 0 {
			k.Logger(ctx).Error("should never retrieve a zero power validator from the staked validators:" + validator.Address.String())
			continue
		}
		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		// check the previous state: if found calculate current power...
		prevStatePowerBytes, found := prevStatePowerMap[valAddrBytes]
		curStatePower := validator.ConsensusPower()
		var curStatePowerBytes []byte
		var err error
		csp := sdk.Int64(curStatePower)
		curStatePowerBytes, err = k.Cdc.MarshalBinaryLengthPrefixed(&csp, ctx.BlockHeight())
		if err != nil {
			panic(err)
		}
		// if not found or the power has changed -> add this validator to the updated list
		if !found || !bytes.Equal(prevStatePowerBytes, curStatePowerBytes) {
			ctx.Logger().Debug(fmt.Sprintf("Updating Validator-Set to Tendermint: %s power changed to %d", validator.Address, validator.ConsensusPower()))
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
		validator, found := k.GetValidator(ctx, valAddrBytes)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("unable to retrieve `nolongerstaked` validator: %s from the world state at height: %d ", hex.EncodeToString(valAddrBytes), ctx.BlockHeight()))
			continue
		}
		// delete from the stake validator index
		k.DeletePrevStateValPower(ctx, validator.GetAddress())
		// add to one of the updates for tendermint
		ctx.Logger().Debug(fmt.Sprintf("Updating Validator-Set to Tendermint: %s is no longer staked, at height %d", validator.Address, ctx.BlockHeight()))
		if k.Cdc.IsAfterValidatorSplitUpgrade(ctx.BlockHeight()) {
			updates = append(updates, validator.ABCIValidatorZeroUpdate())
		} else {
			updates = append(updates, validator.ABCIValidatorUpdate())
		}
		// if validator was force unstaked, delete the validator from the all validators store
		if validator.IsUnstaked() {
			k.DeleteValidator(ctx, validator.Address)
		}
	}
	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetPrevStateValidatorsPower(ctx, totalPower)
	}
	return updates
}

// ValidateValidatorStaking - Check Validator before staking
func (k Keeper) ValidateValidatorStaking(ctx sdk.Ctx, validator types.Validator, amount sdk.BigInt, signerAddress sdk.Address) sdk.Error {

	//check the "new" validator's signature validity
	//will recheck if validator exists
	err, valid := ValidateValidatorMsgSigner(validator, signerAddress, k)
	if !valid {
		return err
	}

	//check that we don't allow nil output if we are after noncustodial upgrade
	//so we won't accept stakes/edits with nil outputAddress
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		if validator.OutputAddress == nil {
			return types.ErrNilOutputAddr(k.codespace)
		}
	}

	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))

	if int64(len(validator.Chains)) > k.MaxChains(ctx) {
		return types.ErrTooManyChains(types.ModuleName)
	}

	// check to see if the public key has already been register for that validator
	val, found := k.GetValidator(ctx, validator.Address)
	if found {
		//check again based on the found "state" validator
		err, valid := ValidateValidatorMsgSigner(val, signerAddress, k)
		if !valid {
			return err
		}
		// edit stake in 6.X upgrade
		if ctx.IsAfterUpgradeHeight() && val.IsStaked() {
			return k.ValidateEditStake(ctx, val, validator, amount, signerAddress)
		}
		if !val.IsUnstaked() { // unstaking or already staked but before the upgrade
			return types.ErrValidatorStatus(k.codespace)
		}
	} else {
		// check the consensus params
		if ctx.ConsensusParams() != nil {
			tmPubKey, err := crypto.CheckConsensusPubKey(validator.PublicKey.PubKey())
			if err != nil {
				return types.ErrValidatorPubKeyTypeNotSupported(k.Codespace(),
					err.Error(),
					ctx.ConsensusParams().Validator.PubKeyTypes)
			}
			if !strings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
				return types.ErrValidatorPubKeyTypeNotSupported(k.Codespace(),
					tmPubKey.Type,
					ctx.ConsensusParams().Validator.PubKeyTypes)
			}
		}
	}
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		if !k.AccountKeeper.HasCoins(ctx, signerAddress, coin) {
			return types.ErrNotEnoughCoins(k.codespace)
		}
	} else {
		if !k.AccountKeeper.HasCoins(ctx, validator.Address, coin) {
			return types.ErrNotEnoughCoins(k.codespace)
		}
	}

	return nil
}

//ValidateValidatorMsgSigner Check Validator Signature
func ValidateValidatorMsgSigner(validator types.Validator, signerAddress sdk.Address, k Keeper) (sdk.Error, bool) {
	//check if outputAddress is defined, if not only the operator/node signature is valid
	if validator.OutputAddress == nil {
		if !signerAddress.Equals(validator.Address) {
			return types.ErrUnauthorizedSigner(k.Codespace()), false
		}
	} else {
		if !signerAddress.Equals(validator.Address) && !signerAddress.Equals(validator.OutputAddress) {
			return types.ErrUnauthorizedSigner(k.Codespace()), false
		}
	}
	return nil, true
}

// ValidateEditStake - Validate the updates to a current staked validator
func (k Keeper) ValidateEditStake(ctx sdk.Ctx, currentValidator, newValidtor types.Validator, amount sdk.BigInt, signer sdk.Address) sdk.Error {
	// ensure not staking less
	diff := amount.Sub(currentValidator.StakedTokens)
	if diff.IsNegative() {
		return types.ErrMinimumEditStake(k.codespace)
	}
	// if stake bump
	if !diff.IsZero() {
		// ensure account has enough coins for bump
		coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), diff))
		if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
			if !k.AccountKeeper.HasCoins(ctx, signer, coin) {
				return types.ErrNotEnoughCoins(k.Codespace())
			}
		} else {
			if !k.AccountKeeper.HasCoins(ctx, currentValidator.Address, coin) {
				return types.ErrNotEnoughCoins(k.Codespace())
			}
		}
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		// ensure output address doesn't change
		if currentValidator.OutputAddress != nil {
			if !newValidtor.OutputAddress.Equals(currentValidator.OutputAddress) {
				fmt.Println(currentValidator.String())
				return types.ErrUnequalOutputAddr(k.Codespace())
			}
		}
		// prevent waiting vals from modifying anything
		if k.IsWaitingValidator(ctx, currentValidator.Address) {
			return types.ErrValidatorWaitingToUnstake(types.ModuleName)
		}
	}
	return nil
}

// StakeValidator - Store ops when a validator stakes
func (k Keeper) StakeValidator(ctx sdk.Ctx, validator types.Validator, amount sdk.BigInt, signer crypto.PublicKey) sdk.Error {
	// edit stake
	if ctx.IsAfterUpgradeHeight() {
		// get Validator to see if edit stake
		val, found := k.GetValidator(ctx, validator.Address)
		if found && val.IsStaked() {
			return k.EditStakeValidator(ctx, val, validator, amount, signer)
		}
	}
	// send the coins from address to staked module account
	err := k.coinsFromUnstakedToStaked(ctx, sdk.Address(signer.Address()), amount)
	if err != nil {
		return err
	}
	// add coins to the staked field
	validator, er := validator.AddStakedTokens(amount)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// set the status to staked
	validator = validator.UpdateStatus(sdk.Staked)
	// save in the validator store
	k.SetValidator(ctx, validator)
	k.SetStakedValidatorByChains(ctx, validator)
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

// EditStakeValidator - Edit an already staked validator with the staking message
func (k Keeper) EditStakeValidator(ctx sdk.Ctx, currentValidator, updatedValidator types.Validator, amount sdk.BigInt, signer crypto.PublicKey) sdk.Error {
	origValForDeletion := currentValidator
	// get the difference in coins
	diff := amount.Sub(currentValidator.StakedTokens)
	// if they bumped the stake amount
	if diff.IsPositive() {
		// send the coins from address to staked module account
		err := k.coinsFromUnstakedToStaked(ctx, sdk.Address(signer.Address()), diff)
		if err != nil {
			return err
		}
		currentValidator.StakedTokens = currentValidator.StakedTokens.Add(diff)
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		if currentValidator.OutputAddress == nil {
			currentValidator.OutputAddress = updatedValidator.OutputAddress
		}
	}
	// update chains
	currentValidator.Chains = updatedValidator.Chains
	// update service url
	currentValidator.ServiceURL = updatedValidator.ServiceURL
	// delete the validator from the staking set
	k.deleteValidatorFromStakingSet(ctx, origValForDeletion)
	// delete the validator from each individual chains set
	k.deleteValidatorForChains(ctx, origValForDeletion)
	// delete in main store
	k.DeleteValidator(ctx, origValForDeletion.Address)
	// save in the validator store
	k.SetValidator(ctx, currentValidator)
	// save the validator by chains
	k.SetStakedValidatorByChains(ctx, currentValidator)
	// patch for june 30 fork
	if ctx.BlockHeight() >= 30040 {
		// reset signing info
		k.ResetValidatorSigningInfo(ctx, currentValidator.Address)
	}
	// clear cache
	k.PocketKeeper.ClearSessionCache()
	// log success
	ctx.Logger().Info("Successfully updated staked validator: " + currentValidator.Address.String())
	return nil
}

// ValidateValidatorBeginUnstaking - Check for validator status
func (k Keeper) ValidateValidatorBeginUnstaking(ctx sdk.Ctx, validator types.Validator) sdk.Error {
	// must be staked to begin unstaking
	if !validator.IsStaked() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) { // allow jailed validators to begin unstaking
		return nil
	}
	if validator.IsJailed() {
		return types.ErrValidatorJailed(k.codespace)
	}
	return nil
}

// WaitToBeginUnstakingValidator - Change validator status to waiting to begin unstaking
func (k Keeper) WaitToBeginUnstakingValidator(ctx sdk.Ctx, validator types.Validator) sdk.Error {
	k.SetWaitingValidator(ctx, validator)
	ctx.Logger().Info("Validator " + validator.Address.String() + " is waiting to begin unstaking until session is over")
	return nil
}

// ReleaseWaitingValidators - Remove UnstakingValidators from store
func (k Keeper) ReleaseWaitingValidators(ctx sdk.Ctx) {
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
			ctx.Logger().Info("Unable to begin unstaking validator " + val.Address.String() + ": " + err.Error())
		}
		k.DeleteWaitingValidator(ctx, val.Address)
	}
}

// BeginUnstakingValidator - Store ops when validator begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingValidator(ctx sdk.Ctx, validator types.Validator) {
	// get params
	params := k.GetParams(ctx)
	// delete the validator from the staking set, as it is technically staked but not going to participate
	k.deleteValidatorFromStakingSet(ctx, validator)
	// delete the validator from each individual chains set
	k.deleteValidatorForChains(ctx, validator)
	// set the status
	validator = validator.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	if validator.UnstakingCompletionTime.IsZero() {
		validator.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	}
	// save the now unstaked validator record and power index
	k.SetValidator(ctx, validator)
	ctx.Logger().Info("Began unstaking validator " + validator.Address.String())
}

// ValidateValidatorFinishUnstaking - Check if validator can finish unstaking
func (k Keeper) ValidateValidatorFinishUnstaking(ctx sdk.Ctx, validator types.Validator) sdk.Error {
	if !validator.IsUnstaking() {
		return types.ErrValidatorStatus(k.codespace)
	}
	if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) { // allow jailed validators to finish unstaking
		return nil
	}
	if validator.IsJailed() {
		return types.ErrValidatorJailed(k.codespace)
	}
	return nil
}

// FinishUnstakingValidator - Store ops to unstake a validator -> called after unstaking time is up
func (k Keeper) FinishUnstakingValidator(ctx sdk.Ctx, validator types.Validator) {
	// delete the validator from the unstaking queue
	k.deleteUnstakingValidator(ctx, validator)
	// amount unstaked = stakedTokens
	amount := validator.StakedTokens
	// send the tokens from staking module account to validator account
	err := k.coinsFromStakedToUnstaked(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		// even if error continue with the unstake
	}
	// removed the staked tokens field from validator structure
	validator, err = validator.RemoveStakedTokens(amount)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
		// even if error continue with the unstake
	}
	// update the status to unstaked
	validator = validator.UpdateStatus(sdk.Unstaked)
	// update the unstaking time
	validator.UnstakingCompletionTime = time.Time{}
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

// LegacyForceValidatorUnstake - Coerce unstake (called when slashed below the minimum)
func (k Keeper) LegacyForceValidatorUnstake(ctx sdk.Ctx, validator types.Validator) sdk.Error {
	k.ClearSessionCache()
	// delete the validator from staking set as they are unstaked
	switch validator.Status {
	case sdk.Staked:
		k.deleteValidatorFromStakingSet(ctx, validator)
		k.deleteValidatorForChains(ctx, validator)
		// don't delete validator to allow for previous power to be properly updated
	case sdk.Unstaking:
		k.deleteUnstakingValidator(ctx, validator)
		k.DeleteValidator(ctx, validator.Address)
	default:
		k.DeleteValidator(ctx, validator.Address)
		return sdk.ErrInternal("should not happen: trying to force unstake an already unstaked validator: " + validator.Address.String())
	}
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, validator.StakedTokens)
	if err != nil {
		return err
	}
	if validator.IsStaked() {
		// remove their tokens from the field
		validator, er := validator.RemoveStakedTokens(validator.StakedTokens)
		if er != nil {
			return sdk.ErrInternal(er.Error())
		}
		// update their status to unstaked
		validator = validator.UpdateStatus(sdk.Unstaked)
		// set the validator in store
		k.SetValidator(ctx, validator)
	}
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

// ForceValidatorUnstake - Coerce unstake (called when slashed below the minimum)
func (k Keeper) ForceValidatorUnstake(ctx sdk.Ctx, validator types.Validator) sdk.Error {
	k.ClearSessionCache()
	// send validator to jail || if already jailed, do nothing
	k.JailValidator(ctx, validator.Address)
	ctx.Logger().Info("Sent Validator to Jail for falling below minimum stake" + validator.Address.String())
	k.SetWaitingValidator(ctx, validator)
	ctx.Logger().Info("Validator is waiting to begin unstaking" + validator.Address.String())
	ctx.Logger().Info("" + validator.Address.String())
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBeginUnstake,
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

// JailValidator - Send a validator to jail
func (k Keeper) JailValidator(ctx sdk.Ctx, addr sdk.Address) {
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		ctx.Logger().Error(fmt.Errorf("cannot find jailed validator: %v at height: %d\n", addr, ctx.BlockHeight()).Error())
		return
	}
	if validator.Jailed {
		ctx.Logger().Debug(fmt.Errorf("cannot jail already jailed validator, validator: %v a height: %d\n", validator, ctx.BlockHeight()).Error())
		return
	}
	if validator.IsUnstaked() {
		ctx.Logger().Info(fmt.Errorf("cannot jail an unstaked validator, likely left in the set to update Tendermint Val Set: %v\n", validator).Error())
		return
	}
	// clear caching for sesssions
	k.ClearSessionCache()
	k.deleteValidatorFromStakingSet(ctx, validator)
	validator.Jailed = true
	k.SetValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Debug(fmt.Sprintf("validator %s jailed", addr))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeJail,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
			sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
		),
	)
}

func (k Keeper) IncrementJailedValidators(ctx sdk.Ctx) {
	store := ctx.KVStore(k.storeKey)
	iterator, _ := sdk.KVStorePrefixIterator(store, types.AllValidatorsKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val, err := k.UnmarshalValidator(ctx, iterator.Value())
		if err != nil {
			ctx.Logger().Error("could not unmarshal validator in IncrementJailedValidators: ", err.Error())
			continue
		}
		if val.IsJailed() {
			addr := val.Address
			signInfo, found := k.GetValidatorSigningInfo(ctx, addr)
			if !found {
				k.Logger(ctx).Error("could not find validator signing info in increment jail validator")
				signInfo = types.ValidatorSigningInfo{
					Address:     addr,
					StartHeight: ctx.BlockHeight(),
				}
			}
			// increase JailedBlockCounter
			signInfo.JailedBlocksCounter++
			// compare against MaxJailedBlocks
			if signInfo.JailedBlocksCounter > k.MaxJailedBlocks(ctx) {
				if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
					err := k.ForceValidatorUnstake(ctx, val)
					if err != nil {
						k.Logger(ctx).Error("could not force unstake in simpleSlash: " + err.Error() + "\nfor validator " + addr.String())
						return
					}
				} else {
					err := k.LegacyForceValidatorUnstake(ctx, val)
					if err != nil {
						k.Logger(ctx).Error("could not force unstake in simpleSlash: " + err.Error() + "\nfor validator " + addr.String())
						return
					}
					k.DeleteValidator(ctx, addr)
				}
			} else {
				k.SetValidatorSigningInfo(ctx, addr, signInfo)
			}
		}
	}
}

// ValidateUnjailMessage - Check unjail message
func (k Keeper) ValidateUnjailMessage(ctx sdk.Ctx, msg types.MsgUnjail) (addr sdk.Address, err sdk.Error) {
	validator, found := k.GetValidator(ctx, msg.ValidatorAddr)
	if !found {
		return nil, types.ErrNoValidatorForAddress(k.Codespace())
	}
	//Check msg Signature
	err, valid := ValidateValidatorMsgSigner(validator, msg.Signer, k)
	if !valid {
		return nil, err
	}

	// cannot be unjailed if no self-delegation exists
	selfDel := validator.GetTokens()
	if selfDel == sdk.ZeroInt() {
		return nil, types.ErrMissingSelfDelegation(k.Codespace())
	}
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
			k.SetWaitingValidator(ctx, validator) // defensive against 'stuck in jail'
		}
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
	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	return
}

// UnjailValidator - Remove a validator from jail
func (k Keeper) UnjailValidator(ctx sdk.Ctx, addr sdk.Address) {
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		ctx.Logger().Error(fmt.Errorf("cannot unjail validator, validator not found: %v at height %d\n", addr, ctx.BlockHeight()).Error())
		return
	}
	if !validator.Jailed {
		k.Logger(ctx).Error(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v at height %d\n", validator, ctx.BlockHeight()))
		return
	}
	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.ResetValidatorSigningInfo(ctx, addr)
	k.Logger(ctx).Info(fmt.Sprintf("validator %s unjailed", addr))
}
