package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) TODOTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string, amount sdk.Int) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	kb, err := am.keybase.Get(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	kb = kb
	msg := types.MsgRelayBatch{}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
