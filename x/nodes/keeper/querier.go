package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryValidators:
			return queryValidators(ctx, req, k)
		case types.QueryValidator:
			return queryValidator(ctx, req, k)
		case types.QueryUnstakingValidators:
			return queryUnstakingValidators(ctx, req, k)
		case types.QueryStakedValidators:
			return queryStakedValidators(ctx, req, k)
		case types.QueryUnstakedValidators:
			return queryUnstakedValidators(ctx, req, k)
		case types.QuerySigningInfo:
			return querySigningInfo(ctx, req, k)
		case types.QuerySigningInfos:
			return querySigningInfos(ctx, req, k)
		case types.QueryStakedPool:
			return queryStakedPool(ctx, k)
		case types.QueryUnstakedPool:
			return queryUnstakedPool(ctx, k)
		case types.QueryDAO:
			return queryDAO(ctx, k)
		case types.QueryAccountBalance:
			return queryAccountBalance(ctx, req, k)
		case types.QueryAccount:
			return queryAccount(ctx, req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func queryValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validators := k.GetAllValidators(ctx)
	start, end := util.Paginate(len(validators), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		validators = []types.Validator{}
	} else {
		validators = validators[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryAccountBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
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

func queryAccount(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAccountParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	acc := k.GetAccount(ctx, params.Address)
	res, err := k.cdc.MarshalJSON(acc)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryUnstakingValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryUnstakingValidatorsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validators := k.getAllUnstakingValidators(ctx)
	start, end := util.Paginate(len(validators), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		validators = []types.Validator{}
	} else {
		validators = validators[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryStakedValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryStakedValidatorsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validators := k.getStakedValidators(ctx)
	start, end := util.Paginate(len(validators), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		validators = []types.Validator{}
	} else {
		validators = validators[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, validators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryUnstakedValidators(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	validators := k.GetAllValidators(ctx)
	var unstakedValidators = make([]types.Validator, 0)
	for _, v := range validators {
		if v.Status == sdk.Unstaked {
			unstakedValidators = append(unstakedValidators, v)
		}
	}
	start, end := util.Paginate(len(unstakedValidators), params.Page, params.Limit, int(k.GetParams(ctx).MaxValidators))
	if start < 0 || end < 0 {
		unstakedValidators = []types.Validator{}
	} else {
		unstakedValidators = unstakedValidators[start:end]
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, unstakedValidators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func queryValidator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
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

func queryDAO(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	daoPool := k.GetDAOTokens(ctx)
	pool := types.NewPool(daoPool)
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

func querySigningInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySigningInfoParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	signingInfo, found := k.GetValidatorSigningInfo(ctx, params.ConsAddress)
	if !found {
		return nil, types.ErrNoSigningInfoFound(types.DefaultCodespace, params.ConsAddress)
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, signingInfo)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}

	return res, nil
}

func querySigningInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
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
