package pocketcore

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO
// NewHandler returns a handler for "pocketCore" type messages.
func NewHandler(keeper PocketCoreKeeper) sdk.Handler {
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
func handleRelayBatchMessage(ctx sdk.Context, keeper PocketCoreKeeper, msg MsgRelayBatch) sdk.Result {
	return sdk.Result{} // return
}
