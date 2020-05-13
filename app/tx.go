package app

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
)

// SendRawTx - Deliver tx bytes to node
func (app PocketCoreApp) SendRawTx(fromAddr string, txBytes []byte) (sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	cliCtx := util.CLIContext{
		Codec:       cdc,
		Client:      tmClient,
		FromAddress: fa,
	}
	cliCtx.BroadcastMode = util.BroadcastSync
	return cliCtx.BroadcastTx(txBytes)
}
