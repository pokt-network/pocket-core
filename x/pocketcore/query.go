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

// "QueryReceipt" - Exported call to query receipt
func QueryReceipt(cdc *codec.Codec, addr sdk.Address, tmNode client.Client, blockchain, appPubKey, receiptType string, sessionBlockHeight, heightOfQuery int64) (*types.Receipt, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(heightOfQuery)
	// setup params
	params := types.QueryReceiptParams{
		Address: addr,
		Header: types.SessionHeader{
			Chain:              blockchain,
			SessionBlockHeight: sessionBlockHeight,
			ApplicationPubKey:  appPubKey,
		},
		Type: receiptType,
	}
	// marshal params
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	// execute abci query
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryReceipt), bz)
	if err != nil {
		return nil, err
	}
	// unmarshal result
	var ps types.Receipt
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

// "QueryReceipts" - Exported call to query receipts
func QueryReceipts(cdc *codec.Codec, tmNode client.Client, addr sdk.Address, height int64) ([]types.Receipt, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	// setup params
	params := types.QueryReceiptsParams{
		Address: addr,
	}
	// marshal params
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	// execute abci query
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryReceipts), bz)
	// unmarshal result
	var ps []types.Receipt
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

// "QueryParams" - Exported call to query params
func QueryParams(cdc *codec.Codec, tmNode client.Client, height int64) (types.Params, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	// execute abci query
	route := fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryParameters)
	bz, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return types.Params{}, err
	}
	// unmarshal result
	var params types.Params
	cdc.MustUnmarshalJSON(bz, &params)
	return params, nil
}

// "QuerySupportedBlockchains" - Exported call to query sbc
func QueryPocketSupportedBlockchains(cdc *codec.Codec, tmNode client.Client, height int64) ([]string, error) {
	var chains []string
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(height)
	// execute abci query
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QuerySupportedBlockchains))
	if err != nil {
		return nil, err
	}
	// unmarshal result
	err = cdc.UnmarshalJSON(res, &chains)
	if err != nil {
		return nil, err
	}
	return chains, nil
}

// "QueryRelay" - Exported call to execute a relay request
func QueryRelay(cdc *codec.Codec, tmNode client.Client, relay types.Relay) (*types.RelayResponse, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(0)
	// setup params
	params := types.QueryRelayParams{
		Relay: relay,
	}
	// marshal params
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	// execute abci query
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryRelay), bz)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("nil response error")
	}
	// unmarshal result
	var response types.RelayResponse
	err = cdc.UnmarshalJSON(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// "QueryChallenge" - Exported call to execute a challenge report
func QueryChallenge(cdc *codec.Codec, tmNode client.Client, challengeProof types.ChallengeProofInvalidData) (*types.ChallengeResponse, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(0)
	// setup params
	params := types.QueryChallengeParams{
		Challenge: challengeProof,
	}
	// marshal params
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	// execute abci query
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryChallenge), bz)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New("nil response error")
	}
	// unmarshal result
	var response types.ChallengeResponse
	err = cdc.UnmarshalJSON(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// "QueryDispatch" - Exported call to execute a dispatch request
func QueryDispatch(cdc *codec.Codec, tmNode client.Client, header types.SessionHeader) (*types.DispatchResponse, error) {
	// generate cli context
	cliCtx := util.NewCLIContext(tmNode, nil, "").WithCodec(cdc).WithHeight(0)
	params := types.QueryDispatchParams{
		SessionHeader: header,
	}
	// marshal params
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	// execute abci query
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryDispatch), bz)
	if err != nil {
		return nil, err
	}
	// unmarshal result
	var response types.DispatchResponse
	err = cdc.UnmarshalJSON(res, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
