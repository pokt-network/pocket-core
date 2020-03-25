package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/libs/common"
)

// validate check called before staking
func (k Keeper) ValidateApplicationStaking(ctx sdk.Ctx, application types.Application, amount sdk.Int) sdk.Error {
	// convert the amount to sdk.Coin
	coin := sdk.NewCoins(sdk.NewCoin(k.StakeDenom(ctx), amount))
	// attempt to get the application from the world state
	app, found := k.GetApplication(ctx, application.Address)
	// if the application exists
	if found {
		// ensure unstaked
		if !app.IsUnstaked() {
			return types.ErrApplicationStatus(k.codespace)
		}
		if app.IsJailed() {
			return types.ErrApplicationJailed(k.codespace)
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
			if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
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
	if !k.coinKeeper.HasCoins(ctx, application.Address, coin) {
		return types.ErrNotEnoughCoins(k.codespace)
	}
	return nil
}

// store ops when a application stakes
func (k Keeper) StakeApplication(ctx sdk.Ctx, application types.Application, amount sdk.Int) sdk.Error {
	// send the coins from address to staked module account
	err := k.coinsFromUnstakedToStaked(ctx, application, amount)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	// add coins to the staked field
	application = application.AddStakedTokens(amount)
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

func (k Keeper) ValidateApplicationBeginUnstaking(ctx sdk.Ctx, application types.Application) sdk.Error {
	// must be staked to begin unstaking
	if !application.IsStaked() {
		return types.ErrApplicationStatus(k.codespace)
	}
	if application.IsJailed() {
		return types.ErrApplicationJailed(k.codespace)
	}
	// sanity check
	if application.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: application trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops when application begins to unstake -> starts the unstaking timer
func (k Keeper) BeginUnstakingApplication(ctx sdk.Ctx, application types.Application) {
	// get params
	params := k.GetParams(ctx)
	// delete the application from the staking set, as it is technically staked but not going to participate
	k.deleteApplicationFromStakingSet(ctx, application)
	// set the status
	application = application.UpdateStatus(sdk.Unstaking)
	// set the unstaking completion time and completion height appropriately
	if application.UnstakingCompletionTime.Second() == 0 {
		application.UnstakingCompletionTime = ctx.BlockHeader().Time.Add(params.UnstakingTime)
	}
	// save the now unstaked application record and power index
	k.SetApplication(ctx, application)
	// Adds to unstaking application queue
	k.SetUnstakingApplication(ctx, application)
	ctx.Logger().Info("Began unstaking App " + application.Address.String())
}

func (k Keeper) ValidateApplicationFinishUnstaking(ctx sdk.Ctx, application types.Application) sdk.Error {
	if !application.IsUnstaking() {
		return types.ErrApplicationStatus(k.codespace)
	}
	if application.IsJailed() {
		return types.ErrApplicationJailed(k.codespace)
	}
	// sanity check
	if application.StakedTokens.LT(sdk.NewInt(k.MinimumStake(ctx))) {
		panic("should not happen: application trying to begin unstaking has less than the minimum stake")
	}
	return nil
}

// store ops to unstake a application -> called after unstaking time is up
func (k Keeper) FinishUnstakingApplication(ctx sdk.Ctx, application types.Application) {
	// delete the application from the unstaking queue
	k.deleteUnstakingApplication(ctx, application)
	// amount unstaked = stakedTokens
	amount := sdk.NewInt(application.StakedTokens.Int64())
	// send the tokens from staking module account to application account
	k.coinsFromStakedToUnstaked(ctx, application)
	// removed the staked tokens field from application structure
	application = application.RemoveStakedTokens(amount)
	// update the status to unstaked
	application = application.UpdateStatus(sdk.Unstaked)
	// reset app relays
	application.MaxRelays = sdk.ZeroInt()
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

// force unstake (called when slashed below the minimum)
func (k Keeper) ForceApplicationUnstake(ctx sdk.Ctx, application types.Application) sdk.Error {
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

// send a application to jail
func (k Keeper) JailApplication(ctx sdk.Ctx, addr sdk.Address) {
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

func (k Keeper) ValidateUnjailMessage(ctx sdk.Ctx, msg types.MsgAppUnjail) (addr sdk.Address, err sdk.Error) {
	application := k.Application(ctx, msg.AppAddr)
	if application == nil {
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

// remove a application from jail
func (k Keeper) UnjailApplication(ctx sdk.Ctx, addr sdk.Address) {
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
