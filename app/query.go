package app

import (
	apps "github.com/pokt-network/pocket-core/x/apps"
	appsTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
	sdk "github.com/pokt-network/posmint/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

// zero for height = latest
func QueryBlock(height int64) (blockJSON []byte, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryBlock(&height)
}

func QueryTx(hash string) (*core_types.ResultTx, error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryTransaction(hash)
}

func QueryHeight() (chainHeight int64, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryChainHeight()
}

func QueryNodeStatus() (*core_types.ResultStatus, error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryNodeStatus()
}

func QueryBalance(addr string, height int64) (balance sdk.Int, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryAccountBalance(cdc, a, height)
}

func QueryAllNodes(height int64) (nodesTypes.Validators, error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryValidators(cdc, height)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryValidator(cdc, a, height)
}

func QueryUnstakingNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryUnstakingValidators(cdc, height)
}

func QueryStakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryStakedValidators(cdc, height)
}

func QueryUnstakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryUnstakedValidators(cdc, height)
}

func QueryNodeParams(height int64) (params nodesTypes.Params, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryPOSParams(cdc, height)
}

func QuerySigningInfo(height int64, addr string) (nodesTypes.ValidatorSigningInfo, error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QuerySigningInfo(cdc, height, sdk.ConsAddress(a))
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QuerySupply(cdc, height)
}

func QueryDaoBalance(height int64) (daoCoins sdk.Int, err error) {
	return (*pcInstance.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryDAO(cdc, height)
}

func QueryAllApps(height int64) (appsTypes.Applications, error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryApplications(cdc, height)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryApplication(cdc, a, height)
}

func QueryUnstakingApps(height int64) (validators appsTypes.Applications, err error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryUnstakingApplications(cdc, height)
}

func QueryStakedApps(height int64) (validators appsTypes.Applications, err error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryStakedApplications(cdc, height)
}

func QueryUnstakedApps(height int64) (validators appsTypes.Applications, err error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryUnstakedApplications(cdc, height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QuerySupply(cdc, height)
}

func QueryAppParams(height int64) (params appsTypes.Params, err error) {
	return (*pcInstance.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryPOSParams(cdc, height)
}

func QueryProofs(addr string, height int64) (proofs []pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryProofs(cdc, a, height)
}

func QueryProof(blockchain, appPubKey, addr string, sessionblockHeight, height int64) (proof *pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryProof(cdc, a, blockchain, appPubKey, sessionblockHeight, height)
}

func QueryPocketSupportedBlockchains(height int64) ([]string, error) {
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryPocketSupportedBlockchains(cdc, height)
}

func QueryPocketParams(height int64) (pocketTypes.Params, error) {
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryParams(cdc, height)
}

func QueryRelay(r pocketTypes.Relay) (*pocketTypes.RelayResponse, error) {
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryRelay(cdc, r)
}

func QueryDispatch(header pocketTypes.SessionHeader) (*pocketTypes.Session, error) {
	return (*pcInstance.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryDispatch(cdc, header)
}
