package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
)

// NewHandler returns a handler for "pocketCore" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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

func handleClaimMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgClaim) sdk.Result {
	// validate the proof message
	if err := validateClaimMsg(ctx, keeper, msg); err != nil {
		return err.Result()
	}
	// set the claim in the world state
	err := keeper.SetClaim(ctx, msg)
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

func handleProofMsg(ctx sdk.Context, k keeper.Keeper, msg types.MsgProof) sdk.Result {
	// validate the claim proof
	addr, proof, err := validateProofMsg(ctx, k, msg)
	if err != nil {
		return err.Result()
	}
	// set the proof in the world state
	er := k.SetInvoice(ctx, addr, types.StoredInvoice{
		SessionHeader:   proof.SessionHeader,
		TotalRelays:     proof.TotalRelays,
		ServicerAddress: addr.String(),
	})
	if er != nil {
		return sdk.ErrInternal(er.Error()).Result()
	}
	// valid claim so award coins for relays
	k.AwardCoinsForRelays(ctx, proof.TotalRelays, addr)
	er = k.DeleteClaim(ctx, addr, proof.SessionHeader)
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

func validateClaimMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgClaim) sdk.Error {
	// if is not a pocket supported blockchain then return not supported error
	if !keeper.IsPocketSupportedBlockchain(ctx.MustGetPrevCtx(msg.SessionBlockHeight), msg.Chain) {
		return types.NewChainNotSupportedErr(types.ModuleName)
	}
	// get the session context
	sessionContext := ctx.MustGetPrevCtx(msg.SessionBlockHeight)
	// get the node from the keeper at the time of the session
	node, found := keeper.GetNode(sessionContext, msg.FromAddress)
	// if not found return not found error
	if !found {
		return types.NewNodeNotFoundErr(types.ModuleName)
	}
	// get the application at the session context
	app, found := keeper.GetAppFromPublicKey(sessionContext, msg.ApplicationPubKey)
	// if not found return not found error
	if !found {
		return types.NewAppNotFoundError(types.ModuleName)
	}
	// get all the available service nodes at the time of the session
	allNodes := keeper.GetAllNodes(sessionContext)
	// get the session node count for the time of thesession
	sessionNodeCount := int(keeper.SessionNodeCount(sessionContext))
	// generate the session
	session, err := types.NewSession(app.GetPublicKey().RawString(), msg.Chain, types.BlockHashFromBlockHeight(ctx, msg.SessionBlockHeight), msg.SessionBlockHeight, allNodes, sessionNodeCount)
	if err != nil {
		return err
	}
	// validate the session
	err = session.Validate(ctx, node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// check if the proof is ready to be claimed, if it's already ready to be claimed, then it's too late to submit cause the secret is revealed
	if keeper.ClaimIsMature(ctx, msg.SessionBlockHeight) {
		return types.NewExpiredProofsSubmissionError(types.ModuleName)
	}
	return nil
}

func validateProofMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgProof) (servicerAddr sdk.Address, proof types.MsgClaim, sdkError sdk.Error) {
	// get the public key from the proof
	pk, err := crypto.NewPublicKey(msg.Leaf.ServicerPubKey)
	if err != nil {
		return nil, types.MsgClaim{}, types.NewPubKeyError(types.ModuleName, err)
	}
	addr := pk.Address()
	// get the unverified proof for the address
	proof, found := keeper.GetClaim(ctx, sdk.Address(addr), types.SessionHeader{
		ApplicationPubKey:  msg.Leaf.Token.ApplicationPublicKey,
		Chain:              msg.Leaf.Blockchain,
		SessionBlockHeight: msg.Leaf.SessionBlockHeight,
	})
	// if the proof is not found for this claim
	if !found {
		return nil, types.MsgClaim{}, types.NewClaimNotFoundError(types.ModuleName)
	}
	// validate the proof sent
	err = keeper.ValidateProof(ctx, proof, msg)
	if err != nil {
		return nil, types.MsgClaim{}, types.NewInvalidProofsError(types.ModuleName)
	}
	// seems good, so return the needed info to the handler
	return sdk.Address(addr), proof, nil
}
