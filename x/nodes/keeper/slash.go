package keeper

import (
	"fmt"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/pokt-network/posmint/types"
)

// BurnForChallenge - Tries to remove coins from account & supply for a challenged validator
func (k Keeper) BurnForChallenge(ctx sdk.Ctx, challenges sdk.Int, address sdk.Address) {
	coins := k.RelaysToTokensMultiplier(ctx).Mul(challenges)
	k.simpleSlash(ctx, address, coins)
}

// simpleSlash - Slash validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) simpleSlash(ctx sdk.Ctx, addr sdk.Address, amount sdk.Int) {
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
func (k Keeper) validateSimpleSlash(ctx sdk.Ctx, addr sdk.Address, amount sdk.Int) types.Validator {
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
		logger.Error(fmt.Errorf("should not be simple slashing unstaked validator: %s", validator.GetAddress()).Error())
		return types.Validator{}
	}
	return validator
}

// slash - Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) slash(ctx sdk.Ctx, addr sdk.Address, infractionHeight, power int64, slashFactor sdk.Dec) {
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
	logger.Info(fmt.Sprintf("validator %s slashed by slash factor of %s; burned %v tokens",
		validator.GetAddress(), slashFactor.String(), tokensToBurn))
}

// validateSlash - Check if slash  is possible
func (k Keeper) validateSlash(ctx sdk.Ctx, addr sdk.Address, infractionHeight int64, power int64, slashFactor sdk.Dec) types.Validator {
	logger := k.Logger(ctx)
	if slashFactor.LT(sdk.ZeroDec()) {
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
		logger.Error(fmt.Errorf("should not be slashing unstaked validator: %s", validator.GetAddress()).Error())
		return types.Validator{}
	}
	return validator
}

// handleDoubleSign - Handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infractionn
func (k Keeper) handleDoubleSign(ctx sdk.Ctx, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	address, _, _, err := k.validateDoubleSign(ctx, addr, infractionHeight, timestamp)
	if err != nil {
		ctx.Logger().Error(err.Error())
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
func (k Keeper) handleValidatorSignature(ctx sdk.Ctx, address crypto.Address, power int64, signed bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()
	addr := sdk.Address(address)
	val, found := k.GetValidator(ctx, addr)
	if !found {
		logger.Error(fmt.Sprintf("error in handleValidatorSignature: validator with address %s not found", addr))
		return
	}
	pubkey := val.PublicKey
	// fetch signing info
	signInfo, isFound := k.GetValidatorSigningInfo(ctx, addr)
	if !isFound {
		logger.Error(fmt.Sprintf("error in handleValidatorSignature: signing info for validator with address %s not found", addr))
		return
	}
	//check if validator is not already jailed
	if !val.IsJailed() {
		// this is a relative index, so it counts blocks the validator *should* have signed
		// will use the 0-value default signing info if not present, except for start height
		index := signInfo.IndexOffset % k.SignedBlocksWindow(ctx)
		signInfo.IndexOffset++
		// Update signed block bit array & counter
		// This counter just tracks the sum of the bit array
		// That way we avoid needing to read/write the whole array each time
		previous := k.valMissedAt(ctx, addr, index)
		missed := !signed
		switch {
		case !previous && missed:
			// Array value has changed from not missed to missed, increment counter
			k.SetValidatorMissedAt(ctx, addr, index, true)
			signInfo.MissedBlocksCounter++
		case previous && !missed:
			// Array value has changed from missed to not missed, decrement counter
			k.SetValidatorMissedAt(ctx, addr, index, false)
			signInfo.MissedBlocksCounter--
		default:
			// Array value at this index has not changed, no need to update counter
		}

		if missed {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeLiveness,
					sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
					sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
					sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
				),
			)
			//show log on first missing block sign after successful sign,
			//if missed again this should not pop up and clutter the logs
			if !previous && missed {
				logger.Info(
					fmt.Sprintf("Absent validator %s (%s) at height %d, %d missed, threshold %d",
						addr, pubkey, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
			}
		}
		minHeight := signInfo.StartHeight + k.SignedBlocksWindow(ctx)
		maxMissed := k.SignedBlocksWindow(ctx) - k.MinSignedPerWindow(ctx)
		// if we are past the minimum height and the validator has missed too many blocks, punish them
		if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {

			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				addr, minHeight, k.MinSignedPerWindow(ctx)))
			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the PrevStateCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSlash,
					sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
					sdk.NewAttribute(types.AttributeKeyJailed, addr.String()),
				),
			)
			k.slash(ctx, addr, distributionHeight, power, k.SlashFractionDowntime(ctx))
			k.JailValidator(ctx, addr)
			signInfo.JailedUntil = ctx.BlockHeader().Time.Add(k.DowntimeJailDuration(ctx))
			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon restaking.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearValidatorMissed(ctx, addr)

		}
		// Set the updated signing info
		k.SetValidatorSigningInfo(ctx, addr, signInfo)
	} else {
		//increase JailedBlockCounter
		signInfo.JailedBlocksCounter++
		//Compare against MaxJailedBlocks
		if signInfo.JailedBlocksCounter >= k.MaxJailedBlocks(ctx) {
			//force unstake orphaned validator
			err := k.ForceValidatorUnstake(ctx, val)
			if err != nil {
				k.Logger(ctx).Error("could not forceUnstake jailed validator: " + err.Error() + "\nfor validator " + addr.String())
			} else {
				signInfo.JailedBlocksCounter = 0
			}
		}
		// Set the updated signing info
		k.SetValidatorSigningInfo(ctx, addr, signInfo)
	}
}
