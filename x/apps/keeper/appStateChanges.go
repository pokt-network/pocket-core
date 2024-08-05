package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/strings"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

func ensurePubKeyTypeSupported(
	ctx sdk.Ctx,
	pubKey crypto.PublicKey,
	codespace sdk.CodespaceType,
) sdk.Error {
	params := ctx.ConsensusParams()
	if params == nil {
		return nil
	}

	// pubKey.PubKey() converts pocket's PublicKey to tendermint's PubKey
	tmPubKey, err := crypto.CheckConsensusPubKey(pubKey.PubKey())
	if err != nil {
		return types.ErrApplicationPubKeyTypeNotSupported(
			codespace,
			err.Error(),
			params.Validator.PubKeyTypes,
		)
	}

	if !strings.StringInSlice(tmPubKey.Type, params.Validator.PubKeyTypes) {
		return types.ErrApplicationPubKeyTypeNotSupported(
			codespace,
			tmPubKey.Type,
			params.Validator.PubKeyTypes,
		)
	}

	return nil
}

// ValidateApplicationStaking - Check application before staking
func (k Keeper) ValidateApplicationStaking(ctx sdk.Ctx, application types.Application, amount sdk.BigInt) sdk.Error {
	// convert the amount to sdk.Coin
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	// Validate the number of chains that can be staked for
	if int64(len(application.Chains)) > k.MaxChains(ctx) {
		return types.ErrTooManyChains(types.ModuleName)
	}
	// attempt to get the application from the world state
	app, found := k.GetApplication(ctx, application.Address)
	// if the application exists
	if found {
		// edit stake in 6.X upgrade
		if ctx.IsAfterUpgradeHeight() && app.IsStaked() {
			return k.ValidateEditStake(ctx, app, amount)
		}
		if !app.IsUnstaked() { // unstaking or already staked but before the upgrade
			return types.ErrApplicationStatus(k.codespace)
		}
	} else {
		if err := ensurePubKeyTypeSupported(
			ctx,
			application.PublicKey,
			k.Codespace(),
		); err != nil {
			return err
		}
	}
	// ensure the amount they are staking is < the minimum stake amount
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.AccountKeeper.HasCoins(ctx, application.Address, coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	if ctx.IsAfterUpgradeHeight() {
		if k.getStakedApplicationsCount(ctx) >= k.MaxApplications(ctx) {
			return types.ErrMaxApplications(k.codespace)
		}
	}
	return nil
}

// ValidateEditStake - Validate the updates to a current staked validator
func (k Keeper) ValidateEditStake(ctx sdk.Ctx, currentApp types.Application, amount sdk.BigInt) sdk.Error {
	// ensure not staking less
	diff := amount.Sub(currentApp.StakedTokens)
	if diff.IsNegative() {
		return types.ErrMinimumEditStake(k.codespace)
	}
	// if stake bump
	if !diff.IsZero() {
		// ensure account has enough coins for bump
		coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), diff))
		if !k.AccountKeeper.HasCoins(ctx, currentApp.Address, coin) {
			return types.ErrNotEnoughCoins(k.Codespace())
		}
	}
	return nil
}

func (k Keeper) ValidateApplicationTransfer(
	ctx sdk.Ctx,
	signer crypto.PublicKey,
	msg types.MsgStake,
) (types.Application, sdk.Error) {
	if !ctx.IsAfterUpgradeHeight() ||
		!k.Cdc.IsAfterAppTransferUpgrade(ctx.BlockHeight()) {
		return types.Application{}, types.ErrApplicationStatus(k.codespace)
	}

	// The signer must be a staked app
	curApp, found := k.GetApplication(ctx, sdk.Address(signer.Address()))
	if !found || !curApp.IsStaked() {
		return types.Application{}, types.ErrApplicationStatus(k.codespace)
	}

	if err := ensurePubKeyTypeSupported(ctx, msg.PubKey, k.Codespace()); err != nil {
		return types.Application{}, err
	}

	// The pubKey of msgStake must not be a pre-existing application
	newAppAddr := sdk.Address(msg.PubKey.Address())
	if _, found = k.GetApplication(ctx, newAppAddr); found {
		return types.Application{}, types.ErrApplicationStatus(k.codespace)
	}

	return curApp, nil
}

// StakeApplication - Store ops when a application stakes
func (k Keeper) StakeApplication(ctx sdk.Ctx, application types.Application, amount sdk.BigInt) sdk.Error {
	// edit stake
	if ctx.IsAfterUpgradeHeight() {
		// get Validator to see if edit stake
		curApp, found := k.GetApplication(ctx, application.Address)
		if found && curApp.IsStaked() {
			return k.EditStakeApplication(ctx, curApp, application, amount)
		}
	}
	// send the coins from address to staked module account
	err := k.coinsFromUnstakedToStaked(ctx, application, amount)
	if err != nil {
		return err
	}
	// add coins to the staked field
	application, er := application.AddStakedTokens(amount)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// calculate relays
	application.MaxRelays = k.CalculateAppRelays(ctx, application)
	// set the status to staked
	application = application.UpdateStatus(sdk.Staked)
	// save in the application store
	k.SetApplication(ctx, application)
	return nil
}

