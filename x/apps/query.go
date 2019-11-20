package pos

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (am AppModule) QueryApplication(cdc *codec.Codec, addr sdk.ValAddress, height int64) (types.Application, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	res, _, err := cliCtx.QueryStore(types.KeyForAppByAllApps(addr), types.StoreKey)
	if err != nil {
		return types.Application{}, err
	}
	if len(res) == 0 {
		return types.Application{}, fmt.Errorf("no application found with address %s", addr)
	}
	return types.MustUnmarshalApplication(cdc, res), nil
}

func (am AppModule) QueryApplications(cdc *codec.Codec, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.AllApplicationsKey, types.StoreKey)
	if err != nil {
		return types.Applications{}, err
	}
	var applications types.Applications
	for _, kv := range resKVs {
		applications = append(applications, types.MustUnmarshalApplication(cdc, kv.Value))
	}
	return applications, nil
}

func (am AppModule) QueryStakedApplications(cdc *codec.Codec, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.StakedAppsKey, types.StoreKey)
	if err != nil {
		return types.Applications{}, err
	}
	var applications types.Applications
	for _, kv := range resKVs {
		applications = append(applications, types.MustUnmarshalApplication(cdc, kv.Value))
	}
	return applications, nil
}

func (am AppModule) QueryUnstakedApplications(cdc *codec.Codec, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakedAppsKey, types.StoreKey)
	if err != nil {
		return types.Applications{}, err
	}
	var applications types.Applications
	for _, kv := range resKVs {
		applications = append(applications, types.MustUnmarshalApplication(cdc, kv.Value))
	}
	return applications, nil
}

func (am AppModule) QueryUnstakingApplications(cdc *codec.Codec, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	resKVs, _, err := cliCtx.QuerySubspace(types.UnstakingAppsKey, types.StoreKey)
	if err != nil {
		return types.Applications{}, err
	}
	var applications types.Applications
	for _, kv := range resKVs {
		applications = append(applications, types.MustUnmarshalApplication(cdc, kv.Value))
	}
	return applications, nil
}

func (am AppModule) QuerySupply(cdc *codec.Codec, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/appStakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/appUnstakedPool", types.StoreKey), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var stakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(stakedPoolBytes, &stakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	var unstakedPool types.StakingPool
	if err := cdc.UnmarshalJSON(unstakedPoolBytes, &unstakedPool); err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	return stakedPool.Tokens, unstakedPool.Tokens, nil
}

func (am AppModule) QueryPOSParams(cdc *codec.Codec, height int64) (types.Params, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

func (am AppModule) QueryBlock(height *int64) ([]byte, error) {
	res, err := rpcclient.NewLocal(am.node).Block(height)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSONIndent(res, "", "  ")
}

// get the current blockchain height
func (am AppModule) QueryChainHeight() (int64, error) {
	client := rpcclient.NewLocal(am.node)
	status, err := client.Status()
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

func (am AppModule) QueryNodeStatus() (*ctypes.ResultStatus, error) {
	res, err := rpcclient.NewLocal(am.GetTendermintNode()).Status()
	if err != nil {
		return nil, nil
	}
	return res, nil
}
