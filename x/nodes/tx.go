package nodes

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) StakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, chains map[string]struct{}, serviceURL string, amount sdk.Int, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	kb, err := am.keybase.Get(sdk.AccAddress(address))
	if err != nil {
		return err
	}
	msg := types.MsgStake{
		Address:    address,
		PubKey:     kb.PubKey, // needed for validator creation
		Value:      amount,
		ServiceURL: serviceURL, // url where pocket service api is hosted
		Chains:     chains,     // non native blockchains
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgBeginUnstake{Address: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnjailTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(address), passphrase).WithCodec(cdc)
	msg := types.MsgUnjail{ValidatorAddr: address}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) Send(cdc *codec.Codec, fromAddr, toAddr sdk.ValAddress, txBuilder auth.TxBuilder, passphrase string, amount sdk.Int) error {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), sdk.AccAddress(fromAddr), passphrase).WithCodec(cdc)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}
