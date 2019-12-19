package app

import (
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// zero for height = latest
func QueryBlock(height int64) (blockJSON []byte, err error) {
	return nodesModule.QueryBlock(&height)
}

func QueryTx(hash string) (*core_types.ResultTx, error) {
	return nodesModule.QueryTransaction(hash)
}

func QueryHeight() (chainHeight int64, err error) {
	return nodesModule.QueryChainHeight()
}

func QueryNodeStatus() (*core_types.ResultStatus, error) {
	return nodesModule.QueryNodeStatus()
}

func QueryBalance(addr string, height int64) (balance sdk.Int, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return nodesModule.QueryAccountBalance(Cdc, a, height)
}

func QueryAllNodes(height int64) (nodesTypes.Validators, error) {
	return nodesModule.QueryValidators(Cdc, height)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return nodesModule.QueryValidator(Cdc, a, height)
}

func QueryUnstakingNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryUnstakingValidators(Cdc, height)
}

func QueryStakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryStakedValidators(Cdc, height)
}

func QueryUnstakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryUnstakedValidators(Cdc, height)
}

func QueryNodeParams(height int64) (params nodesTypes.Params, err error) {
	return nodesModule.QueryPOSParams(Cdc, height)
}

func QuerySigningInfo(height int64, addr string) (nodesTypes.ValidatorSigningInfo, error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	return nodesModule.QuerySigningInfo(Cdc, height, sdk.ConsAddress(a))
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return nodesModule.QuerySupply(Cdc, height)
}

func QueryDaoBalance(height int64) (daoCoins sdk.Int, err error) {
	return nodesModule.QueryDAO(Cdc, height)
}

func QueryAllApps(height int64) (appsTypes.Applications, error) {
	return appsModule.QueryApplications(Cdc, height)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return appsModule.QueryApplication(Cdc, a, height)
}

func QueryUnstakingApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryUnstakingApplications(Cdc, height)
}

func QueryStakedApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryStakedApplications(Cdc, height)
}

func QueryUnstakedApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryUnstakedApplications(Cdc, height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return appsModule.QuerySupply(Cdc, height)
}

func QueryAppParams(height int64) (params appsTypes.Params, err error) {
	return appsModule.QueryPOSParams(Cdc, height)
}

func QueryProofs(addr string, height int64) (proofs []pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocketModule.QueryProofs(Cdc, a, height)
}

func QueryProof(blockchain, appPubKey, addr string, sessionblockHeight, height int64) (proof *pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocketModule.QueryProof(Cdc, a, blockchain, appPubKey, sessionblockHeight, height)
}

func QueryPocketSupportedBlockchains(height int64) ([]string, error) {
	return pocketModule.QueryPocketSupportedBlockchains(Cdc, height)
}

func QueryPocketParams(height int64) (pocketTypes.Params, error) {
	return pocketModule.QueryParams(Cdc, height)
}

func QueryRelay(r pocketTypes.Relay) (*pocketTypes.RelayResponse, error) {
	return pocketModule.QueryRelay(Cdc, r)
}

func QueryDispatch(header pocketTypes.SessionHeader) (*pocketTypes.Session, error) {
	return pocketModule.QueryDispatch(Cdc, header)
}
