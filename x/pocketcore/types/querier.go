package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// query endpoints supported by the staking Querier
const (
	QueryReceipt              = "receipt"
	QueryReceipts             = "receipts"
	QuerySupportedBlockchains = "supportedBlockchains"
	QueryRelay                = "relay"
	QueryDispatch             = "dispatch"
	QueryChallenge            = "challenge"
	QueryParameters           = "parameters"
)

// "QueryRelayParams" - The parameters needed to submit a relay request
type QueryRelayParams struct {
	Relay `json:"relay"`
}

// "QueryChallengeParams" - The parameters needed to submit a challenge request
type QueryChallengeParams struct {
	Challenge ChallengeProofInvalidData `json:"challengeProof"`
}

// "QueryDispatchParams" - The parameters needed to submit a dispatch request
type QueryDispatchParams struct {
	SessionHeader `json:"header"`
}

// "QueryReceiptParams" - The parameters needed to retrieve a receipt obj for a specific instance
type QueryReceiptParams struct {
	Address sdk.Address   `json:"address"`
	Header  SessionHeader `json:"header"`
	Type    string        `json:"type"`
}

// "QueryReceiptsParama" - The parameters needed to retreive receipt objs for an address
type QueryReceiptsParams struct {
	Address sdk.Address `json:"address"`
}