func (k Keeper) EditStakeApplication(ctx sdk.Ctx, application, updatedApplication types.Application, amount sdk.BigInt) sdk.Error {
	origAppForDeletion := application
	// get the difference in coins
	diff := amount.Sub(application.StakedTokens)
	// if they bumped the stake amount
	if diff.IsPositive() {
		// send the coins from address to staked module account
		err := k.coinsFromUnstakedToStaked(ctx, application, diff)
		if err != nil {
			return err
		}
		var er error
		// add coins to the staked field
		application, er = application.AddStakedTokens(diff)
		if er != nil {
			return sdk.ErrInternal(er.Error())
		}
		// update apps max relays
		application.MaxRelays = k.CalculateAppRelays(ctx, application)
	}
	// update chains
	application.Chains = updatedApplication.Chains
	// delete the validator from the staking set
	k.deleteApplicationFromStakingSet(ctx, origAppForDeletion)
	// delete in main store
	k.DeleteApplication(ctx, origAppForDeletion.Address)
	// save in the app store
	k.SetApplication(ctx, application)
	// save the app by chains
	k.SetStakedApplication(ctx, application)
	// clear session cache
	k.PocketKeeper.ClearSessionCache()
	// log success
	ctx.Logger().Info("Successfully updated staked application: " + application.Address.String())
	return nil
}

// ValidateApplicationBeginUnstaking - Check for validator status
func (k Keeper) ValidateApplicationBeginUnstaking(ctx sdk.Ctx, application types.Application) sdk.Error {
	// must be staked to begin unstaking
	if !application.IsStaked() {
		return sdk.ErrInternal(types.ErrApplicationStatus(k.codespace).Error())
	}
	if application.IsJailed() {
		return sdk.ErrInternal(types.ErrApplicationJailed(k.codespace).Error())
	}
	return nil
}

// BeginUnstakingApplication - Store ops when application begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingApplication(ctx sdk.Ctx, application types.Application) {
	// get params
	params := k.GetParams(ctx)
	// delete the application from the staking set, as it is technically staked but not going to participate
	k.deleteApplicationFromStakingSet(ctx, application)
	// set the status
	application = application.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	if application.UnstakingCompletionTime.IsZero() {
		application.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	}
	// save the now unstaked application record and power index
	k.SetApplication(ctx, application)
	ctx.Logger().Info("Began unstaking App " + application.Address.String())
}

// ValidateApplicationFinishUnstaking - Check if application can finish unstaking
func (k Keeper) ValidateApplicationFinishUnstaking(ctx sdk.Ctx, application types.Application) sdk.Error {
	if !application.IsUnstaking() {
		return types.ErrApplicationStatus(k.codespace)
	}
	if application.IsJailed() {
		return types.ErrApplicationJailed(k.codespace)
	}
	return nil
}

// FinishUnstakingApplication - Store ops to unstake a application -> called after unstaking time is up
func (k Keeper) FinishUnstakingApplication(ctx sdk.Ctx, application types.Application) {
	// delete the application from the unstaking queue
	k.deleteUnstakingApplication(ctx, application)
	// amount unstaked = stakedTokens
	amount := application.StakedTokens
	// send the tokens from staking module account to application account
	err := k.coinsFromStakedToUnstaked(ctx, application)
	if err != nil {
		k.Logger(ctx).Error("could not move coins from staked to unstaked for applications module" + err.Error() + "for this app address: " + application.Address.String())
		// continue with the unstaking
	}
	// removed the staked tokens field from application structure
	application, er := application.RemoveStakedTokens(amount)
	if er != nil {
		k.Logger(ctx).Error("could not remove tokens from unstaking application: " + er.Error())
		// continue with the unstaking
	}
	// update the status to unstaked
	application = application.UpdateStatus(sdk.Unstaked)
	// reset app relays
	application.MaxRelays = sdk.ZeroInt()
	// update the unstaking time
	application.UnstakingCompletionTime = time.Time{}
	// update the application in the main store
	k.SetApplication(ctx, application)
	ctx.Logger().Info("Finished unstaking application " + application.Address.String())
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, application.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, application.Address.String()),
		),
	})
}

// LegacyForceApplicationUnstake - Coerce unstake (called when slashed below the minimum)
func (k Keeper) LegacyForceApplicationUnstake(ctx sdk.Ctx, application types.Application) sdk.Error {
	// delete the application from staking set as they are unstaked
	k.deleteApplicationFromStakingSet(ctx, application)
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, application.StakedTokens)
	if err != nil {
		return err
	}
	// remove their tokens from the field
	application, er := application.RemoveStakedTokens(application.StakedTokens)
	if er != nil {
		return sdk.ErrInternal(er.Error())
	}
	// update their status to unstaked
	application = application.UpdateStatus(sdk.Unstaked)
	// reset app relays
	application.MaxRelays = sdk.ZeroInt()
	// set the application in store
	k.SetApplication(ctx, application)
	ctx.Logger().Info("Force Unstaked application " + application.Address.String())
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, application.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, application.Address.String()),
		),
	})
	return nil
}

