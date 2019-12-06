package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/x/params"
)

// POS params default values
const (
	DefaultSessionNodeCount   = 5
	DefaultProofWaitingPeriod = 3
	DefaultClaimExpiration    = 25 // sessions
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
	SessionNodeCount     int64
	ProofWaitingPeriod   int64
	SupportedBlockchains []string
	ClaimExpiration      int64 // per session
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
	if p.SessionNodeCount > 25 {
		return errors.New("too many nodes per session, maximum is 5")
	}
	if p.ProofWaitingPeriod < 1 {
		return errors.New("no waiting period is subject to attack")
	}
	if len(p.SupportedBlockchains) == 0 {
		return errors.New("no supported blockchains")
	}
	if p.ClaimExpiration < p.ProofWaitingPeriod {
		return errors.New("unverified proof expiration is far too short, must be greater than proof waiting period")
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
  ClaimExpiration  			 %d
`,
		p.SessionNodeCount,
		p.ProofWaitingPeriod,
		p.SupportedBlockchains,
		p.ClaimExpiration)
}

// unmarshal the current pos params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	p, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return p
}

// unmarshal the current pos params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}
