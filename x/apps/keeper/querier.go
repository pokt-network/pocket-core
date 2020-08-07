package keeper

import (
	"fmt"
	"math"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth/util"
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

// NewQuerier - creates a query router for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Ctx, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryApplications:
			return queryApplications(ctx, req, k)
		case types.QueryApplication:
			return queryApplication(ctx, req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryAppStakedPool:
			return queryStakedPool(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func queryApplications(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryApplicationsWithOpts
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	applications := k.GetAllApplicationsWithOpts(ctx, params)
	applicationsPage := paginate(params.Page, params.Limit, applications, int(k.GetParams(ctx).MaxApplications))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, applicationsPage)
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
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, stakedTokens)
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
