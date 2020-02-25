package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryApplications:
			return queryApplications(ctx, req, k)
		case types.QueryApplication:
			return queryApplication(ctx, req, k)
		case types.QueryUnstakingApplications:
			return queryUnstakingApplications(ctx, req, k)
		case types.QueryStakedApplications:
			return queryStakedApplications(ctx, req, k)
		case types.QueryUnstakedApplications:
			return queryUnstakedApplications(ctx, req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryAppStakedPool:
			return queryStakedPool(ctx, k)
		case types.QueryAppUnstakedPool:
			return queryUnstakedPool(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func queryApplications(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAppsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	applications := k.GetAllApplications(ctx)
	start, end := util.Paginate(len(applications), params.Page, params.Limit, int(k.GetParams(ctx).MaxApplications))
	if start < 0 || end < 0 {
		applications = []types.Application{}
	} else {
		applications = applications[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applications)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryUnstakingApplications(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryUnstakingApplicationsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	applications := k.getAllUnstakingApplications(ctx)
	start, end := util.Paginate(len(applications), params.Page, params.Limit, int(k.GetParams(ctx).MaxApplications))
	if start < 0 || end < 0 {
		applications = []types.Application{}
	} else {
		applications = applications[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applications)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryStakedApplications(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryStakedApplicationsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	applications := k.getStakedApplications(ctx)
	start, end := util.Paginate(len(applications), params.Page, params.Limit, int(k.GetParams(ctx).MaxApplications))
	if start < 0 || end < 0 {
		applications = []types.Application{}
	} else {
		applications = applications[start:end]
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applications)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryUnstakedApplications(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAppsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	apps := k.GetAllApplications(ctx)
	unstakedApps := make([]types.Application, 0)
	for _, app := range apps {
		if app.Status == sdk.Unstaked {
			unstakedApps = append(unstakedApps, app)
		}
	}
	start, end := util.Paginate(len(unstakedApps), params.Page, params.Limit, int(k.GetParams(ctx).MaxApplications))
	if start < 0 || end < 0 {
		unstakedApps = []types.Application{}
	} else {
		unstakedApps = unstakedApps[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, unstakedApps)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryApplication(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAppParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	application, found := k.GetApplication(ctx, params.Address)
	if !found {
		return nil, types.ErrNoApplicationFound(types.DefaultCodespace)
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, application)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryStakedPool(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	stakedTokens := k.GetStakedTokens(ctx)
	pool := types.StakingPool(types.NewPool(stakedTokens))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryUnstakedPool(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	unstakedTokens := k.GetUnstakedTokens(ctx)
	pool := types.StakingPool(types.NewPool(unstakedTokens))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}
