package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// creates a querier for staking REST endpoints
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryProof:
			return queryProof(ctx, req, k)
		case types.QueryProofs:
			return queryProofs(ctx, req, k)
		case types.QuerySupportedBlockchains:
			return querySupportedBlockchains(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

// query the supported blockchains
func querySupportedBlockchains(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, k.SupportedBlockchains(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

// query the verified proof object for a specific address and header combination
func queryProof(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPORParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	proofSummary, _ := k.GetProof(ctx, params.Address, params.Header)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, proofSummary)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

// query the verified proof object for a particular node address
func queryProofs(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryPORsParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	proofSummary := k.GetAllProofs(ctx, params.Address)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, proofSummary)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}
