package nodes

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

func StakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, chains []string, serviceURL string, amount sdk.Int, kp keys.KeyPair, passphrase string) (*sdk.TxResponse, error) {
	fromAddr := kp.GetAddress()
	msg := types.MsgStake{
		PublicKey:  kp.PublicKey,
		Value:      amount,
		ServiceURL: serviceURL, // url where pocket service api is hosted
		Chains:     chains,     // non native blockchains
	}
	txBuilder, cliCtx := newTx(cdc, msg, fromAddr, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgBeginUnstake{Address: address}
	txBuilder, cliCtx := newTx(cdc, msg, address, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func UnjailTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgUnjail{ValidatorAddr: address}
	txBuilder, cliCtx := newTx(cdc, msg, address, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func Send(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddr, toAddr sdk.Address, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	txBuilder, cliCtx := newTx(cdc, msg, fromAddr, tmNode, keybase, passphrase)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func RawTx(cdc *codec.Codec, tmNode client.Client, fromAddr sdk.Address, txBytes []byte) (sdk.TxResponse, error) {
	cliCtx := util.CLIContext{
		Codec:       cdc,
		Client:      tmNode,
		FromAddress: fromAddr,
	}
	cliCtx.BroadcastMode = util.BroadcastSync
	return cliCtx.BroadcastTx(txBytes)
}

func newTx(cdc *codec.Codec, msg sdk.Msg, fromAddr sdk.Address, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		panic(err)
	}
	chainID := genDoc.Genesis.ChainID
	cliCtx = util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	cliCtx.BroadcastMode = util.BroadcastSync
	accGetter := auth.NewAccountRetriever(cliCtx)
	err = accGetter.EnsureExists(fromAddr)
	if err != nil {
		panic(err)
	}
	account, err := accGetter.GetAccount(fromAddr)
	if err != nil {
		panic(err)
	}
	fee := sdk.NewInt(types.NodeFeeMap[msg.Type()])
	if account.GetCoins().AmountOf(sdk.DefaultStakeDenom).LTE(fee) { // todo get stake denom
		panic(fmt.Sprintf("insufficient funds: the fee needed is %v", fee))
	}
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		account.GetAccountNumber(),
		account.GetSequence(),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, fee))).WithKeybase(keybase)
	return
}
