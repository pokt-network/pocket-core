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
	return nodesModule.QueryAccountBalance(cdc, a, height)
}

func QueryAllNodes(height int64) (nodesTypes.Validators, error) {
	return nodesModule.QueryValidators(cdc, height)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return nodesModule.QueryValidator(cdc, a, height)
}

func QueryUnstakingNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryUnstakingValidators(cdc, height)
}

func QueryStakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryStakedValidators(cdc, height)
}

func QueryUnstakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return nodesModule.QueryUnstakedValidators(cdc, height)
}

func QueryNodeParams(height int64) (params nodesTypes.Params, err error) {
	return nodesModule.QueryPOSParams(cdc, height)
}

func QuerySigningInfo(height int64, addr string) (nodesTypes.ValidatorSigningInfo, error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	return nodesModule.QuerySigningInfo(cdc, height, sdk.ConsAddress(a))
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return nodesModule.QuerySupply(cdc, height)
}

func QueryDaoBalance(height int64) (daoCoins sdk.Int, err error) {
	return nodesModule.QueryDAO(cdc, height)
}

func QueryAllApps(height int64) (appsTypes.Applications, error) {
	return appsModule.QueryApplications(cdc, height)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return appsModule.QueryApplication(cdc, a, height)
}

func QueryUnstakingApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryUnstakingApplications(cdc, height)
}

func QueryStakedApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryStakedApplications(cdc, height)
}

func QueryUnstakedApps(height int64) (validators appsTypes.Applications, err error) {
	return appsModule.QueryUnstakedApplications(cdc, height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return appsModule.QuerySupply(cdc, height)
}

func QuerAppParams(height int64) (params appsTypes.Params, err error) {
	return appsModule.QueryPOSParams(cdc, height)
}

func QueryProofs(addr string, height int64) (proofs []pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocketModule.QueryProofs(cdc, a, height)
}

func QueryProof(blockchain, appPubKey, addr string, sessionblockHeight, height int64) (proof *pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return pocketModule.QueryProof(cdc, a, blockchain, appPubKey, sessionblockHeight, height)
}

func QueryPocketSupportedBlockchains(height int64) ([]string, error){
	return pocketModule.QueryPocketSupportedBlockchains(cdc, height)
}

func QueryPocketParams(height int64) (pocketTypes.Params, error){
	return pocketModule.QueryParams(cdc, height)
}
