package pos

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

func StakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, chains []string, amount sdk.Int, kp keys.KeyPair, passphrase string) (*sdk.TxResponse, error) {
	fromAddr := kp.GetAddress()
	txBuilder, cliCtx := newTx(cdc, fromAddr, tmNode, keybase, passphrase)
	msg := types.MsgAppStake{
		PubKey: kp.PublicKey,
		Value:  amount,
		Chains: chains, // non native blockchains
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	txBuilder, cliCtx := newTx(cdc, address, tmNode, keybase, passphrase)
	msg := types.MsgBeginAppUnstake{Address: address}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func newTx(cdc *codec.Codec, fromAddr sdk.Address, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	chainID := genDoc.Genesis.ChainID
	cliCtx = util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	cliCtx.BroadcastMode = util.BroadcastSync
	accGetter := auth.NewAccountRetriever(cliCtx)
	err = accGetter.EnsureExists(fromAddr)
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		account.GetAccountNumber(),
		account.GetSequence(),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin("pokt", sdk.NewInt(10)))).WithKeybase(keybase)
	return
}
