package nodes

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"reflect"
	"time"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.Msg, signer crypto.PublicKey) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		// convert to value for switch consistency
		if reflect.ValueOf(msg).Kind() == reflect.Ptr {
			msg = reflect.Indirect(reflect.ValueOf(msg)).Interface().(sdk.Msg)
		}
		if k.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
			switch msg := msg.(type) {
			case types.MsgBeginUnstake:
				return handleMsgBeginUnstake(ctx, msg, k)
			case types.MsgUnjail:
				return handleMsgUnjail(ctx, msg, k)
			case types.MsgSend:
				return handleMsgSend(ctx, msg, k)
			case types.MsgStake:
				return handleStake(ctx, msg, k, signer)
			default:
				errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
				return sdk.ErrUnknownRequest(errMsg).Result()
			}
		} else {
			switch msg := msg.(type) {
			case types.LegacyMsgBeginUnstake:
				return legacyHandleMsgBeginUnstake(ctx, msg, k)
			case types.LegacyMsgUnjail:
				return legacyHandleMsgUnjail(ctx, msg, k)
			case types.MsgSend:
				return handleMsgSend(ctx, msg, k)
			case types.LegacyMsgStake:
				return legacyHandleMsgStake(ctx, msg, k, signer)
			default:
				errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
				return sdk.ErrUnknownRequest(errMsg).Result()
			}
		}
	}
}

func handleStake(ctx sdk.Ctx, msg types.MsgStake, k keeper.Keeper, signer crypto.PublicKey) sdk.Result {
	defer sdk.TimeTrack(time.Now())
	pk := msg.PublicKey
	addr := pk.Address()
	// create validator object using the message fields
	validator := types.NewValidator(sdk.Address(addr), pk, msg.Chains, msg.ServiceUrl, sdk.ZeroInt(), msg.Output)
	// check if they can stake
	if err := k.ValidateValidatorStaking(ctx, validator, msg.Value, sdk.Address(signer.Address())); err != nil {
		if sdk.ShowTimeTrackData {
			result := err.Result()
			fmt.Println(result.String())
		}
		return err.Result()
	}
	// change the validator state to staked
	err := k.StakeValidator(ctx, validator, msg.Value, signer)
	if err != nil {
		if sdk.ShowTimeTrackData {
			result := err.Result()
			fmt.Println(result.String())
		}
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
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
	defer sdk.TimeTrack(time.Now())

	ctx.Logger().Info("Begin Unstaking Message received from " + msg.Address.String())
	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	validator, found := k.GetValidator(ctx, msg.Address)
	if !found {
		return types.ErrNoValidatorFound(k.Codespace()).Result()
	}
	err, valid := keeper.ValidateValidatorMsgSigner(validator, msg.Signer, k)
	if !valid {
		return err.Result()
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
func handleMsgUnjail(ctx sdk.Ctx, msg types.MsgUnjail, k keeper.Keeper) sdk.Result {
	defer sdk.TimeTrack(time.Now())

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

func handleMsgSend(ctx sdk.Ctx, msg types.MsgSend, k keeper.Keeper) sdk.Result {
	defer sdk.TimeTrack(time.Now())

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

func legacyHandleMsgBeginUnstake(ctx sdk.Ctx, msg types.LegacyMsgBeginUnstake, k keeper.Keeper) sdk.Result {
	m := types.MsgBeginUnstake{
		Address: msg.Address,
		Signer:  msg.Address,
	}
	return handleMsgBeginUnstake(ctx, m, k)
}

func legacyHandleMsgUnjail(ctx sdk.Ctx, msg types.LegacyMsgUnjail, k keeper.Keeper) sdk.Result {
	m := types.MsgUnjail{
		ValidatorAddr: msg.ValidatorAddr,
		Signer:        msg.ValidatorAddr,
	}
	return handleMsgUnjail(ctx, m, k)
}

func legacyHandleMsgStake(ctx sdk.Ctx, msg types.LegacyMsgStake, k keeper.Keeper, signer crypto.PublicKey) sdk.Result {
	m := types.MsgStake{
		PublicKey:  msg.PublicKey,
		Chains:     msg.Chains,
		Value:      msg.Value,
		ServiceUrl: msg.ServiceUrl,
	}
	return handleStake(ctx, m, k, signer)
}
