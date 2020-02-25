package pos

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

func QueryApplication(cdc *codec.Codec, tmNode client.Client, addr sdk.Address, height int64) (types.Application, error) {
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
	applications := make(types.Applications, 0)
	for _, kv := range resKVs {
		applications = append(applications, types.MustUnmarshalApplication(cdc, kv.Value))
	}
	return applications, nil
}

func QueryStakedApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryStakedApplicationsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryStakedApplications), bz)
	if err != nil {
		return types.Applications{}, err
	}
	apps := types.Applications{}
	err = cdc.UnmarshalJSON(res, &apps)
	if err != nil {
		return apps, err
	}
	return apps, nil
}

func QueryUnstakedApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryStakedApplicationsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUnstakedApplications), bz)
	if err != nil {
		return types.Applications{}, err
	}
	apps := types.Applications{}
	err = cdc.UnmarshalJSON(res, &apps)
	if err != nil {
		return apps, err
	}
	return apps, nil
}

func QueryUnstakingApplications(cdc *codec.Codec, tmNode client.Client, height int64) (types.Applications, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryStakedApplicationsParams{
		Page:  1,
		Limit: 10000,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryUnstakingApplications), bz)
	if err != nil {
		return types.Applications{}, err
	}
	apps := types.Applications{}
	err = cdc.UnmarshalJSON(res, &apps)
	if err != nil {
		return apps, err
	}
	return apps, nil
}

func QuerySupply(cdc *codec.Codec, tmNode client.Client, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAppStakedPool), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryAppUnstakedPool), nil)
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
