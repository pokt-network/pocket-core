package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
)

// register the application in the necessary stores in the world state
func (k Keeper) RegisterApplication(ctx sdk.Context, application types.Application) {
	k.BeforeApplicationRegistered(ctx, application.Address)
	k.SetApplication(ctx, application)                     // store application here (master list)
	k.AfterApplicationRegistered(ctx, application.Address) // call after hook
}

// validate check called before staking
func (k Keeper) ValidateApplicationStaking(ctx sdk.Context, application types.Application, amount sdk.Int) sdk.Error {
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	if !application.IsUnstaked() {
		return types.ErrApplicationStatus(k.codespace)
	}
	if amount.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return types.ErrMinimumStake(k.codespace)
	}
	if !k.coinKeeper.HasCoins(ctx, sdk.Address(application.Address), coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// store ops when a application stakes
func (k Keeper) StakeApplication(ctx sdk.Context, application types.Application, amount sdk.Int) sdk.Error {
	// call the before hook
	k.BeforeApplicationStaked(ctx, application.GetAddress(), application.Address)
	// send the coins from address to staked module account
	err := k.coinsFromUnstakedToStaked(ctx, application, amount)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	// add coins to the staked field
	application.AddStakedTokens(amount)
	// calculate relays
	application.MaxRelays = k.CalculateAppRelays(ctx, application)
	// set the status to staked
	application = application.UpdateStatus(sdk.Staked)
	// save in the application store
	k.SetApplication(ctx, application)
	// save in the staked store
	k.SetStakedApplication(ctx, application)
	// call the after hook
	k.AfterApplicationStaked(ctx, application.GetAddress(), application.Address)
	return nil
}

func (k Keeper) ValidateApplicationBeginUnstaking(ctx sdk.Context, application types.Application) sdk.Error {
	// must be staked to begin unstaking
	if !application.IsStaked() {
		return types.ErrApplicationStatus(k.codespace)
	}
	// sanity check
	if application.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: application trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops when application begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingApplication(ctx sdk.Context, application types.Application) sdk.Error {
	// call before unstaking hook
	k.BeforeApplicationBeginUnstaking(ctx, application.GetAddress(), application.Address)
	// get params
	params := k.GetParams(ctx)
	// delete the application from the staking set, as it is technically staked but not going to participate
	k.deleteApplicationFromStakingSet(ctx, application)
	// set the status
	application = application.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	application.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	// save the now unstaked application record and power index
	k.SetApplication(ctx, application)
	// Adds to unstaking application queue
	k.SetUnstakingApplication(ctx, application)
	// call after hook
	k.AfterApplicationBeginUnstaking(ctx, application.GetAddress(), application.Address)
	return nil
}

func (k Keeper) ValidateApplicationFinishUnstaking(ctx sdk.Context, application types.Application) sdk.Error {
	if !application.IsUnstaking() {
		return types.ErrApplicationStatus(k.codespace)
	}
	// sanity check
	if application.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: application trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops to unstake a application -> called after unstaking time is up
func (k Keeper) FinishUnstakingApplication(ctx sdk.Context, application types.Application) sdk.Error {
	// call the before hook
	k.BeforeApplicationUnstaked(ctx, application.GetAddress(), application.Address)
	// delete the application from the unstaking queue
	k.deleteUnstakingApplication(ctx, application)
	// amount unstaked = stakedTokens
	amount := sdk.NewInt(application.StakedTokens.Int64())
	// removed the staked tokens field from application structure
	application = application.RemoveStakedTokens(amount)
	// send the tokens from staking module account to application account
	k.coinsFromStakedToUnstaked(ctx, application)
	// update the status to unstaked
	application = application.UpdateStatus(sdk.Unstaked)
	// reset app relays
	application.MaxRelays = sdk.ZeroInt()
	// update the application in the main store
	k.SetApplication(ctx, application)
	// call the after hook
	k.AfterApplicationUnstaked(ctx, application.GetAddress(), application.Address)
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

// force unstake (called when slashed below the minimum)
func (k Keeper) ForceApplicationUnstake(ctx sdk.Context, application types.Application) sdk.Error {
	// call the before unstaked hook
	k.BeforeApplicationUnstaked(ctx, application.GetAddress(), application.Address)
	// delete the application from staking set as they are unstaked
	k.deleteApplicationFromStakingSet(ctx, application)
	// amount unstaked = stakedTokens
	err := k.burnStakedTokens(ctx, application.StakedTokens)
	if err != nil {
		return err
	}
	// remove their tokens from the field
	application = application.RemoveStakedTokens(application.StakedTokens)
	// update their status to unstaked
	application = application.UpdateStatus(sdk.Unstaked)
	// reset app relays
	application.MaxRelays = sdk.ZeroInt()
	// set the application in store
	k.SetApplication(ctx, application)
	// call after hook
	k.AfterApplicationUnstaked(ctx, application.GetAddress(), application.Address)
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

// send a application to jail
func (k Keeper) JailApplication(ctx sdk.Context, addr sdk.Address) {
	application := k.mustGetApplicationByConsAddr(ctx, addr)
	if application.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed application, application: %v\n", application))
	}
	application.Jailed = true
	k.SetApplication(ctx, application)
	k.deleteApplicationFromStakingSet(ctx, application)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("application %s jailed", addr))
}

// remove a application from jail
func (k Keeper) UnjailApplication(ctx sdk.Context, addr sdk.Address) {
	application := k.mustGetApplicationByConsAddr(ctx, addr)
	if !application.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed application, application: %v\n", application))
	}
	application.Jailed = false
	k.SetApplication(ctx, application)
	k.SetStakedApplication(ctx, application)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("application %s unjailed", addr))
}
