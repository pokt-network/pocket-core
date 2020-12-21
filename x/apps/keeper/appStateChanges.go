package keeper

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/strings"
	"time"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

// ValidateApplicationStaking - Check application before staking
func (k Keeper) ValidateApplicationStaking(ctx sdk.Ctx, application types.Application, amount sdk.BigInt) sdk.Error {
	// convert the amount to sdk.Coin
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	if int64(len(application.Chains)) > k.MaxChains(ctx) {
		return types.ErrTooManyChains(types.ModuleName)
	}
	// attempt to get the application from the world state
	app, found := k.GetApplication(ctx, application.Address)
	// if the application exists
	if found {
		// ensure unstaked
		if !app.IsUnstaked() {
			return types.ErrApplicationStatus(k.codespace)
		}
	} else {
		// ensure public key type is supported
		if ctx.ConsensusParams() != nil {
			tmPubKey, err := crypto.CheckConsensusPubKey(application.PublicKey.PubKey())
			if err != nil {
				return types.ErrApplicationPubKeyTypeNotSupported(k.Codespace(),
					err.Error(),
					ctx.ConsensusParams().Validator.PubKeyTypes)
			}
			if !strings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
				return types.ErrApplicationPubKeyTypeNotSupported(k.Codespace(),
					tmPubKey.Type,
					ctx.ConsensusParams().Validator.PubKeyTypes)
			}
		}
	}
	// ensure the amount they are staking is < the minimum stake amount
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.AccountsKeeper.HasCoins(ctx, application.Address, coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// StakeApplication - Store ops when a application stakes
func (k Keeper) StakeApplication(ctx sdk.Ctx, application types.Application, amount sdk.BigInt) sdk.Error {
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
	// save in the staked store
	k.SetStakedApplication(ctx, application)
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
	// Adds to unstaking application queue
	k.SetUnstakingApplication(ctx, application)
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

// ForceApplicationUnstake - Coerce unstake (called when slashed below the minimum)
func (k Keeper) ForceApplicationUnstake(ctx sdk.Ctx, application types.Application) sdk.Error {
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
	k.SetStakedApplication(ctx, application)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("application %s unjailed", addr))
}
