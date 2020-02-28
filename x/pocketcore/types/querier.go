package types

import (
	sdk "github.com/pokt-network/posmint/types"
)

// query endpoints supported by the staking Querier
const (
	QueryReceipt              = "receipt"
	QueryReceipts             = "receipts"
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

type QueryReceiptParams struct {
	Address sdk.Address   `json:"address"`
	Header  SessionHeader `json:"header"`
}

type QueryReceiptsParams struct {
	Address sdk.Address `json:"address"`
}
