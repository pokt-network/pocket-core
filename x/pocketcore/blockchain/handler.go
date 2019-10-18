package blockchain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// todo account handlers

// NewHandler returns a handler for "blockchain" type messages.
func NewHandler(nodeKeeper NodeKeeper, applicationKeeper ApplicationKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRelayBatch:
			return handleRelayBatchMessage(ctx, nodeKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized blockchain Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleRelayBatchMessage(ctx sdk.Context, keeper NodeKeeper, msg MsgRelayBatch) sdk.Result {
	return sdk.Result{} // return
}
