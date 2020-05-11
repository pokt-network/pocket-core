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
	return types.UnmarshalApplication(cdc, res)
}

func QueryApplications(cdc *codec.Codec, tmNode client.Client, height int64, opts types.QueryApplicationsWithOpts) (types.ApplicationsPage, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	opts.Page, opts.Limit = checkPagination(opts.Page, opts.Limit)
	bz, err := cdc.MarshalJSON(opts)
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

func QuerySupply(cdc *codec.Codec, tmNode client.Client, height int64) (stakedCoins sdk.Int, err error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	stakedPoolBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf(customQuery, types.StoreKey, types.QueryAppStakedPool), nil)
	if err != nil {
		return sdk.Int{}, err
	}
	if err := cdc.UnmarshalJSON(stakedPoolBytes, &stakedCoins); err != nil {
		return sdk.Int{}, err
	}
	return
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
