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
	return nodes.QueryBlock(getTMClient(), height)
}

func QueryTx(hash string) (*core_types.ResultTx, error) {
	return nodes.QueryTransaction(getTMClient(), hash)
}

func QueryAccountTxs(addr string, page, perPage int) (*core_types.ResultTxSearch, error) {
	return nodes.QueryAccountTransactions(getTMClient(), addr, page, perPage)
}

func QueryBlockTxs(height int64, page, perPage int) (*core_types.ResultTxSearch, error) {
	return nodes.QueryBlockTransactions(getTMClient(), height, page, perPage)
}

func QueryHeight() (chainHeight int64, err error) {
	return nodes.QueryChainHeight(getTMClient())
}

func QueryNodeStatus() (*core_types.ResultStatus, error) {
	return nodes.QueryNodeStatus(getTMClient())
}

func QueryBalance(addr string, height int64) (balance sdk.Int, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return nodes.QueryAccountBalance(Codec(), getTMClient(), a, height)
}

func QueryAccount(addr string, height int64) (account *auth.BaseAccount, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return nodes.QueryAccount(Codec(), getTMClient(), a, height)
}

func QueryNodes(height int64, opts nodesTypes.QueryValidatorsParams) (nodesTypes.ValidatorsPage, error) {
	return nodes.QueryValidators(Codec(), getTMClient(), height, opts)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return nodes.QueryValidator(Codec(), getTMClient(), a, height)
}

func QueryNodeParams(height int64) (params nodesTypes.Params, err error) {
	return nodes.QueryPOSParams(Codec(), getTMClient(), height)
}

func QuerySigningInfo(height int64, addr string) (nodesTypes.ValidatorSigningInfo, error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	return nodes.QuerySigningInfo(Codec(), getTMClient(), height, a)
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return nodes.QuerySupply(Codec(), getTMClient(), height)
}

func QueryDaoBalance(height int64) (daoCoins sdk.Int, err error) {
	return gov.QueryDAO(Codec(), getTMClient(), height)
}

func QueryDaoOwner(height int64) (daoOwner sdk.Address, err error) {
	return gov.QueryDAOOwner(Codec(), getTMClient(), height)
}

func QueryUpgrade(height int64) (upgrade types.Upgrade, err error) {
	return gov.QueryUpgrade(Codec(), getTMClient(), height)
}

func QueryACL(height int64) (acl types.ACL, err error) {
	return gov.QueryACL(Codec(), getTMClient(), height)
}

func QueryApps(height int64, opts appsTypes.QueryApplicationsWithOpts) (appsTypes.ApplicationsPage, error) {
	return apps.QueryApplications(Codec(), getTMClient(), height, opts)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return apps.QueryApplication(Codec(), getTMClient(), a, height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return apps.QuerySupply(Codec(), getTMClient(), height)
}

func QueryAppParams(height int64) (params appsTypes.Params, err error) {
	return apps.QueryPOSParams(Codec(), getTMClient(), height)
}

func QueryReceipts(addr string, height int64) (proofs []pocketTypes.Receipt, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocket.QueryReceipts(Codec(), getTMClient(), a, height)
}

func QueryReceipt(blockchain, appPubKey, addr, receiptType string, sessionblockHeight, height int64) (proof *pocketTypes.Receipt, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocket.QueryReceipt(Codec(), a, getTMClient(), blockchain, appPubKey, receiptType, sessionblockHeight, height)
}

func QueryPocketSupportedBlockchains(height int64) ([]string, error) {
	return pocket.QueryPocketSupportedBlockchains(Codec(), getTMClient(), height)
}

func QueryPocketParams(height int64) (pocketTypes.Params, error) {
	return pocket.QueryParams(Codec(), getTMClient(), height)
}

func QueryRelay(r pocketTypes.Relay) (*pocketTypes.RelayResponse, error) {
	return pocket.QueryRelay(Codec(), getTMClient(), r)
}

func QueryChallenge(c pocketTypes.ChallengeProofInvalidData) (*pocketTypes.ChallengeResponse, error) {
	return pocket.QueryChallenge(Codec(), getTMClient(), c)
}

func QueryDispatch(header pocketTypes.SessionHeader) (*pocketTypes.DispatchResponse, error) {
	return pocket.QueryDispatch(Codec(), getTMClient(), header)
}

func QueryState() (appState json.RawMessage, err error) {
	return pca.ExportAppState(false, nil)
}
