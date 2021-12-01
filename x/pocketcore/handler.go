package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"reflect"
	"time"
)

// "NewHandler" - Returns a handler for "pocketCore" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Ctx, msg sdk.Msg, _ crypto.PublicKey) sdk.Result {
		if ctx.IsAfterUpgradeHeight() {
			ctx = ctx.WithEventManager(sdk.NewEventManager())
		}
		// convert to value for switch consistency
		if reflect.ValueOf(msg).Kind() == reflect.Ptr {
			msg = reflect.Indirect(reflect.ValueOf(msg)).Interface().(sdk.Msg)
		}
		switch msg := msg.(type) {
		// handle claim message
		case types.MsgClaim:
			return handleClaimMsg(ctx, keeper, msg)
		// handle legacy proof message
		case types.MsgProof:
			return handleProofMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore ProtoMsg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// "handleClaimMsg" - General handler for the claim message
func handleClaimMsg(ctx sdk.Ctx, k keeper.Keeper, msg types.MsgClaim) sdk.Result {
	defer sdk.TimeTrack(time.Now())
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
	defer sdk.TimeTrack(time.Now())

	// validate the claim claim
	addr, claim, err := k.ValidateProof(ctx, proof)
	if err != nil {
		if err.Code() == types.CodeInvalidMerkleVerifyError && !claim.IsEmpty() {
			// delete local evidence
			processSelf(ctx, k, proof.GetSigners()[0], claim.SessionHeader, claim.EvidenceType, sdk.ZeroInt())
			return err.Result()
		}
		if err.Code() == types.CodeReplayAttackError && !claim.IsEmpty() {
			// delete local evidence
			processSelf(ctx, k, proof.GetSigners()[0], claim.SessionHeader, claim.EvidenceType, sdk.ZeroInt())
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
	tokens, err := k.ExecuteProof(ctx, proof, claim)
	if err != nil {
		return err.Result()
	}
	// delete local evidence
	processSelf(ctx, k, proof.GetSigners()[0], claim.SessionHeader, claim.EvidenceType, tokens)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProof,
			sdk.NewAttribute(types.AttributeKeyValidator, addr.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func processSelf(ctx sdk.Ctx, k keeper.Keeper, signer sdk.Address, header types.SessionHeader, evidenceType types.EvidenceType, tokens sdk.BigInt) {
	// delete local evidence
	if signer.Equals(k.GetSelfAddress(ctx)) {
		err := types.DeleteEvidence(header, evidenceType)
		if err != nil {
			ctx.Logger().Error("Unable to delete evidence: " + err.Error())
		}
		if !tokens.IsZero() {
			types.GlobalServiceMetric().AddUPOKTEarnedFor(header.Chain, float64(tokens.Int64()))
		}
	}
}
