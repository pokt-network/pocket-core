package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// "NewHandler" - Returns a handler for "pocketCore" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		// handle claim message
		case types.MsgClaim:
			return handleClaimMsg(ctx, keeper, msg)
		// handle proof message
		case types.MsgProof:
			return handleProofMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// "handleClaimMsg" - General handler for the claim message
func handleClaimMsg(ctx sdk.Ctx, k keeper.Keeper, msg types.MsgClaim) sdk.Result {
	// validate the claim message
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

// "handleProofMsg" - General handler for the proof message
func handleProofMsg(ctx sdk.Ctx, k keeper.Keeper, proof types.MsgProof) sdk.Result {
	// validate the claim claim
	addr, claim, err := k.ValidateProof(ctx, proof)
	if err != nil {
		if err.Code() == types.CodeReplayAttackError {
			// if is a replay attack, handle accordingly
			k.HandleReplayAttack(ctx, addr, sdk.NewInt(claim.TotalProofs))
			err := k.DeleteClaim(ctx, addr, claim.SessionHeader, claim.EvidenceType)
			if err != nil {
				ctx.Logger().Error("Could not delete claim from world state after replay attack detected", "Address", claim.FromAddress)
			}
		}
		return err.Result()
	}
	// valid claim message so execute according to type
	err = k.ExecuteProof(ctx, proof, claim)
	if err != nil {
		return err.Result()
	}
	// set the claim in the world state
	er := k.SetReceipt(ctx, addr, types.Receipt{
		SessionHeader:   claim.SessionHeader,
		Total:           claim.TotalProofs,
		ServicerAddress: addr.String(),
		EvidenceType:    proof.Leaf.EvidenceType(),
	})
	if er != nil {
		return sdk.ErrInternal(er.Error()).Result()
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
