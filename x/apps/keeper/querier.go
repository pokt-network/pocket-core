package keeper

import (
	"fmt"
	"math"

	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	abci "github.com/tendermint/tendermint/abci/types"
)

func paginate(page, limit int, validators []types.Application, MaxValidators int) types.ApplicationsPage {
	validatorsLen := len(validators)
	start, end := util.Paginate(validatorsLen, page, limit, MaxValidators)

	if start < 0 || end < 0 {
		validators = []types.Application{}
	} else {
		validators = validators[start:end]
	}
	totalPages := int(math.Ceil(float64(validatorsLen) / float64(end-start)))
	if totalPages < 1 {
		totalPages = 1
	}
	applicationsPage := types.ApplicationsPage{Result: validators, Total: totalPages, Page: page}
	return applicationsPage
}

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Ctx, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
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

func queryApplications(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAppsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	applications := k.GetAllApplications(ctx)
	applicationsPage := paginate(params.Page, params.Limit, applications, int(k.GetParams(ctx).MaxApplications))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applicationsPage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryUnstakingApplications(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryUnstakingApplicationsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	applications := k.getAllUnstakingApplications(ctx)
	applicationsPage := paginate(params.Page, params.Limit, applications, int(k.GetParams(ctx).MaxApplications))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applicationsPage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryStakedApplications(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryStakedApplicationsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	applications := k.getStakedApplications(ctx)
	applicationsPage := paginate(params.Page, params.Limit, applications, int(k.GetParams(ctx).MaxApplications))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applicationsPage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryUnstakedApplications(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
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
	unstakedAppsPage := paginate(params.Page, params.Limit, unstakedApps, int(k.GetParams(ctx).MaxApplications))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, unstakedAppsPage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryApplication(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
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

func queryStakedPool(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	stakedTokens := k.GetStakedTokens(ctx)
	pool := types.StakingPool(types.NewPool(stakedTokens))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryUnstakedPool(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	unstakedTokens := k.GetUnstakedTokens(ctx)
	pool := types.StakingPool(types.NewPool(unstakedTokens))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, pool)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}

func queryParameters(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	params := k.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return res, nil
}
