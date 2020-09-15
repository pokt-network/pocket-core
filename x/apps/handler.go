package pos

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/keeper"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

var _ sdk.Msg = types.MsgAppStake{}

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.LegacyMsg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgApplicationStake:
			return handleStake(ctx, *msg, k)
		case *types.MsgBeginAppUnstake:
			return handleMsgBeginUnstake(ctx, *msg, k)
		case *types.MsgAppUnjail:
			return handleMsgUnjail(ctx, *msg, k)
		case types.MsgApplicationStake:
			return handleStake(ctx, msg, k)
		case types.MsgBeginAppUnstake:
			return handleMsgBeginUnstake(ctx, msg, k)
		case types.MsgAppUnjail:
			return handleMsgUnjail(ctx, msg, k)
		case types.MsgAppStake:
			return handleLegacyMsgStake(ctx, msg, k)
		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStake(ctx sdk.Ctx, msg types.MsgApplicationStake, k keeper.Keeper) sdk.Result {
	pk, er := crypto.NewPublicKey(msg.PubKey)
	if er != nil {
		return sdk.ErrInvalidPubKey(er.Error()).Result()
	}
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

func handleMsgBeginUnstake(ctx sdk.Ctx, msg types.MsgBeginAppUnstake, k keeper.Keeper) sdk.Result {
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
func handleMsgUnjail(ctx sdk.Ctx, msg types.MsgAppUnjail, k keeper.Keeper) sdk.Result {
	consAddr, err := k.ValidateUnjailMessage(ctx, msg)
	if err != nil {
		return err.Result()
	}
	k.UnjailApplication(ctx, consAddr)
	ctx.EventManager().EmitEvent(
		sdk.Event(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AppAddr.String()),
		)),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Legacy Apps Amino Handlers below
func handleLegacyMsgStake(ctx sdk.Ctx, msg types.MsgAppStake, k keeper.Keeper) sdk.Result {
	if !ctx.IsAfterUpgradeHeight() {
		return handleStake(ctx, msg.ToProto(), k)
	}
	return sdk.ErrInternal("cannot execute a legacy msg: AppStake after upgrade height").Result()
}
