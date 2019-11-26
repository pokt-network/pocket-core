package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) ProofBatchTx(cdc *codec.Codec, cliCtx util.CLIContext, txBuilder auth.TxBuilder, truncatedPOR types.ProofOfRelay) error {
	msg := types.MsgProofOfRelays{
		ProofOfRelay: truncatedPOR,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
