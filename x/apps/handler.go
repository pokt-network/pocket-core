package pos

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/keeper"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"reflect"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.Msg, _ crypto.PublicKey) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		// convert to value for switch consistency
		if reflect.ValueOf(msg).Kind() == reflect.Ptr {
			msg = reflect.Indirect(reflect.ValueOf(msg)).Interface().(sdk.Msg)
		}
		switch msg := msg.(type) {
		case types.MsgStake:
			return handleStake(ctx, msg, k)
		case types.MsgBeginUnstake:
			return handleMsgBeginUnstake(ctx, msg, k)
		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStake(ctx sdk.Ctx, msg types.MsgStake, k keeper.Keeper) sdk.Result {
	pk := msg.PubKey
	addr := pk.Address()
	ctx.Logger().Info("Begin Staking App Message received from " + sdk.Address(pk.Address()).String())
	// create application object using the message fields
	application := types.NewApplication(sdk.Address(addr), pk, msg.Chains, sdk.ZeroInt())
	ctx.Logger().Info("Validate App Can Stake " + sdk.Address(addr).String())
	// check if they can stake
	if err := k.ValidateApplicationStaking(ctx, application, msg.Value); err != nil {
		ctx.Logger().Error(fmt.Sprintf("Validate App Can Stake Error, at height: %d with address: %s", ctx.BlockHeight(), sdk.Address(addr).String()))
		return err.Result()
	}
	ctx.Logger().Info("Change App state to Staked " + sdk.Address(addr).String())
	// change the application state to staked
	err := k.StakeApplication(ctx, application, msg.Value)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateApplication,
			sdk.NewAttribute(types.AttributeKeyApplication, sdk.Address(addr).String()),
		),
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.Address(addr).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.Address(addr).String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBeginUnstake(ctx sdk.Ctx, msg types.MsgBeginUnstake, k keeper.Keeper) sdk.Result {
	application, found := k.GetApplication(ctx, msg.Address)
	if !found {
		ctx.Logger().Error(fmt.Sprintf("App Not Found at height: %d", ctx.BlockHeight()) + msg.Address.String())
		return types.ErrNoApplicationFound(k.Codespace()).Result()
	}
	if err := k.ValidateApplicationBeginUnstaking(ctx, application); err != nil {
		ctx.Logger().Error(fmt.Sprintf("App Unstake Validation Not Successful, at height: %d", ctx.BlockHeight()) + msg.Address.String())
		return err.Result()
	}
	ctx.Logger().Info("Starting to Unstake App " + msg.Address.String())
	k.BeginUnstakingApplication(ctx, application)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBeginUnstake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Applications must submit a transaction to unjail itself after todo
// having been jailed (and thus unstaked) for downtime
func handleMsgUnjail(ctx sdk.Ctx, msg types.MsgUnjail, k keeper.Keeper) sdk.Result {
	consAddr, err := k.ValidateUnjailMessage(ctx, msg)
	if err != nil {
		return err.Result()
	}
	k.UnjailApplication(ctx, consAddr)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AppAddr.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
