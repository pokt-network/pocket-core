package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"net/url"
)

func SendTransaction(fromAddr, toAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
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
	return nodes.Send(Codec(), getTMClient(), MustGetKeybase(), fa, ta, passphrase, amount)
}

func SendRawTx(fromAddr string, txBytes []byte) (sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return nodes.RawTx(Codec(), getTMClient(), fa, txBytes)
}

func StakeNode(chains []string, serviceUrl, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kp, err := (MustGetKeybase()).Get(fa)
	if err != nil {
		return nil, err
	}
	for _, chain := range chains {
		err := pocketTypes.HashVerification(chain)
		if err != nil {
			return nil, err
		}
	}
	if amount.LTE(sdk.NewInt(0)) {
		return nil, sdk.ErrInternal("must stake above zero")
	}
	_, err = url.ParseRequestURI(serviceUrl)
	if err != nil {
		return nil, err
	}
	return nodes.StakeTx(Codec(), getTMClient(), MustGetKeybase(), chains, serviceUrl, amount, kp, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodes.UnstakeTx(Codec(), getTMClient(), MustGetKeybase(), fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodes.UnjailTx(Codec(), getTMClient(), MustGetKeybase(), fa, passphrase)
}

func StakeApp(chains []string, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kp, err := (MustGetKeybase()).Get(fa)
	if err != nil {
		return nil, err
	}
	for _, chain := range chains {
		err := pocketTypes.HashVerification(chain)
		if err != nil {
			return nil, err
		}
	}
	if amount.LTE(sdk.NewInt(0)) {
		return nil, sdk.ErrInternal("must stake above zero")
	}
	return apps.StakeTx(Codec(), getTMClient(), MustGetKeybase(), chains, amount, kp, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return apps.UnstakeTx(Codec(), getTMClient(), MustGetKeybase(), fa, passphrase)
}
