package keeper

import (
	"fmt"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/pokt-network/pocket-core/types"
)

// BurnForChallenge - Tries to remove coins from account & supply for a challenged validator
func (k Keeper) BurnForChallenge(ctx sdk.Ctx, challenges sdk.BigInt, address sdk.Address) {
	coins := k.RelaysToTokensMultiplier(ctx).Mul(challenges)
	k.simpleSlash(ctx, address, coins)
}

// simpleSlash - Slash validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) simpleSlash(ctx sdk.Ctx, addr sdk.Address, amount sdk.BigInt) {
	// error check slash
	validator := k.validateSimpleSlash(ctx, addr, amount)
	if validator.Address.Empty() {
		return // invalid simple slash
	}
	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(amount, validator.StakedTokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.
	validator, err := k.removeValidatorTokens(ctx, validator, tokensToBurn)
	if err != nil {
		k.Logger(ctx).Error("could not remove staked tokens in simpleSlash: " + err.Error() + "\nfor validator " + addr.String())
		return
	}
	err = k.burnStakedTokens(ctx, tokensToBurn)
	if err != nil {
		k.Logger(ctx).Error("could not burn staked tokens in simpleSlash: " + err.Error() + "\nfor validator " + addr.String())
		return
	}
	// if falls below minimum force burn all of the stake
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		err := k.ForceValidatorUnstake(ctx, validator)
		if err != nil {
			k.Logger(ctx).Error("could not burn forceUnstake in simpleSlash: " + err.Error() + "\nfor validator " + addr.String())
			return
		}
	}
	// Log that a slash occurred
	ctx.Logger().Info(fmt.Sprintf("validator %s simple slashed; burned %s tokens",
		validator.GetAddress(), amount.String()))
}

// validateSimpleSlash - Check if simpleSlash is possible
func (k Keeper) validateSimpleSlash(ctx sdk.Ctx, addr sdk.Address, amount sdk.BigInt) types.Validator {
	logger := k.Logger(ctx)
	if amount.LTE(sdk.ZeroInt()) {
		k.Logger(ctx).Error(fmt.Errorf("attempted to simple slash with a negative slash factor: %v", amount).Error())
		return types.Validator{}
	}
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		logger.Error(fmt.Sprintf( // could've been overslashed and removed
			"WARNING: Ignored attempt to simple slash a nonexistent validator with address %s, we recommend you investigate immediately",
			addr))
		return types.Validator{}
	}
	// should not be slashing an unstaked validator
	if validator.IsUnstaked() {
		logger.Debug(fmt.Errorf("should not be simple slashing unstaked validator: %s", validator.GetAddress()).Error())
		return types.Validator{}
	}
	return validator
}

// slash - Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) slash(ctx sdk.Ctx, addr sdk.Address, infractionHeight, power int64, slashFactor sdk.BigDec) {
	// error check slash
	validator := k.validateSlash(ctx, addr, infractionHeight, power, slashFactor)
	if validator.Address == nil {
		return // invalid slash
	}
	logger := k.Logger(ctx)
	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := sdk.TokensFromConsensusPower(power)
	slashAmount := amount.ToDec().Mul(slashFactor).TruncateInt()
	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(slashAmount, validator.StakedTokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.
	// Deduct from validator's staked tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	validator, err := k.removeValidatorTokens(ctx, validator, tokensToBurn)
	if err != nil {
		k.Logger(ctx).Error("could not remove staked tokens in slash: " + err.Error() + "\nfor validator " + addr.String())
		return
	}
	err = k.burnStakedTokens(ctx, tokensToBurn)
	if err != nil {
		k.Logger(ctx).Error("could not burn staked tokens in slash: " + err.Error() + "\nfor validator " + addr.String())
		return
	}
	// if falls below minimum force burn all of the stake
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		err := k.ForceValidatorUnstake(ctx, validator)
		if err != nil {
			k.Logger(ctx).Error("could not forceUnstake in slash: " + err.Error() + "\nfor validator " + addr.String())
			return
		}
	}
	// Log that a slash occurred
	logger.Debug(fmt.Sprintf("validator %s slashed by slash factor of %s; burned %v tokens",
		validator.GetAddress(), slashFactor.String(), tokensToBurn))
}

// validateSlash - Check if slash  is possible
func (k Keeper) validateSlash(ctx sdk.Ctx, addr sdk.Address, infractionHeight int64, power int64, slashFactor sdk.BigDec) types.Validator {
	logger := k.Logger(ctx)
	if slashFactor.LTE(sdk.ZeroDec()) {
		k.Logger(ctx).Error(fmt.Errorf("attempted to simple slash with a negative slash factor: %v", slashFactor).Error())
		return types.Validator{}
	}
	if infractionHeight > ctx.BlockHeight() {
		logger.Error(fmt.Sprintf("impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))
		return types.Validator{}
	}
	validator, found := k.GetValidator(ctx, addr)
	if !found {
		logger.Error(fmt.Sprintf( // could've been overslashed and removed
			"WARNING: Ignored attempt to slash a nonexistent validator with address %s, we recommend you investigate immediately",
			addr))
		return types.Validator{}
	}
	// should not be slashing an unstaked validator
	if validator.IsUnstaked() {
		logger.Debug(fmt.Errorf("should not be slashing unstaked validator: %s", validator.GetAddress()).Error())
		return types.Validator{}
	}
	return validator
}

// handleDoubleSign - Handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infractionn
func (k Keeper) handleDoubleSign(ctx sdk.Ctx, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	address, _, _, err := k.validateDoubleSign(ctx, addr, infractionHeight, timestamp)
	if err != nil {
		ctx.Logger().Error(err.Error() + fmt.Sprintf(" at height: %d", ctx.BlockHeight()))
		return
	}
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay
	// get the percentage slash penalty fraction
	fraction := k.SlashFractionDoubleSign(ctx)
	// slash validator
	// `power` is the int64 power of the validator as provided to/by Tendermint. This value is validator.StakedTokens as
	// sent to Tendermint via ABCI, and now received as evidence. The fraction is passed in to separately to slash
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
			sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
		),
	)
	k.slash(ctx, address, distributionHeight, power, fraction)
	// todo fix once tendermint is patched
}

