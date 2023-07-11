package mesh

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

var (
	ChainNotFoundStatusType     = "chain_not_found"
	ServicerNotFoundStatusType  = "servicer_not_found"
	BadRequest                  = "bad_request"
	NotifyRequestErrorType      = "notify_request"
	NotifyResponseErrorType     = "notify_response"
	ChainStatusType             = "chain"
	GetSessionErrorType         = "get_session"
	InvalidSessionType          = "invalid_session"
	SessionHeightOutOfRangeType = "session_height_out_of_range"
	AuthorizationHeader         = "Authorization"
)

// SdkErrorResponse - response error format for re-implemented endpoints.
type SdkErrorResponse struct {
	Code      sdk.CodeType      `json:"code"`
	Codespace sdk.CodespaceType `json:"codespace"`
	Error     string            `json:"message"`
}

// HealthResponse - response payload of /v1/mesh/health
type HealthResponse struct {
	Version   string `json:"version"`
	Servicers int    `json:"servicers"`
	FullNodes int    `json:"full_nodes"`
}

// RPCRelayResult response payload of /v1/client/relay
type RPCRelayResult struct {
	Success  bool                          `json:"signature"`
	Error    error                         `json:"error"`
	Dispatch *pocketTypes.DispatchResponse `json:"dispatch"`
}

// RPCSessionResult - response payload of /v1/private/mesh/session
type RPCSessionResult struct {
	Success         bool              `json:"success"`
	Error           *SdkErrorResponse `json:"error"`
	Dispatch        *DispatchResponse `json:"dispatch"`
	RemainingRelays json.Number       `json:"remaining_relays"`
}

// RPCRelayResponse - response payload of /v1/private/mesh/relay
type RPCRelayResponse struct {
	Success  bool              `json:"signature"`
	Error    *SdkErrorResponse `json:"error"`
	Dispatch *DispatchResponse `json:"dispatch"`
}

// CheckPayload - payload used to call /v1/private/mesh/check
type CheckPayload struct {
	Servicers []string `json:"servicers"`
	Chains    []string `json:"Chains"`
}

// CheckResponse - response payload of /v1/private/mesh/check
type CheckResponse struct {
	Success                    bool               `json:"success"`
	Error                      *SdkErrorResponse  `json:"error"`
	Status                     app.HealthResponse `json:"status"`
	BlocksPerSession           int64              `json:"blocks_per_session"`
	Servicers                  bool               `json:"servicers"`
	Chains                     bool               `json:"chains"`
	WrongServicers             []string           `json:"wrong_servicers"`
	WrongChains                []string           `json:"wrong_chains"`
	ClientSessionSyncAllowance int64              `json:"client_session_sync_allowance"`
}
