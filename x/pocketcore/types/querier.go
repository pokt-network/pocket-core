package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryProof                = "proof"
	QueryProofs               = "proofs"
	QuerySupportedBlockchains = "supportedBlockchains"
)

type QueryPORParams struct {
	Address sdk.ValAddress
	Header  SessionHeader
}

type QueryPORsParams struct {
	Address sdk.ValAddress
}
