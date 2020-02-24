package nodes

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgStake:
			return handleStake(ctx, msg, k)
		case types.MsgBeginUnstake:
			return handleMsgBeginUnstake(ctx, msg, k)
		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k)
		case types.MsgSend:
			return handleMsgSend(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStake(ctx sdk.Context, msg types.MsgStake, k keeper.Keeper) sdk.Result {
	// create validator object using the message fields
	validator := types.NewValidator(sdk.Address(msg.PublicKey.Address()), msg.PublicKey, msg.Chains, msg.ServiceURL, sdk.ZeroInt())
	// check if they can stake
	if err := k.ValidateValidatorStaking(ctx, validator, msg.Value); err != nil {
		return err.Result()
	}
	// change the validator state to staked
	err := k.StakeValidator(ctx, validator, msg.Value)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStake,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.Address(msg.PublicKey.Address()).String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.Address(msg.PublicKey.Address()).String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBeginUnstake(ctx sdk.Context, msg types.MsgBeginUnstake, k keeper.Keeper) sdk.Result {
	ctx.Logger().Info("Begin Unstaking Message received from " + msg.Address.String())
	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	validator, found := k.GetValidator(ctx, msg.Address)
	if !found {
		return types.ErrNoValidatorFound(k.Codespace()).Result()
	}
	if err := k.ValidateValidatorBeginUnstaking(ctx, validator); err != nil {
		return err.Result()
	}
	if err := k.WaitToBeginUnstakingValidator(ctx, validator); err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWaitingToBeginUnstaking,
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

// Validators must submit a transaction to unjail itself after todo
// having been jailed (and thus unstaked) for downtime
func handleMsgUnjail(ctx sdk.Context, msg types.MsgUnjail, k keeper.Keeper) sdk.Result {
	ctx.Logger().Info("Unjail Message received from " + msg.ValidatorAddr.String())
	addr, err := k.ValidateUnjailMessage(ctx, msg)
	if err != nil {
		return err.Result()
	}
	k.UnjailValidator(ctx, addr)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSend(ctx sdk.Context, msg types.MsgSend, k keeper.Keeper) sdk.Result {
	ctx.Logger().Info("Send Message from " + msg.FromAddress.String() + " received")
	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}
