package pocketcore

import (
	"encoding/hex"
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
		case types.MsgProof:
			return handleProofMsg(ctx, keeper, msg)
		case types.MsgClaimProof:
			return handleClaimProofMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleProofMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgProof) sdk.Result {
	// validate the proof message
	if err := validateProofMsg(ctx, keeper, msg); err != nil {
		return err.Result()
	}
	// set the unverified proof in the world state
	keeper.SetUnverifiedProof(ctx, msg)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnverifiedProof,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.FromAddress.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleClaimProofMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgClaimProof) sdk.Result {
	// validate the claim proof
	addr, totalRelays, err := validateClaimProofMsg(ctx, keeper, msg)
	if err != nil {
		return err.Result()
	}
	// valid claim so award coins for relays
	keeper.AwardCoinsForRelays(ctx, totalRelays, addr)
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimedProof,
			sdk.NewAttribute(types.AttributeKeyValidator, addr.String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func validateProofMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgProof) sdk.Error {
	// if is not a pocket supported blockchain then return not supported error
	if !keeper.IsPocketSupportedBlockchain(ctx.WithBlockHeight(msg.SessionBlockHeight), msg.Chain) {
		return types.NewChainNotSupportedErr(types.ModuleName)
	}
	// get the session context
	sessionContext := ctx.WithBlockHeight(msg.SessionBlockHeight)
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
	session, err := types.NewSession(hex.EncodeToString(app.GetConsPubKey().Bytes()), msg.Chain, types.BlockHashFromBlockHeight(ctx, msg.SessionBlockHeight), msg.SessionBlockHeight, allNodes, sessionNodeCount)
	if err != nil {
		return err
	}
	// validate the session
	err = session.Validate(ctx, node, app, sessionNodeCount)
	if err != nil {
		return err
	}
	// check if the proof is ready to be claimed, if it's already ready to be claimed, then it's too late to submit cause the secret is revealed
	if keeper.ProofIsReadyToClaim(ctx, msg.SessionBlockHeight) {
		return types.NewExpiredProofsSubmissionError(types.ModuleName)
	}
	return nil
}

func validateClaimProofMsg(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgClaimProof) (servicerAddr sdk.ValAddress, totalRelays int64, sdkError sdk.Error) {
	// get the public key from the proof
	pk, err := crypto.NewPublicKey(msg.LeafNode.ServicerPubKey)
	if err != nil {
		return nil, 0, types.NewPubKeyError(types.ModuleName, err)
	}
	addr := pk.Address()
	// get the unverified proof for the address
	proof, found := keeper.GetUnverfiedProof(ctx, addr, types.Header{
		ApplicationPubKey:  msg.LeafNode.Token.ApplicationPublicKey,
		Chain:              msg.LeafNode.Blockchain,
		SessionBlockHeight: msg.LeafNode.SessionBlockHeight,
	})
	// if the proof is not found for this claim
	if !found {
		return nil, 0, types.NewUnverifiedProofNotFoundError(types.ModuleName)
	}
	// validate the proof sent
	err = keeper.ValidateProof(ctx, proof, msg)
	if err != nil {
		return nil, 0, types.NewInvalidProofsError(types.ModuleName)
	}
	// seems good, so return the needed info to the handler
	return addr, proof.TotalRelays, nil
}
