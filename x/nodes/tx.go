package nodes

import (
	"fmt"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/auth/util"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/rpc/client"
)

func StakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, chains []string, serviceURL string, amount sdk.Int, kp keys.KeyPair, passphrase string) (*sdk.TxResponse, error) {
	if cdc.IsAfterUpgrade() {
		fromAddr := kp.GetAddress()
		msg := types.MsgNodeStake{
			Publickey:  kp.PublicKey.RawString(),
			Value:      amount,
			ServiceUrl: serviceURL, // url where pocket service api is hosted
			Chains:     chains,     // non native blockchains
		}
		txBuilder, cliCtx, err := newTx(cdc, &msg, fromAddr, tmNode, keybase, passphrase)
		if err != nil {
			return nil, err
		}
		err = msg.ValidateBasic()
		if err != nil {
			return nil, err
		}
		return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, &msg)
	} else {
		return LegacyStakeTx(cdc, tmNode, keybase, chains, serviceURL, amount, kp, passphrase)
	}
}

func UnstakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgBeginUnstake{Address: address}
	txBuilder, cliCtx, err := newTx(cdc, &msg, address, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, &msg)
}

func UnjailTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, address sdk.Address, passphrase string) (*sdk.TxResponse, error) {
	msg := types.MsgUnjail{ValidatorAddr: address}
	txBuilder, cliCtx, err := newTx(cdc, &msg, address, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, &msg)
}

func Send(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, fromAddr, toAddr sdk.Address, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	msg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      amount,
	}
	txBuilder, cliCtx, err := newTx(cdc, &msg, fromAddr, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, &msg)
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
func newTx(cdc *codec.Codec, msg sdk.Msg, fromAddr sdk.Address, tmNode client.Client, keybase keys.Keybase, passphrase string) (txBuilder auth.TxBuilder, cliCtx util.CLIContext, err error) {
	genDoc, err := tmNode.Genesis()
	if err != nil {
		return
	}
	chainID := genDoc.Genesis.ChainID

	kp, err := keybase.Get(fromAddr)
	if err != nil {
		return
	}
	privkey, err := mintkey.UnarmorDecryptPrivKey(kp.PrivKeyArmor, passphrase)
	if err != nil {
		return
	}
	cliCtx = util.NewCLIContext(tmNode, fromAddr, passphrase).WithCodec(cdc)
	cliCtx.BroadcastMode = util.BroadcastSync
	cliCtx.PrivateKey = privkey
	account, err := cliCtx.GetAccount(fromAddr)
	if err != nil {
		return
	}
	fee := msg.GetFee()
	if account.GetCoins().AmountOf(sdk.DefaultStakeDenom).LTE(fee) { // todo get stake denom
		_ = fmt.Errorf("insufficient funds: the fee needed is %v", fee)
		return
	}
	txBuilder = auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"",
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, fee))).WithKeybase(keybase)
	return
}

func LegacyStakeTx(cdc *codec.Codec, tmNode client.Client, keybase keys.Keybase, chains []string, serviceURL string, amount sdk.Int, kp keys.KeyPair, passphrase string) (*sdk.TxResponse, error) {
	fromAddr := kp.GetAddress()
	msg := types.MsgStake{
		PublicKey:  kp.PublicKey,
		Value:      amount,
		ServiceURL: serviceURL, // url where pocket service api is hosted
		Chains:     chains,     // non native blockchains
	}
	txBuilder, cliCtx, err := newTx(cdc, msg, fromAddr, tmNode, keybase, passphrase)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, msg)
}
