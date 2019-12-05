package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryProofsSummary         = "proof_summary"
	QueryProofsSummaries       = "proof_summaries"
	QueryProofsSummariesForApp = "proof_summaries_for_app"
	QuerySupportedBlockchains  = "supportedBlockchains"
)

type QueryPORParams struct {
	Address sdk.ValAddress
	Header  Header
}

type QueryPORsParams struct {
	Address sdk.ValAddress
}

type QueryPORsAppParams struct {
	Address   sdk.ValAddress
	AppPubKey string
}
