package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) QueryProofOfRelay(cdc *codec.Codec, addr sdk.ValAddress, blockchain, appPubKey string, sessionBlockHeight, heightOfQuery int64) (*types.ProofOfRelay, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(heightOfQuery)
	params := types.QueryPORParams{
		Address: addr,
		Header: types.PORHeader{
			Chain:              blockchain,
			SessionBlockHeight: sessionBlockHeight,
			ApplicationPubKey:  appPubKey,
		},
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummary), bz)
	var ps types.ProofOfRelay
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func (am AppModule) QueryAllPORs(cdc *codec.Codec, addr sdk.ValAddress, height int64) ([]types.ProofOfRelay, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryPORsParams{
		Address: addr,
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummaries), bz)
	var ps []types.ProofOfRelay
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (am AppModule) QueryAllPORsForApp(cdc *codec.Codec, addr sdk.ValAddress, appPubKey string, height int64) ([]types.ProofOfRelay, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryPORsAppParams{
		Address:   addr,
		AppPubKey: appPubKey,
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummariesForApp), bz)
	var ps []types.ProofOfRelay
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (am AppModule) QueryPocketSupportedBlockchains(cdc *codec.Codec, height int64) ([]string, error) {
	var chains []string
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QuerySupportedBlockchains))
	if err != nil {
		return nil, err
	}
	err = cdc.UnmarshalJSON(res, &chains)
	return chains, nil
}
