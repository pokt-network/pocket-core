package mesh

import (
	"encoding/json"
	"github.com/pokt-network/pocket-core/app"
	sdk "github.com/pokt-network/pocket-core/types"
	pocketTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

var (
	InternalStatusType = "internal"
	NotifyStatusType   = "notify"
	ChainStatusType    = "chain"
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
	Success        bool               `json:"success"`
	Error          *SdkErrorResponse  `json:"error"`
	Status         app.HealthResponse `json:"status"`
	Servicers      bool               `json:"servicers"`
	Chains         bool               `json:"Chains"`
	WrongServicers []string           `json:"wrong_servicers"`
	WrongChains    []string           `json:"wrong_chains"`
}
