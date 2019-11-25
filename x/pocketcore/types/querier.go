package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryProofsSummary         = "proof_summary"
	QueryProofsSummaries       = "proof_summaries"
	QueryProofsSummariesForApp = "proof_summaries_for_app"
)

type QueryProofSummaryParams struct {
	Address sdk.ValAddress
	Header  ProofsHeader
}

type QueryProofSummariesParams struct {
	Address   sdk.ValAddress
}

type QueryProofSummariesForAppParams struct {
	Address   sdk.ValAddress
	AppPubKey string
}
