package baseapp

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/pokt-network/pocket-core/types"
)

var testQuerier = func(_ sdk.Ctx, _ []string, _ abci.RequestQuery) (res []byte, err sdk.Error) {
	return nil, nil
}
