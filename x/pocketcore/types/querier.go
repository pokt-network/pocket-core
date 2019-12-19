package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryProof                = "proof"
	QueryProofs               = "proofs"
	QuerySupportedBlockchains = "supportedBlockchains"
	QueryRelay                = "relay"
	QueryDispatch             = "dispatch"
	QueryParameters           = "parameters"
)

type QueryRelayParams struct {
	Relay `json:"relay"`
}

type QueryDispatchParams struct {
	SessionHeader `json:"header"`
}

type QueryPORParams struct {
	Address sdk.ValAddress `json:"address"`
	Header  SessionHeader  `json:"header"`
}

type QueryPORsParams struct {
	Address sdk.ValAddress `json:"address"`
}
