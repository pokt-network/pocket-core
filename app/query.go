package app

import (
	"encoding/json"

	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/pokt-network/posmint/x/gov/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// zero for height = latest
func QueryBlock(height *int64) (blockJSON []byte, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err := nodes.QueryBlock(tmClient, height)
	return res, err
}

func QueryTx(hash string) (res *core_types.ResultTx, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryTransaction(tmClient, hash)
	return res, err
}

func QueryAccountTxs(addr string, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryAccountTransactions(tmClient, addr, page, perPage, false, prove)
	return res, err
}
func QueryRecipientTxs(addr string, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryAccountTransactions(tmClient, addr, page, perPage, true, prove)
	return res, err
}

func QueryBlockTxs(height int64, page, perPage int, prove bool) (res *core_types.ResultTxSearch, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryBlockTransactions(tmClient, height, page, perPage, prove)
	return
}

func QueryHeight() (res int64, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryChainHeight(tmClient)
	return
}

func QueryNodeStatus() (res *core_types.ResultStatus, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryNodeStatus(tmClient)
	return
}

func QueryBalance(addr string, height int64) (res sdk.Int, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return sdk.NewInt(0), err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryAccountBalance(Codec(), tmClient, a, height)
	return
}

func QueryAccount(addr string, height int64) (res *auth.BaseAccount, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryAccount(Codec(), tmClient, a, height)
	return
}

func QueryNodes(height int64, opts nodesTypes.QueryValidatorsParams) (res nodesTypes.ValidatorsPage, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryValidators(Codec(), tmClient, height, opts)
	return
}

func QueryNode(addr string, height int64) (res nodesTypes.Validator, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return res, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryValidator(Codec(), tmClient, a, height)
	return
}

func QueryNodeParams(height int64) (res nodesTypes.Params, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QueryPOSParams(Codec(), tmClient, height)
	return
}

func QuerySigningInfo(height int64, addr string) (res nodesTypes.ValidatorSigningInfo, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = nodes.QuerySigningInfo(Codec(), tmClient, height, a)
	return
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	staked, unstaked, err = nodes.QuerySupply(Codec(), tmClient, height)
	return
}

func QueryDaoBalance(height int64) (res sdk.Int, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = gov.QueryDAO(Codec(), tmClient, height)
	return
}

func QueryDaoOwner(height int64) (res sdk.Address, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = gov.QueryDAOOwner(Codec(), tmClient, height)
	return
}

func QueryUpgrade(height int64) (res types.Upgrade, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = gov.QueryUpgrade(Codec(), tmClient, height)
	return
}

func QueryACL(height int64) (res types.ACL, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = gov.QueryACL(Codec(), tmClient, height)
	return
}

func QueryApps(height int64, opts appsTypes.QueryApplicationsWithOpts) (res appsTypes.ApplicationsPage, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = apps.QueryApplications(Codec(), tmClient, height, opts)
	return
}

func QueryApp(addr string, height int64) (res appsTypes.Application, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return res, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = apps.QueryApplication(Codec(), tmClient, a, height)
	return
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	staked, err = apps.QuerySupply(Codec(), tmClient, height)
	return
}

func QueryAppParams(height int64) (res appsTypes.Params, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = apps.QueryPOSParams(Codec(), tmClient, height)
	return
}

func QueryReceipts(addr string, height int64) (res []pocketTypes.Receipt, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryReceipts(Codec(), tmClient, a, height)
	return
}

func QueryReceipt(blockchain, appPubKey, addr, receiptType string, sessionblockHeight, height int64) (res *pocketTypes.Receipt, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryReceipt(Codec(), a, tmClient, blockchain, appPubKey, receiptType, sessionblockHeight, height)
	return
}

func QueryPocketSupportedBlockchains(height int64) (res []string, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryPocketSupportedBlockchains(Codec(), tmClient, height)
	return
}

func QueryPocketParams(height int64) (res pocketTypes.Params, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryParams(Codec(), tmClient, height)
	return
}

func QueryRelay(r pocketTypes.Relay) (res *pocketTypes.RelayResponse, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryRelay(Codec(), tmClient, r)
	return
}

func QueryChallenge(c pocketTypes.ChallengeProofInvalidData) (res *pocketTypes.ChallengeResponse, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryChallenge(Codec(), tmClient, c)
	return
}

func QueryDispatch(header pocketTypes.SessionHeader) (res *pocketTypes.DispatchResponse, err error) {
	tmClient := getTMClient()
	defer tmClient.Stop()
	res, err = pocket.QueryDispatch(Codec(), tmClient, header)
	return
}

func QueryState() (appState json.RawMessage, err error) {
	return pca.ExportAppState(false, nil)
}
