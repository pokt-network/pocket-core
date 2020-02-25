package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// zero for height = latest
func QueryBlock(height *int64) (blockJSON []byte, err error) {
	return nodes.QueryBlock(getTMClient(), height)
}

func QueryTx(hash string) (*core_types.ResultTx, error) {
	return nodes.QueryTransaction(getTMClient(), hash)
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

func QueryAllNodes(height int64) (nodesTypes.Validators, error) {
	return nodes.QueryValidators(Codec(), getTMClient(), height)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return nodes.QueryValidator(Codec(), getTMClient(), a, height)
}

func QueryUnstakingNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodes.QueryUnstakingValidators(Codec(), getTMClient(), height)
}

func QueryStakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodes.QueryStakedValidators(Codec(), getTMClient(), height)
}

func QueryUnstakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodes.QueryUnstakedValidators(Codec(), getTMClient(), height)
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
	return nodes.QueryDAO(Codec(), getTMClient(), height)
}

func QueryAllApps(height int64) (appsTypes.Applications, error) {
	return apps.QueryApplications(Codec(), getTMClient(), height)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return apps.QueryApplication(Codec(), getTMClient(), a, height)
}

func QueryUnstakingApps(height int64) (validators appsTypes.Applications, err error) {
	return apps.QueryUnstakingApplications(Codec(), getTMClient(), height)
}

func QueryStakedApps(height int64) (validators appsTypes.Applications, err error) {
	return apps.QueryStakedApplications(Codec(), getTMClient(), height)
}

func QueryUnstakedApps(height int64) (validators appsTypes.Applications, err error) {
	return apps.QueryUnstakedApplications(Codec(), getTMClient(), height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return apps.QuerySupply(Codec(), getTMClient(), height)
}

func QueryAppParams(height int64) (params appsTypes.Params, err error) {
	return apps.QueryPOSParams(Codec(), getTMClient(), height)
}

func QueryProofs(addr string, height int64) (proofs []pocketTypes.StoredInvoice, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocket.QueryProofs(Codec(), getTMClient(), a, height)
}

func QueryProof(blockchain, appPubKey, addr string, sessionblockHeight, height int64) (proof *pocketTypes.StoredInvoice, err error) {
	a, err := sdk.AddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocket.QueryProof(Codec(), a, getTMClient(), blockchain, appPubKey, sessionblockHeight, height)
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

func QueryDispatch(header pocketTypes.SessionHeader) (*pocketTypes.Session, error) {
	return pocket.QueryDispatch(Codec(), getTMClient(), header)
}
