package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pokt-network/posmint/x/params"
)

// POS params default values
const (
	// DefaultParamspace for params keeper
	DefaultParamspace         = ModuleName
	DefaultSessionNodeCount   = int64(5)
	DefaultProofWaitingPeriod = int64(3)
	DefaultClaimExpiration    = int64(100) // sessions
)

var (
	DefaultSupportedBlockchains []string // todo add defaults
)

// nolint - Keys for parameter access
var (
	KeySessionNodeCount     = []byte("SessionNodeCount")
	KeyProofWaitingPeriod   = []byte("ProofWaitingPeriod")
	KeySupportedBlockchains = []byte("SupportedBlockchains")
	KeyClaimExpiration      = []byte("ClaimExpiration")
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for pos module
type Params struct {
	SessionNodeCount     int64    `json:"session_node_count"`
	ProofWaitingPeriod   int64    `json:"proof_waiting_period"`
	SupportedBlockchains []string `json:"supported_blockchains"`
	ClaimExpiration      int64    `json:"claim_expiration"` // per session
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeySessionNodeCount, Value: &p.SessionNodeCount},
		{Key: KeyProofWaitingPeriod, Value: &p.ProofWaitingPeriod},
		{Key: KeySupportedBlockchains, Value: &p.SupportedBlockchains},
		{Key: KeyClaimExpiration, Value: &p.ClaimExpiration},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		SessionNodeCount:     DefaultSessionNodeCount,
		ProofWaitingPeriod:   DefaultProofWaitingPeriod,
		SupportedBlockchains: DefaultSupportedBlockchains,
		ClaimExpiration:      DefaultClaimExpiration,
	}
}

// validate a set of params
func (p Params) Validate() error {
	if p.SessionNodeCount > 25 || p.SessionNodeCount < 1 {
		return errors.New("Invalid session node count")
	}
	if p.ProofWaitingPeriod < 2 {
		return errors.New("waiting period must be at least 2 sessions")
	}
	if len(p.SupportedBlockchains) == 0 {
		return errors.New("no supported blockchains")
	}
	for _, chain := range p.SupportedBlockchains {
		if err := HashVerification(chain); err != nil {
			return err
		}
	}
	if p.ClaimExpiration < 0 {
		return errors.New("invalid claim expiration")
	}
	if p.ClaimExpiration < p.ProofWaitingPeriod {
		return errors.New("unverified RelayProof expiration is far too short, must be greater than RelayProof waiting period")
	}
	return nil
}

// Checks the equality of two param objects
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// HashString returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  SessionNodeCount:          %d
  ProofWaitingPeriod:        %d
  Supported Blockchains      %v
  ClaimExpiration            %d
`,
		p.SessionNodeCount,
		p.ProofWaitingPeriod,
		p.SupportedBlockchains,
		p.ClaimExpiration)
}
