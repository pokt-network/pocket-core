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
	return (*pcInstance.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).Send(cdc, fa, ta, passphrase, amount)
}

func SendRawTx(fromAddr string, txBytes []byte) (sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return (*pcInstance.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).RawTx(cdc, fa, txBytes)
}

func StakeNode(chains []string, serviceUrl, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).StakeTx(cdc, chains, serviceUrl, amount, fa, passphrase)
}

func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).UnstakeTx(cdc, fa, passphrase)
}

func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(nodeTypes.ModuleName)).(nodes.AppModule).UnjailTx(cdc, fa, passphrase)
}

func StakeApp(chains []string, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(appTypes.ModuleName)).(apps.AppModule).StakeTx(cdc, chains, amount, fa, passphrase)
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.ValAddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(appTypes.ModuleName)).(apps.AppModule).UnstakeTx(cdc, fa, passphrase)
}
