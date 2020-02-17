package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/libs/common"
	tmTypes "github.com/tendermint/tendermint/types"
)

// validate check called before staking
func (k Keeper) ValidateApplicationStaking(ctx sdk.Context, application types.Application, amount sdk.Int) sdk.Error {
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
		// if the application does not exist
	} else {
		// ensure public key type is supported
		if ctx.ConsensusParams() != nil {
			tmPubKey := tmTypes.TM2PB.PubKey(application.PublicKey.PubKey())
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
func (k Keeper) StakeApplication(ctx sdk.Context, application types.Application, amount sdk.Int) sdk.Error {
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

func (k Keeper) ValidateUnjailMessage(ctx sdk.Context, msg types.MsgAppUnjail) (addr sdk.Address, err sdk.Error) {
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
