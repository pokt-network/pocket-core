package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
)

// POS params default values
const (
	DefaultUnstakingTime                      = time.Hour * 24 * 7 * 3
	DefaultMaxValidators               uint64 = 100000
	DefaultMinStake                    int64  = 1
	DefaultBaseProposerAwardPercentage        = 90
	DefaultMaxEvidenceAge                     = 60 * 2 * time.Second
	DefaultSignedBlocksWindow                 = int64(100)
	DefaultDowntimeJailDuration               = 60 * 10 * time.Second
)

// nolint - Keys for parameter access
var (
	KeyUnstakingTime               = []byte("UnstakingTime")
	KeyMaxValidators               = []byte("MaxValidators")
	KeyStakeDenom                  = []byte("StakeDenom")
	KeyStakeMinimum                = []byte("StakeMinimum")
	KeyProposerRewardPercentage    = []byte("ProposerRewardPercentage")
	KeyMaxEvidenceAge              = []byte("MaxEvidenceAge")
	KeySignedBlocksWindow          = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow          = []byte("MinSignedPerWindow")
	KeyDowntimeJailDuration        = []byte("DowntimeJailDuration")
	KeySlashFractionDoubleSign     = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime       = []byte("SlashFractionDowntime")
	DoubleSignJailEndTime          = time.Unix(253402300799, 0) // forever
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 1)
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(100))
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for pos module
type Params struct {
	UnstakingTime            time.Duration `json:"unstaking_time" yaml:"unstaking_time"`           // duration of unstaking
	MaxValidators            uint64        `json:"max_validators" yaml:"max_validators"`           // maximum number of validators
	StakeDenom               string        `json:"stake_denom" yaml:"stake_denom"`                 // bondable coin denomination
	StakeMinimum             int64         `json:"stake_minimum" yaml:"stake_minimum"`             // minimum amount needed to stake
	ProposerRewardPercentage int8          `json:"base_proposer_award" yaml:"base_proposer_award"` // minimum award for the proposer
	// slashing params
	MaxEvidenceAge          time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window" yaml:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window" yaml:"min_signed_per_window"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration" yaml:"downtime_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUnstakingTime, Value: &p.UnstakingTime},
		{Key: KeyMaxValidators, Value: &p.MaxValidators},
		{Key: KeyStakeDenom, Value: &p.StakeDenom},
		{Key: KeyStakeMinimum, Value: &p.StakeMinimum},
		{Key: KeyMaxEvidenceAge, Value: &p.MaxEvidenceAge},
		{Key: KeySignedBlocksWindow, Value: &p.SignedBlocksWindow},
		{Key: KeyMinSignedPerWindow, Value: &p.MinSignedPerWindow},
		{Key: KeyDowntimeJailDuration, Value: &p.DowntimeJailDuration},
		{Key: KeySlashFractionDoubleSign, Value: &p.SlashFractionDoubleSign},
		{Key: KeySlashFractionDowntime, Value: &p.SlashFractionDowntime},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		UnstakingTime:            DefaultUnstakingTime,
		MaxValidators:            DefaultMaxValidators,
		StakeMinimum:             DefaultMinStake,
		StakeDenom:               sdk.DefaultBondDenom,
		ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
		MaxEvidenceAge:           DefaultMaxEvidenceAge,
		SignedBlocksWindow:       DefaultSignedBlocksWindow,
		MinSignedPerWindow:       DefaultMinSignedPerWindow,
		DowntimeJailDuration:     DefaultDowntimeJailDuration,
		SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
		SlashFractionDowntime:    DefaultSlashFractionDowntime,
	}
}

// validate a set of params
func (p Params) Validate() error {
	if p.StakeDenom == "" {
		return fmt.Errorf("staking parameter StakeDenom can't be an empty string")
	}
	if p.MaxValidators == 0 {
		return fmt.Errorf("staking parameter MaxValidators must be a positive integer")
	}
	if p.StakeMinimum < DefaultMinStake {
		return fmt.Errorf("staking parameter StakeMimimum must be a positive integer")
	}
	if p.ProposerRewardPercentage < 0 || p.ProposerRewardPercentage > 100 {
		return fmt.Errorf("base proposer award is a percentage and must be between 0 and 100")
	}
	return nil
}

// Checks the equality of two param objects
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unstaking Time:          %s
  Max Validators:          %d
  Stake Coin Denom:        %s
  Minimum Stake:     	   %d
  Base Proposer Award:     %d
  MaxEvidenceAge:          %s
  SignedBlocksWindow:      %d
  MinSignedPerWindow:      %s
  DowntimeJailDuration:    %s
  SlashFractionDoubleSign: %s
  SlashFractionDowntime:   %s`,
		p.UnstakingTime,
		p.MaxValidators,
		p.StakeDenom,
		p.StakeMinimum,
		p.ProposerRewardPercentage,
		p.MaxEvidenceAge,
		p.SignedBlocksWindow,
		p.MinSignedPerWindow,
		p.DowntimeJailDuration,
		p.SlashFractionDoubleSign,
		p.SlashFractionDowntime)
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
