package app

import (
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
	return nodesModule.Send(Cdc, fa, ta, passphrase, amount)
}

func StakeNode(chains map[string]struct{}, serviceUrl, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodesModule.StakeTx(Cdc, chains, serviceUrl, amount, fa, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodesModule.UnstakeTx(Cdc, fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return nodesModule.UnjailTx(Cdc, fa, passphrase)
}

func StakeApp(chains map[string]struct{}, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return appsModule.StakeTx(Cdc, chains, amount, fa, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return appsModule.UnstakeTx(Cdc, fa, passphrase)
}
