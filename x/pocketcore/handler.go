package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// NewHandler returns a handler for "pocketCore" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgClaim:
			return handleClaimMsg(ctx, keeper, msg)
		case types.MsgProof:
			return handleProofMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleClaimMsg(ctx sdk.Ctx, k keeper.Keeper, msg types.MsgClaim) sdk.Result {
	// validate the proof message
	if err := k.ValidateClaim(ctx, msg); err != nil {
		return err.Result()
	}
	// set the claim in the world state
	err := k.SetClaim(ctx, msg)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.FromAddress.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleProofMsg(ctx sdk.Ctx, k keeper.Keeper, proof types.MsgProof) sdk.Result {
	// validate the claim claim
	addr, claim, err := k.ValidateProof(ctx, proof)
	if err != nil {
		return err.Result()
	}
	// set the claim in the world state
	er := k.SetReceipt(ctx, addr, types.Receipt{
		SessionHeader:   claim.SessionHeader,
		TotalRelays:     claim.TotalProofs,
		ServicerAddress: addr.String(),
	})
	if er != nil {
		return sdk.ErrInternal(er.Error()).Result()
	}
	// valid claim message so execute according to type
	err = k.ExecuteProof(ctx, proof, claim)
	if err != nil {
		return err.Result()
	}
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProof,
			sdk.NewAttribute(types.AttributeKeyValidator, addr.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