// ForceValidatorUnstake - Coerce unstake (called when slashed below the minimum)
func (k Keeper) ForceApplicationUnstake(ctx sdk.Ctx, application types.Application) sdk.Error {
	if !ctx.IsAfterUpgradeHeight() {
		return k.LegacyForceApplicationUnstake(ctx, application)
	}
	switch application.Status {
	case sdk.Staked:
		k.deleteApplicationFromStakingSet(ctx, application)
	case sdk.Unstaking:
		k.deleteUnstakingApplication(ctx, application)
		k.DeleteApplication(ctx, application.Address)
	default:
		k.DeleteApplication(ctx, application.Address)
		return sdk.ErrInternal("should not happen: trying to force unstake an already unstaked application: " + application.Address.String())
	}
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, application.StakedTokens)
	if err != nil {
		return err
	}
	if application.IsStaked() {
		// remove their tokens from the field
		validator, er := application.RemoveStakedTokens(application.StakedTokens)
		if er != nil {
			return sdk.ErrInternal(er.Error())
		}
		// update their status to unstaked
		validator = validator.UpdateStatus(sdk.Unstaked)
		// set the validator in store
		k.SetApplication(ctx, validator)
	}
	ctx.Logger().Info("Force Unstaked validator " + application.Address.String())
	return nil
}

// Transfer the ownership by adding a new app and delete the current app
func (k Keeper) TransferApplication(
	ctx sdk.Ctx,
	curApp types.Application,
	newAppPubKey crypto.PublicKey,
) {
	// Add a new staked application
	// Since we don't change the staked amount, we can inherit all fields
	// and don't need to move tokens.
	newApp := curApp.UpdateStatus(sdk.Staked)
	newApp.Address = sdk.Address(newAppPubKey.Address())
	newApp.PublicKey = newAppPubKey
	k.SetApplication(ctx, newApp)

	// Delete the current application
	// Before UpgradeCodecHeight, we used to keep the app with Unstaked status,
	// but after UpgradeCodecHeight, we just delete the app.
	// (See unstakeAllMatureApplications calling DeleteApplication)
	k.deleteApplicationFromStakingSet(ctx, curApp)
	k.DeleteApplication(ctx, curApp.Address)
}

// JailApplication - Send a application to jail
func (k Keeper) JailApplication(ctx sdk.Ctx, addr sdk.Address) {
	application, found := k.GetApplication(ctx, addr)
	if !found {
		k.Logger(ctx).Error(fmt.Errorf("application %s is attempted jailed but not found in all applications store", addr).Error())
		return
	}
	if application.Jailed {
		k.Logger(ctx).Error(fmt.Sprintf("cannot jail already jailed application, application: %v\n", application))
		return
	}
	application.Jailed = true
	k.SetApplication(ctx, application)
	k.deleteApplicationFromStakingSet(ctx, application)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("application %s jailed", addr))
}

func (k Keeper) IncrementJailedApplications(ctx sdk.Ctx) {
	// TODO
}

// ValidateUnjailMessage - Check unjail message
func (k Keeper) ValidateUnjailMessage(ctx sdk.Ctx, msg types.MsgUnjail) (addr sdk.Address, err sdk.Error) {
	application, found := k.GetApplication(ctx, msg.AppAddr)
	if !found {
		return nil, types.ErrNoApplicationForAddress(k.Codespace())
	}
	// cannot be unjailed if not staked
	stake := application.GetTokens()
	if stake == sdk.ZeroInt() {
		return nil, types.ErrMissingAppStake(k.Codespace())
	}
	if application.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) { // TODO look into this state change (stuck in jail)
		return nil, types.ErrStakeTooLow(k.Codespace())
	}
	// cannot be unjailed if not jailed
	if !application.IsJailed() {
		return nil, types.ErrApplicationNotJailed(k.Codespace())
	}
	return
}

// UnjailApplication - Remove a application from jail
func (k Keeper) UnjailApplication(ctx sdk.Ctx, addr sdk.Address) {
	application, found := k.GetApplication(ctx, addr)
	if !found {
		k.Logger(ctx).Error(fmt.Errorf("application %s is attempted jailed but not found in all applications store", addr).Error())
		return
	}
	if !application.Jailed {
		k.Logger(ctx).Error(fmt.Sprintf("cannot unjail already unjailed application, application: %v\n", application))
		return
	}
	application.Jailed = false
	k.SetApplication(ctx, application)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("application %s unjailed", addr))
}
