package keeper

import (
	"fmt"
	"math"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/util"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Ctx, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryValidators:
			return queryValidators(ctx, req, k)
		case types.QueryValidator:
			return queryValidator(ctx, req, k)
		case types.QuerySigningInfo:
			return querySigningInfo(ctx, req, k)
		case types.QuerySigningInfos:
			return querySigningInfos(ctx, req, k)
		case types.QueryStakedPool:
			return queryStakedPool(ctx, k)
		case types.QueryAccountBalance:
			return queryAccountBalance(ctx, req, k)
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryTotalSupply:
			return queryTotalSupply(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func paginate(page, limit int, validators []types.Validator, MaxValidators int) types.ValidatorsPage {
	validatorsLen := len(validators)
	start, end := util.Paginate(validatorsLen, page, limit, MaxValidators)

	if start < 0 || end < 0 {
		validators = []types.Validator{}
	} else {
		validators = validators[start:end]
	}
	totalPages := int(math.Ceil(float64(validatorsLen) / float64(end-start)))
	if totalPages < 1 {
		totalPages = 1
	}
	validatorsPage := types.ValidatorsPage{Result: validators, Total: totalPages, Page: page}
	return validatorsPage
}

func queryValidators(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validators := k.GetAllValidatorsWithOpts(ctx, params)
	validatorsPage := paginate(params.Page, params.Limit, validators, int(k.MaxValidators(ctx)))
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validatorsPage)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryAccountBalance(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountBalanceParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	balance := k.GetBalance(ctx, params.Address)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, balance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryAccount(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	acc := k.GetAccount(ctx, params.Address)
	res, err := k.Cdc.MarshalJSON(acc)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryValidator(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validator, found := k.GetValidator(ctx, params.Address)
	if !found {
		return nil, types.ErrNoValidatorFound(types.DefaultCodespace)
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validator)
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

func queryTotalSupply(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	stakedTokens := k.TotalTokens(ctx)
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

func querySigningInfo(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySigningInfoParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	signingInfo, found := k.GetValidatorSigningInfo(ctx, params.Address)
	if !found {
		return nil, types.ErrNoSigningInfoFound(types.DefaultCodespace, params.Address)
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, signingInfo)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func querySigningInfos(ctx sdk.Ctx, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySigningInfosParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	var signingInfos = make([]types.ValidatorSigningInfo, 0)
	k.IterateAndExecuteOverValSigningInfo(ctx, func(consAddr sdk.Address, info types.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})
	start, end := util.Paginate(len(signingInfos), params.Page, params.Limit, int(k.MaxValidators(ctx)))
	if start < 0 || end < 0 {
		signingInfos = []types.ValidatorSigningInfo{}
	} else {
		signingInfos = signingInfos[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, signingInfos)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}
