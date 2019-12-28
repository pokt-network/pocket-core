package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	sdk "github.com/pokt-network/posmint/types"
)

func SendTransaction(fromAddr, toAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	ta, err := sdk.ValAddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}
	return nodes.Send(GetCodec(), GetTendermintClient(), GetKeybase(), fa, ta, passphrase, amount)
}

func SendRawTx(fromAddr string, txBytes []byte) (sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return nodes.RawTx(GetCodec(), GetTendermintClient(), fa, txBytes)
}

func StakeNode(chains []string, serviceUrl, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kp, err := (GetKeybase()).Get(sdk.AccAddress(fa))
	if err != nil {
		return nil, err
	}
	return nodes.StakeTx(GetCodec(), GetTendermintClient(), GetKeybase(), chains, serviceUrl, amount, kp, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodes.UnstakeTx(GetCodec(), GetTendermintClient(), GetKeybase(), fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodes.UnjailTx(GetCodec(), GetTendermintClient(), GetKeybase(), fa, passphrase)
}

func StakeApp(chains []string, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	kp, err := (GetKeybase()).Get(sdk.AccAddress(fa))
	if err != nil {
		return nil, err
	}
	return apps.StakeTx(GetCodec(), GetTendermintClient(), GetKeybase(), chains, amount, kp, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return apps.UnstakeTx(GetCodec(), GetTendermintClient(), GetKeybase(), fa, passphrase)
}
