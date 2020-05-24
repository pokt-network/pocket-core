package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/crypto/keys"
	//"github.com/pokt-network/posmint/crypto/keys/mintkey"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	authTypes "github.com/pokt-network/posmint/x/auth/types"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/tendermint/tendermint/libs/common"
)

// SendTransaction - Deliver Transaction to node
func SendTransaction(fromAddr, toAddr, passphrase, chainID string, amount sdk.Int) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	ta, err := sdk.AddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}
	if amount.LTE(sdk.ZeroInt()) {
		return nil, sdk.ErrInternal("must send above 0")
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	msg := nodeTypes.MsgSend{
		FromAddress: fa,
		ToAddress:   ta,
		Amount:      amount,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// StakeNode - Deliver Stake message to node
func StakeNode(chains []string, serviceURL, fromAddr, passphrase, chainID string, amount sdk.Int) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	kp, err := kb.Get(fa)
	if err != nil {
		return nil, err
	}
	for _, chain := range chains {
		err := pocketTypes.NetworkIdentifierVerification(chain)
		if err != nil {
			return nil, err
		}
	}
	if amount.LTE(sdk.NewInt(0)) {
		return nil, sdk.ErrInternal("must stake above zero")
	}
	err = nodeTypes.ValidateServiceURL(serviceURL)
	if err != nil {
		return nil, err
	}
	msg := nodeTypes.MsgStake{
		PublicKey:  kp.PublicKey,
		Chains:     chains,
		Value:      amount,
		ServiceURL: serviceURL,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// UnstakeNode - start unstaking message to node
func UnstakeNode(fromAddr, passphrase, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	msg := nodeTypes.MsgBeginUnstake{
		Address: fa,
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// UnjailNode - Remove node from jail
func UnjailNode(fromAddr, passphrase, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	msg := nodeTypes.MsgUnjail{
		ValidatorAddr: fa,
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func StakeApp(chains []string, fromAddr, passphrase, chainID string, amount sdk.Int) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	kp, err := kb.Get(fa)
	if err != nil {
		return nil, err
	}
	for _, chain := range chains {
		fmt.Println(chain)
		err := pocketTypes.NetworkIdentifierVerification(chain)
		if err != nil {
			return nil, err
		}
	}
	if amount.LTE(sdk.NewInt(0)) {
		return nil, sdk.ErrInternal("must stake above zero")
	}
	msg := appsType.MsgAppStake{
		PubKey: kp.PublicKey,
		Chains: chains,
		Value:  amount,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func UnstakeApp(fromAddr, passphrase, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	msg := appsType.MsgBeginAppUnstake{
		Address: fa,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func DAOTx(fromAddr, toAddr, passphrase string, amount sdk.Int, action, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	ta, err := sdk.AddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	msg := govTypes.MsgDAOTransfer{
		FromAddress: fa,
		ToAddress:   ta,
		Amount:      amount,
		Action:      action,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func ChangeParam(fromAddr, paramACLKey string, paramValue interface{}, passphrase, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	msg := govTypes.MsgChangeParam{
		FromAddress: fa,
		ParamKey:    paramACLKey,
		ParamVal:    paramValue,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func Upgrade(fromAddr string, upgrade govTypes.Upgrade, passphrase, chainID string) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}
	msg := govTypes.MsgUpgrade{
		Address: fa,
		Upgrade: upgrade,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), msg, fa, chainID, kb, passphrase)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func newTxBz(cdc *codec.Codec, msg sdk.Msg, fromAddr sdk.Address, chainID string, keybase keys.Keybase, passphrase string) (transactionBz []byte, err error) {
	if err != nil {
		return nil, err
	}
	// fee
	fee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, msg.GetFee()))
	// entroyp
	entropy := common.RandInt64()
	signBytes, err := auth.StdSignBytes(chainID, entropy, fee, msg, "")
	if err != nil {
		return nil, err
	}
	sig, pubKey, err := keybase.Sign(fromAddr, passphrase, signBytes)
	if err != nil {
		return nil, err
	}
	s := authTypes.StdSignature{PublicKey: pubKey, Signature: sig}
	tx := authTypes.NewStdTx(msg, fee, s, "", entropy)
	return auth.DefaultTxEncoder(cdc)(tx)
}
