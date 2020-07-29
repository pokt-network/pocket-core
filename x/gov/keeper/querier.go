package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Ctx, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryACL:
			return queryACL(ctx, k)
		case types.QueryDAO:
			return queryDAO(ctx, k)
		case types.QueryDAOOwner:
			return queryDAOOwner(ctx, k)
		case types.QueryUpgrade:
			return queryUpgrade(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown gov query endpoint")
		}
	}
}

func queryACL(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	acl := k.GetACL(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, acl)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryDAO(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	balance := k.GetDAOTokens(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, balance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryDAOOwner(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	owner := k.GetDAOOwner(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, owner)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryUpgrade(ctx sdk.Ctx, k Keeper) ([]byte, sdk.Error) {
	upgrade := k.GetUpgrade(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, upgrade)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}
