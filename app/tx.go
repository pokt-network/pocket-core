package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
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
	return (*app.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).Send(Cdc, fa, ta, passphrase, amount)
}

func StakeNode(chains []string, serviceUrl, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).StakeTx(Cdc, chains, serviceUrl, amount, fa, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).UnstakeTx(Cdc, fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).UnjailTx(Cdc, fa, passphrase)
}

func StakeApp(chains []string, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(appTypes.ModuleName)).(apps.AppModule).StakeTx(Cdc, chains, amount, fa, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(appTypes.ModuleName)).(apps.AppModule).UnstakeTx(Cdc, fa, passphrase)
}
