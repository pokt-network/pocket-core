package pocketcore

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRelayBatch:
			return handleRelayBatchMessage(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized pocketcore Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleRelayBatchMessage(ctx sdk.Context, keeper Keeper, msg MsgRelayBatch) sdk.Result {
	return sdk.Result{} // return
}