// validateDoubleSign - Check if double signature occurred
func (k Keeper) validateDoubleSign(ctx sdk.Ctx, addr crypto.Address, infractionHeight int64, timestamp time.Time) (address sdk.Address, signInfo types.ValidatorSigningInfo, validator exported.ValidatorI, err sdk.Error) {
	val, found := k.GetValidator(ctx, sdk.Address(addr))
	if !found || val.IsUnstaked() {
		// Ignore evidence that cannot be handled.
		err = types.ErrCantHandleEvidence(k.Codespace())
		return
	}
	pubkey := val.PublicKey
	// calculate the age of the evidence
	t := ctx.BlockHeader().Time
	age := t.Sub(timestamp)
	// Reject evidence if the double-sign is too old
	if age > k.MaxEvidenceAge(ctx) {
		// Ignore evidence that cannot be handled.
		err = sdk.ErrInternal(fmt.Errorf("ignored double sign from %s at height %d, age of %d past max age of %d", sdk.Address(addr), infractionHeight, age, k.MaxEvidenceAge(ctx)).Error())
		return
	}
	// fetch the validator signing info
	signInfo, found = k.GetValidatorSigningInfo(ctx, sdk.Address(addr))
	if !found {
		err = sdk.ErrInternal(fmt.Sprintf("WARNING: Ignored attempt to slash a nonexistent validator with address %s, we recommend you investigate immediately", addr))
		return
	}
	// double sign confirmed
	k.Logger(ctx).Info(fmt.Sprintf("confirmed double sign from %s at height %d, age of %d", sdk.Address(pubkey.Address()), infractionHeight, age))
	return sdk.Address(addr), signInfo, val, nil
}

// handleValidatorSignature - Handle a validator signature, must be called once per validator per block
func (k Keeper) handleValidatorSignature(ctx sdk.Ctx, addr sdk.Address, power int64, signed bool, signedBlocksWindow, minSignedPerWindow int64, downtimeJailDuration time.Duration, slashFractionDowtime sdk.BigDec) {
	_, found := k.GetValidator(ctx, addr)
	if !found {
		ctx.Logger().Info(fmt.Sprintf("in handleValidatorSignature: validator with addr %s not found, "+
			"this is usually due to the 2 block delay between Tendermint and the baseapp", addr))
		return
	}
	// fetch signing info
	signInfo, isFound := k.GetValidatorSigningInfo(ctx, addr)
	if !isFound {
		ctx.Logger().Error(fmt.Sprintf("error in handleValidatorSignature: signing info for validator with addr %s not found, at height %d", addr, ctx.BlockHeight()))
		// patch for june 30 fork
		if ctx.BlockHeight() >= 30040 {
			// reset signing info
			k.ResetValidatorSigningInfo(ctx, addr)
		}
		return
	}
	// reset the validator signing info every blocks window
	if ctx.BlockHeight()%signedBlocksWindow == 0 {
		signInfo.ResetSigningInfo()
		// clear the validator missed at
		k.clearValidatorMissed(ctx, addr)
	}
	// Update signed block bit array
	previous := k.valMissedAt(ctx, addr, signInfo.Index)
	switch {
	case !previous && !signed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedAt(ctx, addr, signInfo.Index, true)
		signInfo.MissedBlocksCounter++
		//ctx.Logger().Info(fmt.Sprintf("Absent validator %s at height %d, %d missed, threshold %d", addr, ctx.BlockHeight(), signInfo.MissedBlocksCounter, minSignedPerWindow))
	case previous && signed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedAt(ctx, addr, signInfo.Index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}
	// update the sign info index
	signInfo.Index++
	// calculate the max missed blocks
	maxMissed := signedBlocksWindow - minSignedPerWindow
	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if signInfo.MissedBlocksCounter > maxMissed {
		// Downtime confirmed: slash and jail the validator
		// ctx.Logger().Info(fmt.Sprintf("Validator %s missed more than the max signed blocks: %d", addr, signedBlocksWindow-minSignedPerWindow))
		// height where the infraction occured
		slashHeight := ctx.BlockHeight() - sdk.ValidatorUpdateDelay - 1
		// slash them based on their power
		k.slash(ctx, addr, slashHeight, power, slashFractionDowtime)
		// reset the signing info
		signInfo.ResetSigningInfo()
		// clear the validator missed at
		k.clearValidatorMissed(ctx, addr)
		// jail the validator to prevent consensus problems
		k.JailValidator(ctx, addr)
		// set the jail time duration
		signInfo.JailedUntil = ctx.BlockHeader().Time.Add(downtimeJailDuration)
	}
	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, addr, signInfo)
}
