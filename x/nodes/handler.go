package nodes

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
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

// These functions assume everything has been authenticated,
// now we just perform action and save
func handleStake(ctx sdk.Context, msg types.MsgStake, k keeper.Keeper) sdk.Result {
	if _, found := k.GetValidator(ctx, sdk.Address(msg.PublicKey.Address())); found {
		return stakeRegisteredValidator(ctx, msg, k)
	} else {
		return stakeNewValidator(ctx, msg, k)
	}
}

func stakeNewValidator(ctx sdk.Context, msg types.MsgStake, k keeper.Keeper) sdk.Result {
	// check to see if teh public key has already been register for that validator
	if _, found := k.GetValidator(ctx, sdk.Address(msg.PublicKey.Address())); found {
		return types.ErrValidatorPubKeyExists(k.Codespace()).Result()
	}
	// check the consensus params
	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PublicKey.PubKey())
		if !common.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return types.ErrValidatorPubKeyTypeNotSupported(k.Codespace(),
				tmPubKey.Type,
				ctx.ConsensusParams().Validator.PubKeyTypes).Result()
		}
	}
	// create validator object using the message fields
	validator := types.NewValidator(sdk.Address(msg.PublicKey.Address()), msg.PublicKey, msg.Chains, msg.ServiceURL, msg.Value)
	validator.Status = sdk.Unstaked
	// check if they can stake
	if err := k.ValidateValidatorStaking(ctx, validator, msg.Value); err != nil {
		return err.Result()
	}
	// register the validator in the world state
	k.RegisterValidator(ctx, validator)
	// change the validator state to staked
	err := k.StakeValidator(ctx, validator, msg.Value)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateValidator,
			sdk.NewAttribute(types.AttributeKeyValidator, sdk.Address(msg.PublicKey.Address()).String()),
		),
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

func stakeRegisteredValidator(ctx sdk.Context, msg types.MsgStake, k keeper.Keeper) sdk.Result {
	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the validator account and global shares are updated within here
	validator, found := k.GetValidator(ctx, sdk.Address(msg.PublicKey.Address()))
	if !found {
		return types.ErrNoValidatorFound(k.Codespace()).Result()
	}
	err := k.ValidateValidatorStaking(ctx, validator, msg.Value)
	if err != nil {
		return err.Result()
	}
	if k.IsWaitingValidator(ctx, validator.Address) {
		return types.ErrValidatorWaitingToUnstake(types.DefaultCodespace).Result()
	}
	err = k.StakeValidator(ctx, validator, msg.Value)
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
	consAddr, err := validateUnjailMessage(ctx, msg, k)
	if err != nil {
		return err.Result()
	}
	k.UnjailValidator(ctx, consAddr)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func validateUnjailMessage(ctx sdk.Context, msg types.MsgUnjail, k keeper.Keeper) (consAddr sdk.Address, err sdk.Error) {
	validator := k.Validator(ctx, msg.ValidatorAddr)
	if validator == nil {
		return nil, types.ErrNoValidatorForAddress(k.Codespace())
	}
	// cannot be unjailed if no self-delegation exists
	selfDel := validator.GetTokens()
	if selfDel == sdk.ZeroInt() {
		return nil, types.ErrMissingSelfDelegation(k.Codespace())
	}
	if validator.GetTokens().LT(sdk.NewInt(k.MinimumStake(ctx))) {
		return nil, types.ErrSelfDelegationTooLowToUnjail(k.Codespace())
	}
	// cannot be unjailed if not jailed
	if !validator.IsJailed() {
		return nil, types.ErrValidatorNotJailed(k.Codespace())
	}
	consAddr = sdk.Address(validator.GetPublicKey().Address())
	info, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		return nil, types.ErrNoValidatorForAddress(k.Codespace())
	}
	// cannot be unjailed if tombstoned
	if info.Tombstoned {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	// cannot be unjailed until out of jail
	if ctx.BlockHeader().Time.Before(info.JailedUntil) {
		return nil, types.ErrValidatorJailed(k.Codespace())
	}
	return
}

func handleMsgSend(ctx sdk.Context, msg types.MsgSend, k keeper.Keeper) sdk.Result {
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
