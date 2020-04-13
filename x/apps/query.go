package pos

import (
	"fmt"

	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

const customQuery = "custom/%s/%s"

func checkPagination(page, limit int) (int, int) {
	if page < 0 {
		page = 1
	}
	if limit < 0 {
		limit = 10000
	}
	return page, limit
}
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

func QueryApplications(cdc *codec.Codec, tmNode client.Client, height int64, page, limit int) (types.ApplicationsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	page, limit = checkPagination(page, limit)
	params := types.QueryStakedApplicationsParams{
		Page:  page,
		Limit: limit,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	ApplicationsPage := types.ApplicationsPage{}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryApplications), bz)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	err = cdc.UnmarshalJSON(res, &ApplicationsPage)
	if err != nil {
		return ApplicationsPage, err
	}

	return ApplicationsPage, nil
}

func QueryStakedApplications(cdc *codec.Codec, tmNode client.Client, height int64, page, limit int) (types.ApplicationsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	page, limit = checkPagination(page, limit)
	params := types.QueryStakedApplicationsParams{
		Page:  page,
		Limit: limit,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryStakedApplications), bz)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	appsPage := types.ApplicationsPage{}
	err = cdc.UnmarshalJSON(res, &appsPage)
	if err != nil {
		return appsPage, err
	}
	return appsPage, nil
}

func QueryUnstakedApplications(cdc *codec.Codec, tmNode client.Client, height int64, page, limit int) (types.ApplicationsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	page, limit = checkPagination(page, limit)
	params := types.QueryStakedApplicationsParams{
		Page:  page,
		Limit: limit,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryUnstakedApplications), bz)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	appsPage := types.ApplicationsPage{}
	err = cdc.UnmarshalJSON(res, &appsPage)
	if err != nil {
		return appsPage, err
	}
	return appsPage, nil
}

func QueryUnstakingApplications(cdc *codec.Codec, tmNode client.Client, height int64, page, limit int) (types.ApplicationsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	page, limit = checkPagination(page, limit)
	params := types.QueryStakedApplicationsParams{
		Page:  page,
		Limit: limit,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryUnstakingApplications), bz)
	if err != nil {
		return types.ApplicationsPage{}, err
	}
	appsPage := types.ApplicationsPage{}
	err = cdc.UnmarshalJSON(res, &appsPage)
	if err != nil {
		return appsPage, err
	}
	return appsPage, nil
}

func QuerySupply(cdc *codec.Codec, tmNode client.Client, height int64) (stakedCoins sdk.Int, unstakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryAppStakedPool), nil)
	if err != nil {
		return sdk.Int{}, sdk.Int{}, err
	}
	unstakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryAppUnstakedPool), nil)
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
	route := fmt.Sprintf(customQuery, types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}
