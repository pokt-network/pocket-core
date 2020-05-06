package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/pokt-network/posmint/x/gov/types"
)

// SendTransaction - Deliver Transaction to node
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
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.Send(Codec(), tmClient, MustGetKeybase(), fa, ta, passphrase, amount)
	return res, err
}

// SendRawTx - Deliver tx bytes to node
func SendRawTx(fromAddr string, txBytes []byte) (sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.RawTx(Codec(), tmClient, fa, txBytes)
	return res, err
}

// StakeNode - Deliver Stake message to node
func StakeNode(chains []string, serviceURL, fromAddr, passphrase string, amount sdk.Int) (*sdk.TxResponse, error) {
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
	err = nodesTypes.ValidateServiceURL(serviceURL)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.StakeTx(Codec(), tmClient, MustGetKeybase(), chains, serviceURL, amount, kp, passphrase)
	return res, err
}

// UnstakeNode - start unstaking message to node
func UnstakeNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.UnstakeTx(Codec(), tmClient, MustGetKeybase(), fa, passphrase)
	return res, err
}

// UnjailNode - Remove node from jail
func UnjailNode(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.UnjailTx(Codec(), tmClient, MustGetKeybase(), fa, passphrase)
	return res, err
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
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := apps.StakeTx(Codec(), tmClient, MustGetKeybase(), chains, amount, kp, passphrase)
	return res, err
}

func UnstakeApp(fromAddr, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := apps.UnstakeTx(Codec(), tmClient, MustGetKeybase(), fa, passphrase)
	return res, err
}

func DAOTx(fromAddr, toAddr, passphrase string, amount sdk.Int, action string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	ta, err := sdk.AddressFromHex(toAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := gov.DAOTransferTx(Codec(), tmClient, MustGetKeybase(), fa, ta, amount, action, passphrase)
	return res, err
}

func ChangeParam(fromAddr, paramACLKey string, paramValue interface{}, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := gov.ChangeParamsTx(Codec(), tmClient, MustGetKeybase(), fa, paramACLKey, paramValue, passphrase)
	return res, err
}

func Upgrade(fromAddr string, upgrade types.Upgrade, passphrase string) (*sdk.TxResponse, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := gov.UpgradeTx(Codec(), tmClient, MustGetKeybase(), fa, upgrade, passphrase)
	return res, err
}
