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
		case types.QueryInvoice:
			return queryInvoice(ctx, req, k)
		case types.QueryInvoices:
			return queryInvoices(ctx, req, k)
		case types.QuerySupportedBlockchains:
			return querySupportedBlockchains(ctx, req, k)
		case types.QueryParameters:
			return queryParameters(ctx, k)
		case types.QueryRelay:
			return queryRelay(ctx, req, k)
		case types.QueryDispatch:
			return queryDispatch(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func queryRelay(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryRelayParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	response, er := k.HandleRelay(ctx, params.Relay)
	if er != nil {
		return nil, er
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, response)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

func queryDispatch(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDispatchParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	response, er := k.Dispatch(ctx, params.SessionHeader)
	if er != nil {
		return nil, er
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, *response)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
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

// query the supported blockchains
func querySupportedBlockchains(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, k.SupportedBlockchains(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

// query the verified proof object for a specific address and header combination
func queryInvoice(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryInvoiceParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	invoice, _ := k.GetInvoice(ctx, params.Address, params.Header)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, invoice)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}

// query the verified proof object for a particular node address
func queryInvoices(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryInvoicesParams
	err := types.ModuleCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	invoices, err := k.GetInvoices(ctx, params.Address)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("an error occured retrieving the invoices: %s", err))
	}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, invoices)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to JSON marshal result: %s", err.Error()))
	}
	return res, nil
}
