package pocketcore

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth/util"
)

func (am AppModule) QueryProofSummary(cdc *codec.Codec, addr sdk.ValAddress, blockchain, sessionBlockHash, appPubKey string, heightOfQuery int64) (*types.ProofSummary, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(heightOfQuery)
	params := types.QueryProofSummaryParams{
		Address: addr,
		Header: types.ProofsHeader{
			Chain:             blockchain,
			SessionBlockHash:  sessionBlockHash,
			ApplicationPubKey: appPubKey,
		},
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummary), bz)
	var ps types.ProofSummary
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return &ps, nil
}

func (am AppModule) QueryAllProofSummaries(cdc *codec.Codec, addr sdk.ValAddress, height int64) ([]types.ProofSummary, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryProofSummariesParams{
		Address: addr,
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummaries), bz)
	var ps []types.ProofSummary
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (am AppModule) QueryAllProofSummariesForApp(cdc *codec.Codec, addr sdk.ValAddress, appPubKey string, height int64) ([]types.ProofSummary, error) {
	cliCtx := util.NewCLIContext(am.GetTendermintNode(), nil, "").WithCodec(cdc).WithHeight(height)
	params := types.QueryProofSummariesForAppParams{
		Address:   addr,
		AppPubKey: appPubKey,
	}
	bz, err := cdc.MarshalBinaryBare(params)
	if err != nil {
		return nil, err
	}
	proofSummaryBz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.StoreKey, types.QueryProofsSummariesForApp), bz)
	var ps []types.ProofSummary
	err = cdc.UnmarshalJSON(proofSummaryBz, &ps)
	if err != nil {
		return nil, err
	}
	return ps, nil
}
