package pocketcore

import (
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
	"github.com/tendermint/tendermint/rpc/client"
)

func QueryProof(cdc *codec.Codec, addr sdk.Address, tmNode client.Client, blockchain, appPubKey string, sessionBlockHeight, heightOfQuery int64) (*types.StoredInvoice, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(heightOfQuery)
	params := types.QueryInvoiceParams{
		Address: addr,
		Header: types.SessionHeader{
			Chain:              blockchain,
			SessionBlockHeight: sessionBlockHeight,
			ApplicationPubKey:  appPubKey,
		},
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryInvoice), bz)
	if err != nil {
		return nil, err
	}
	var ps types.StoredInvoice
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func QueryProofs(cdc *codec.Codec, tmNode client.Client, addr sdk.Address, height int64) ([]types.StoredInvoice, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryInvoicesParams{
		Address: addr,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryInvoices), bz)
	var ps []types.StoredInvoice
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func QueryParams(cdc *codec.Codec, tmNode client.Client, height int64) (types.Params, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

func QueryPocketSupportedBlockchains(cdc *codec.Codec, tmNode client.Client, height int64) ([]string, error) {
	var chains []string
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QuerySupportedBlockchains))
	if err != nil {
		return nil, err
	}
	err = cdc.UnmarshalJSON(res, &chains)
	if err != nil {
		return nil, err
	}
	return chains, nil
}

func QueryRelay(cdc *codec.Codec, tmNode client.Client, relay types.Relay) (*types.RelayResponse, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(0)
	params := types.QueryRelayParams{
		Relay: relay,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryRelay), bz)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("nil response error")
	}
	var response types.RelayResponse
	err = cdc.UnmarshalJSON(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func QueryDispatch(cdc *codec.Codec, tmNode client.Client, header types.SessionHeader) (*types.Session, error) {
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(0)
	params := types.QueryDispatchParams{
		SessionHeader: header,
	}
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryDispatch), bz)
	if err != nil {
		return nil, err
	}
	var response types.Session
	err = cdc.UnmarshalJSON(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
