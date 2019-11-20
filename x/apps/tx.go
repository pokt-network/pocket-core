package pos

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) StakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, chains map[string]struct{}, amount sdk.Int, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	kb, err := am.keybase.Get(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	msg := types.MsgAppStake{
		Address: address,
		PubKey:  kb.PubKey, // needed for validator creation
		Value:   amount,
		Chains:  chains, // non native blockchains
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgBeginAppUnstake{Address: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

//func (am AppModule) UnjailTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
//	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
//	msg := types.MsgAppUnjail{AppAddr: address}
//	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
//}
