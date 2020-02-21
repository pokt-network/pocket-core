package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"time"

	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) BurnValidator(ctx sdk.Context, address sdk.Address, severityPercentage sdk.Dec) {
	curBurn, _ := k.getValidatorBurn(ctx, address)
	newSeverity := curBurn.Add(severityPercentage)
	k.setValidatorBurn(ctx, newSeverity, address)
}

// slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
func (k Keeper) slash(ctx sdk.Context, consAddr sdk.Address, infractionHeight, power int64, slashFactor sdk.Dec) {
	// error check slash
	validator := k.validateSlash(ctx, consAddr, infractionHeight, power, slashFactor)
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
	validator = k.removeValidatorTokens(ctx, validator, tokensToBurn)
	err := k.burnStakedTokens(ctx, tokensToBurn)
	if err != nil {
		panic(err)
	}
	// if falls below minimum force burn all of the stake
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		err := k.ForceValidatorUnstake(ctx, validator)
		if err != nil {
			panic(err)
		}
	}
	// Log that a slash occurred
	logger.Info(fmt.Sprintf("validator %s slashed by slash factor of %s; burned %v tokens",
		validator.GetAddress(), slashFactor.String(), tokensToBurn))
}

func (k Keeper) validateSlash(ctx sdk.Context, addr sdk.Address, infractionHeight int64, power int64, slashFactor sdk.Dec) types.Validator {
	logger := k.Logger(ctx)
	if slashFactor.LT(sdk.ZeroDec()) {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}
	if infractionHeight > ctx.BlockHeight() {
		panic(fmt.Errorf( // Can't slash infractions in the future
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))
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
		panic(fmt.Errorf("should not be slashing unstaked validator: %s", validator.GetAddress()))
	}
	return validator
}

// handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infraction
func (k Keeper) handleDoubleSign(ctx sdk.Context, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	address, _, _, err := k.validateDoubleSign(ctx, addr, infractionHeight, timestamp)
	if err != nil {
		panic(err)
	}
	// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height.
	// Note that this *can* result in a negative "distributionHeight", up to -ValidatorUpdateDelay,
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

func (k Keeper) validateDoubleSign(ctx sdk.Context, addr crypto.Address, infractionHeight int64, timestamp time.Time) (address sdk.Address, signInfo types.ValidatorSigningInfo, validator exported.ValidatorI, err sdk.Error) {
	logger := k.Logger(ctx)
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
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, age of %d past max age of %d",
			sdk.Address(addr), infractionHeight, age, k.MaxEvidenceAge(ctx)))
		return
	}
	// fetch the validator signing info
	signInfo, found = k.GetValidatorSigningInfo(ctx, sdk.Address(addr))
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", sdk.Address(addr)))
	}
	// validator is already tombstoned
	if signInfo.Tombstoned {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already tombstoned", sdk.Address(pubkey.Address()), infractionHeight))
		err = types.ErrValidatorTombstoned(k.Codespace())
		return
	}
	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", sdk.Address(pubkey.Address()), infractionHeight, age))
	return sdk.Address(addr), signInfo, val, nil
}

// handle a validator signature, must be called once per validator per block
func (k Keeper) handleValidatorSignature(ctx sdk.Context, address crypto.Address, power int64, signed bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()
	addr := sdk.Address(address)
	val, found := k.GetValidator(ctx, addr)
	if !found {
		panic(fmt.Sprintf("Validator consensus-address %s not found", addr))
	}
	pubkey := val.PublicKey
	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, addr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", addr))
	}
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
		logger.Info(
			fmt.Sprintf("Absent validator %s (%s) at height %d, %d missed, threshold %d", addr, pubkey, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}
	minHeight := signInfo.StartHeight + k.SignedBlocksWindow(ctx)
	maxMissed := k.SignedBlocksWindow(ctx) - k.MinSignedPerWindow(ctx)
	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator, found := k.GetValidator(ctx, addr)
		if found && !validator.IsJailed() {
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
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed", addr),
			)
		}
	}
	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, addr, signInfo)
}

func (k Keeper) getBurnFromSeverity(ctx sdk.Context, address sdk.Address, severityPercentage sdk.Dec) sdk.Int {
	val := k.mustGetValidator(ctx, address)
	amount := sdk.TokensFromConsensusPower(val.ConsensusPower())
	slashAmount := amount.ToDec().Mul(severityPercentage).TruncateInt()
	return slashAmount
}

// called on begin blocker
func (k Keeper) burnValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BurnValidatorKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		severity := sdk.Dec{}
		address := sdk.Address(types.AddressFromKey(iterator.Key()))
		amino.MustUnmarshalBinaryBare(iterator.Value(), &severity)
		val := k.mustGetValidator(ctx, address)
		k.slash(ctx, address, ctx.BlockHeight(), val.ConsensusPower(), severity)
		// remove from the burn store
		store.Delete(iterator.Key())
	}
}

// store functions used to keep track of a validator burn
func (k Keeper) setValidatorBurn(ctx sdk.Context, amount sdk.Dec, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyForValidatorBurn(address), amino.MustMarshalBinaryBare(amount))
}

func (k Keeper) getValidatorBurn(ctx sdk.Context, address sdk.Address) (coins sdk.Dec, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(types.KeyForValidatorBurn(address))
	if value == nil {
		return coins, false
	}
	found = true
	k.cdc.MustUnmarshalBinaryBare(value, &coins)
	return
}

func (k Keeper) deleteValidatorBurn(ctx sdk.Context, address sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyForValidatorBurn(address))
}
