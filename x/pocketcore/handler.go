package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
)

// NewHandler returns a handler for "pocketCore" type messages.
func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgProofOfRelays:
			return handleProofBatchMessage(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleProofBatchMessage(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgProofOfRelays) sdk.Result {
	sessionContext := ctx.WithBlockHeight(msg.SessionBlockHeight)
	nPubKey := msg.ProofOfRelay.Proofs[0].ServicerPubKey
	node, found := keeper.GetNodeFromPublicKey(sessionContext, nPubKey)
	if !found {
		return types.NewNodeNotFoundErr(types.ModuleName).Result()
	}
	// get the application at the session context
	app, found := keeper.GetAppFromPublicKey(sessionContext, msg.ProofOfRelay.Proofs[0].Token.ApplicationPublicKey)
	if !found {
		return types.NewAppNotFoundError(types.ModuleName).Result()
	}
	// get all the available service nodes at the time of the session
	allNodes := keeper.GetAllNodes(sessionContext)
	// verify session
	_, err := types.SessionVerification(ctx, node, app, msg.Chain, msg.SessionBlockHeight, int(keeper.SessionNodeCount(ctx)), allNodes)
	if err != nil {
		return err.Result()
	}
	// generate reqProofsIndex to compare
	if err := keeper.ValidateProofs(sessionContext, nPubKey, msg.ProofOfRelay); err != nil {
		return types.NewInvalidProofsError(types.ModuleName).Result()
	}
	// store proofs summary into state
	keeper.SetProofOfRelay(ctx, node.GetAddress(), msg.ProofOfRelay)
	// Award the coins for the batch
	keeper.AwardCoinsForRelays(ctx, msg.TotalRelays, node.GetAddress())
	// create the event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProofBatch,
			sdk.NewAttribute(types.AttributeKeyValidator, node.GetAddress().String()),
		),
	})
	return sdk.Result{Events: ctx.EventManager().Events()}
}
