package types

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/pokt-network/pocket-core/types"
)

// POS params default values
const (
	// DefaultParamspace for params keeper
	DefaultParamspace                 = ModuleName
	DefaultSessionNodeCount           = int64(5)   // default number of nodes in a session
	DefaultClaimSubmissionWindow      = int64(3)   // default sessions to submit a claim
	DefaultClaimExpiration            = int64(100) // default sessions to expire claims
	DefaultReplayAttackBurnMultiplier = int64(3)   // default replay attack burn multiplier
	DefaultMinimumNumberOfProofs      = int64(5)   // default minimum number of proofs
	// DefaultBlockByteSize              = int64(4000000) // default block size in bytes at genesis; 4MB
	DefaultBlockByteSize = int64(100000) // default block size in bytes at genesis; 0.1MB

)

var (
	DefaultSupportedBlockchains   = []string{"0001"}
	KeySessionNodeCount           = []byte("SessionNodeCount")
	KeyClaimSubmissionWindow      = []byte("ClaimSubmissionWindow")
	KeySupportedBlockchains       = []byte("SupportedBlockchains")
	KeyClaimExpiration            = []byte("ClaimExpiration")
	KeyReplayAttackBurnMultiplier = []byte("ReplayAttackBurnMultiplier")
	KeyMinimumNumberOfProofs      = []byte("MinimumNumberOfProofs")
	KeyBlockByteSize              = []byte("BlockByteSize")
)

var _ types.ParamSet = (*Params)(nil)

// "Params" - defines the governance set, high level settings for pocketcore module
type Params struct {
	SessionNodeCount           int64    `json:"session_node_count"`
	ClaimSubmissionWindow      int64    `json:"proof_waiting_period"`
	SupportedBlockchains       []string `json:"supported_blockchains"`
	ClaimExpiration            int64    `json:"claim_expiration"` // per session
	ReplayAttackBurnMultiplier int64    `json:"replay_attack_burn_multiplier"`
	MinimumNumberOfProofs      int64    `json:"minimum_number_of_proofs"`
	BlockByteSize              int64    `json:"block_byte_size,omitempty"`
}

// "ParamSetPairs" - returns an kv params object
// Note: Implements params.ParamSet
func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		{Key: KeySessionNodeCount, Value: &p.SessionNodeCount},
		{Key: KeyClaimSubmissionWindow, Value: &p.ClaimSubmissionWindow},
		{Key: KeySupportedBlockchains, Value: &p.SupportedBlockchains},
		{Key: KeyClaimExpiration, Value: &p.ClaimExpiration},
		{Key: KeyReplayAttackBurnMultiplier, Value: p.ReplayAttackBurnMultiplier},
		{Key: KeyMinimumNumberOfProofs, Value: p.MinimumNumberOfProofs},
		{Key: KeyBlockByteSize, Value: p.BlockByteSize},
	}
}

// "DefaultParams" - Returns a default set of parameters
func DefaultParams() Params {
	defaultGenesisParams := DefaultGenesisParams()
	defaultGenesisParams.BlockByteSize = DefaultBlockByteSize
	return defaultGenesisParams
}

// "DefaultParams" - Returns a default set of parameters at genesis
func DefaultGenesisParams() Params {
	return Params{
		SessionNodeCount:           DefaultSessionNodeCount,
		ClaimSubmissionWindow:      DefaultClaimSubmissionWindow,
		SupportedBlockchains:       DefaultSupportedBlockchains,
		ClaimExpiration:            DefaultClaimExpiration,
		ReplayAttackBurnMultiplier: DefaultReplayAttackBurnMultiplier,
		MinimumNumberOfProofs:      DefaultMinimumNumberOfProofs,
	}
}

// "Validate" - Validate a set of params
func (p Params) Validate() error {
	// session count constraints
	if p.SessionNodeCount > 25 || p.SessionNodeCount < 1 {
		return errors.New("invalid session node count")
	}
	// claim submission window constraints
	if p.ClaimSubmissionWindow < 2 {
		return errors.New("waiting period must be at least 2 sessions")
	}
	// verify each supported blockchain
	for _, chain := range p.SupportedBlockchains {
		if err := NetworkIdentifierVerification(chain); err != nil {
			return err
		}
	}
	// ensure replay attack burn multiplier is above 0
	if p.ReplayAttackBurnMultiplier < 0 {
		return errors.New("invalid replay attack burn multiplier")
	}
	// ensure claim expiration
	if p.ClaimExpiration < 0 {
		return errors.New("invalid claim expiration")
	}
	if p.ClaimExpiration < p.ClaimSubmissionWindow {
		return errors.New("unverified Proof expiration is far too short, must be greater than Proof waiting period")
	}
	return nil
}

// "Equal" - Checks the equality of two param objects
func (p Params) Equal(p2 Params) bool {
	return reflect.DeepEqual(p, p2)
}

// "String" -  returns a human readable string representation of the parameters
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  SessionNodeCount:          %d
  ClaimSubmissionWindow:        %d
  Supported Blockchains      %v
  ClaimExpiration            %d
  ReplayAttackBurnMultiplier %d
  BlockByteSize %d
`,
		p.SessionNodeCount,
		p.ClaimSubmissionWindow,
		p.SupportedBlockchains,
		p.ClaimExpiration,
		p.ReplayAttackBurnMultiplier,
		p.BlockByteSize)
}
