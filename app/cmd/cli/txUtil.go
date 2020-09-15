package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto/keys"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/tendermint/tendermint/libs/rand"

	//"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
)

// SendTransaction - Deliver Transaction to node
func SendTransaction(fromAddr, toAddr, passphrase, chainID string, amount sdk.Int, fees int64, memo string) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, memo)
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// StakeNode - Deliver Stake message to node
func StakeNode(chains []string, serviceURL, fromAddr, passphrase, chainID string, amount sdk.Int, fees int64) (*rpc.SendRawTxParams, error) {
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
	msg := nodeTypes.MsgNodeStake{
		Publickey:  kp.PublicKey.RawString(),
		Chains:     chains,
		Value:      amount,
		ServiceUrl: serviceURL,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// UnstakeNode - start unstaking message to node
func UnstakeNode(fromAddr, passphrase, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

// UnjailNode - Remove node from jail
func UnjailNode(fromAddr, passphrase, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func StakeApp(chains []string, fromAddr, passphrase, chainID string, amount sdk.Int, fees int64) (*rpc.SendRawTxParams, error) {
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
	msg := appsType.MsgApplicationStake{
		PubKey: kp.PublicKey.String(),
		Chains: chains,
		Value:  amount,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func UnstakeApp(fromAddr, passphrase, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func DAOTx(fromAddr, toAddr, passphrase string, amount sdk.Int, action, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func ChangeParam(fromAddr, paramACLKey string, paramValue json.RawMessage, passphrase, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kb, err := app.GetKeybase()
	if err != nil {
		return nil, err
	}

	valueBytes, err := app.Codec().MarshalJSON(paramValue)
	if err != nil {
		return nil, err

	}
	msg := govTypes.MsgChangeParam{
		FromAddress: fa,
		ParamKey:    paramACLKey,
		ParamVal:    valueBytes,
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func Upgrade(fromAddr string, upgrade govTypes.Upgrade, passphrase, chainID string, fees int64) (*rpc.SendRawTxParams, error) {
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
	txBz, err := newTxBz(app.Codec(), &msg, fa, chainID, kb, passphrase, fees, "")
	if err != nil {
		return nil, err
	}
	return &rpc.SendRawTxParams{
		Addr:        fromAddr,
		RawHexBytes: hex.EncodeToString(txBz),
	}, nil
}

func newTxBz(cdc *codec.Codec, msg sdk.Msg, fromAddr sdk.Address, chainID string, keybase keys.Keybase, passphrase string, fee int64, memo string) (transactionBz []byte, err error) {
	// fees
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(fee)))
	// entroyp
	entropy := rand.Int64()
	signBytes, err := auth.StdSignBytes(chainID, entropy, fees, msg, memo)
	if err != nil {
		return nil, err
	}
	sig, pubKey, err := keybase.Sign(fromAddr, passphrase, signBytes)
	if err != nil {
		return nil, err
	}
	s := authTypes.StdSignature{PublicKey: pubKey.RawString(), Signature: sig}
	tx := authTypes.NewTx(msg, fees, s, memo, entropy, cdc.IsAfterUpgrade())
	return auth.DefaultTxEncoder(cdc)(tx)
}
