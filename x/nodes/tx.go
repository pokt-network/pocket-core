package nodes

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) StakeTx(cdc *codec.Codec, chains map[string]struct{}, serviceURL string, amount sdk.Int, address sdk.ValAddress, passphrase string) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgStake{
		Address:    address,
		PubKey:     am.node.PrivValidator().GetPubKey(),
		Value:      amount,
		ServiceURL: serviceURL, // url where pocket service api is hosted
		Chains:     chains,     // non native blockchains
	}
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, address sdk.ValAddress, passphrase string) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgBeginUnstake{Address: address}
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnjailTx(cdc *codec.Codec, address sdk.ValAddress, passphrase string) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgUnjail{ValidatorAddr: address}
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) Send(cdc *codec.Codec, fromAddr, toAddr sdk.ValAddress, passphrase string, amount sdk.Int) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func newTx(cdc *codec.Codec, am AppModule, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	chainID := am.node.GenesisDoc().ChainID
	fromAddr := sdk.AccAddress(am.node.PrivValidator().GetPubKey().Address())
	cliCtx = util.NewCLIContext(am.node, fromAddr, passphrase).WithCodec(cdc)
	accGetter := auth.NewAccountRetriever(cliCtx)
	err := accGetter.EnsureExists(fromAddr)
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	fee := auth.NewStdFee(9000, sdk.NewCoins(sdk.NewInt64Coin("pokt", 0))) // todo get stake denom
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		account.GetAccountNumber(),
		account.GetSequence(),
		fee.Gas,
		1,
		false,
		chainID,
		"",
		fee.Amount,
		fee.GasPrices(),
	).WithKeybase(am.keybase)
	return
}
