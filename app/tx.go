package app

import (
	sdk "github.com/pokt-network/posmint/types"
)

func SendTransaction(fromAddr, toAddr, passphrase string, amount sdk.Int) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	ta, err := sdk.ValAddressFromHex(toAddr)
	if err != nil {
		return err
	}
	return nodesModule.Send(cdc, fa, ta, passphrase, amount)
}

func StakeNode(chains map[string]struct{}, serviceUrl, fromAddr, passphrase string, amount sdk.Int) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	return nodesModule.StakeTx(cdc, chains, serviceUrl, amount, fa, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	return nodesModule.UnstakeTx(cdc, fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	return nodesModule.UnjailTx(cdc, fa, passphrase)
}

func StakeApp(chains map[string]struct{}, serviceUrl, fromAddr, passphrase string, amount sdk.Int) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	return appsModule.StakeTx(cdc, chains, amount, fa, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) error {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return err
	}
	return appsModule.UnstakeTx(cdc, fa, passphrase)
}
