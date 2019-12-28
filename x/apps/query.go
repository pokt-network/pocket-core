package pos

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

func QueryApplication(cdc *codec.Codec, tmNode client.Client, addr sdk.ValAddress, height int64) (types.Application, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	res, _, err := cliCtx.QueryStore(types.KeyForAppByAllApps(addr), types.StoreKey)
	if err != nil {
		return types.Application{}, err
	}
	if len(res) == 0 {
		return types.Application{}, fmt.Errorf("no application found with address %s", addr)
	}
	return types.MustUnmarshalApplication(cdc, res), nil
}

func QueryApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
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

func QueryStakedApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
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

func QueryUnstakedApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
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

func QueryUnstakingApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
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

func QuerySupply(cdc *codec.Codec, tmNode client.Client, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
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

func QueryPOSParams(cdc *codec.Codec, tmNode client.Client, height int64) (types.Params, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}
