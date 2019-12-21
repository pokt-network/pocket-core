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
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryBlock(&height)
}

func QueryTx(hash string) (*core_types.ResultTx, error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryTransaction(hash)
}

func QueryHeight() (chainHeight int64, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryChainHeight()
}

func QueryNodeStatus() (*core_types.ResultStatus, error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryNodeStatus()
}

func QueryBalance(addr string, height int64) (balance sdk.Int, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return sdk.NewInt(0), err
	}
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryAccountBalance(Cdc, a, height)
}

func QueryAllNodes(height int64) (nodesTypes.Validators, error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryValidators(Cdc, height)
}

func QueryNode(addr string, height int64) (validator nodesTypes.Validator, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryValidator(Cdc, a, height)
}

func QueryUnstakingNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryUnstakingValidators(Cdc, height)
}

func QueryStakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryStakedValidators(Cdc, height)
}

func QueryUnstakedNodes(height int64) (validators nodesTypes.Validators, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryUnstakedValidators(Cdc, height)
}

func QueryNodeParams(height int64) (params nodesTypes.Params, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryPOSParams(Cdc, height)
}

func QuerySigningInfo(height int64, addr string) (nodesTypes.ValidatorSigningInfo, error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nodesTypes.ValidatorSigningInfo{}, err
	}
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QuerySigningInfo(Cdc, height, sdk.ConsAddress(a))
}

func QueryTotalNodeCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QuerySupply(Cdc, height)
}

func QueryDaoBalance(height int64) (daoCoins sdk.Int, err error) {
	return (*app.mm.GetModule(nodesTypes.ModuleName)).(nodes.AppModule).QueryDAO(Cdc, height)
}

func QueryAllApps(height int64) (appsTypes.Applications, error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryApplications(Cdc, height)
}

func QueryApp(addr string, height int64) (validator appsTypes.Application, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return validator, err
	}
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryApplication(Cdc, a, height)
}

func QueryUnstakingApps(height int64) (validators appsTypes.Applications, err error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryUnstakingApplications(Cdc, height)
}

func QueryStakedApps(height int64) (validators appsTypes.Applications, err error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryStakedApplications(Cdc, height)
}

func QueryUnstakedApps(height int64) (validators appsTypes.Applications, err error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryUnstakedApplications(Cdc, height)
}

func QueryTotalAppCoins(height int64) (staked sdk.Int, unstaked sdk.Int, err error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QuerySupply(Cdc, height)
}

func QueryAppParams(height int64) (params appsTypes.Params, err error) {
	return (*app.mm.GetModule(appsTypes.ModuleName)).(apps.AppModule).QueryPOSParams(Cdc, height)
}

func QueryProofs(addr string, height int64) (proofs []pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryProofs(Cdc, a, height)
}

func QueryProof(blockchain, appPubKey, addr string, sessionblockHeight, height int64) (proof *pocketTypes.StoredProof, err error) {
	a, err := sdk.ValAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryProof(Cdc, a, blockchain, appPubKey, sessionblockHeight, height)
}

func QueryPocketSupportedBlockchains(height int64) ([]string, error) {
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryPocketSupportedBlockchains(Cdc, height)
}

func QueryPocketParams(height int64) (pocketTypes.Params, error) {
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryParams(Cdc, height)
}

func QueryRelay(r pocketTypes.Relay) (*pocketTypes.RelayResponse, error) {
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryRelay(Cdc, r)
}

func QueryDispatch(header pocketTypes.SessionHeader) (*pocketTypes.Session, error) {
	return (*app.mm.GetModule(pocketTypes.ModuleName)).(pocket.AppModule).QueryDispatch(Cdc, header)
}
