package baseapp

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

var testHandler = func(_ sdk.Ctx, _ sdk.Msg) sdk.Result {
	return sdk.Result{}
}
