package pos

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) StakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, chains map[string]struct{}, amount sdk.Int, address sdk.ValAddress, passphrase string) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgAppStake{
		Address: address,
		PubKey:  am.node.PrivValidator().GetPubKey(),
		Value:   amount,
		Chains:  chains, // non native blockchains
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) UnstakeTx(cdc *codec.Codec, txBuilder auth.TxBuilder, address sdk.ValAddress, passphrase string) error {
	txBuilder, cliCtx := newTx(cdc, am, passphrase)
	msg := types.MsgBeginAppUnstake{Address: address}
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
	fee := auth.NewStdFee(9000, sdk.NewCoins(sdk.NewInt64Coin("pokt", 0)))
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
