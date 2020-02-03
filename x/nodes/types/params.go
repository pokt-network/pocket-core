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
	// DefaultParamspace for params keeper
	DefaultParamspace                  = ModuleName
	DefaultUnstakingTime               = time.Hour * 24 * 7 * 3
	DefaultMaxValidators        uint64 = 100000
	DefaultMinStake             int64  = 1000000
	DefaultMaxEvidenceAge              = 60 * 2 * time.Second
	DefaultSignedBlocksWindow          = int64(100)
	DefaultDowntimeJailDuration        = 60 * 10 * time.Second
	DefaultSessionBlocktime            = 25
	DefaultProposerAllocation          = 1
	DefaultDAOAllocation               = 10
)

// nolint - Keys for parameter access
var (
	KeyUnstakingTime               = []byte("UnstakingTime")
	KeyMaxValidators               = []byte("MaxValidators")
	KeyStakeDenom                  = []byte("StakeDenom")
	KeyStakeMinimum                = []byte("StakeMinimum")
	KeyMaxEvidenceAge              = []byte("MaxEvidenceAge")
	KeySignedBlocksWindow          = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow          = []byte("MinSignedPerWindow")
	KeyDowntimeJailDuration        = []byte("DowntimeJailDuration")
	KeySlashFractionDoubleSign     = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime       = []byte("SlashFractionDowntime")
	KeySessionBlock                = []byte("SessionBlockFrequency")
	KeyDAOAllocation               = []byte("DAOAllocation")
	KeyProposerAllocation          = []byte("ProposerPercentage")
	DoubleSignJailEndTime          = time.Unix(253402300799, 0) // forever
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 1)
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(100))
)

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for pos module
type Params struct {
	UnstakingTime         time.Duration `json:"unstaking_time" yaml:"unstaking_time"`                   // how much time must pass between the begin_unstaking_tx and the node going to -> unstaked status
	MaxValidators         uint64        `json:"max_validators" yaml:"max_validators"`                   // maximum number of validators in the network at any given block
	StakeDenom            string        `json:"stake_denom" yaml:"stake_denom"`                         // the monetary denomination of the coins in the network `uPOKT` or `uAtom` or `Wei`
	StakeMinimum          int64         `json:"stake_minimum" yaml:"stake_minimum"`                     // minimum amount of `uPOKT` needed to stake in the network as a node
	SessionBlockFrequency int64         `json:"session_block_frequency" yaml:"session_block_frequency"` // how many blocks are in a session (pocket network unit)
	DAOAllocation         int64         `json:"dao_allocation" yaml:"dao_allocation"`
	ProposerAllocation    int64         `json:"proposer_allocation" yaml:"proposer_allocation"`
	// slashing params
	MaxEvidenceAge          time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`                     // maximum age of tendermint evidence that is still valid (currently not implemented in Cosmos or Pocket-Core)
	SignedBlocksWindow      int64         `json:"signed_blocks_window" yaml:"signed_blocks_window"`             // window of time in blocks (unit) used for signature verification -> specifically in not signing (missing) blocks
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window" yaml:"min_signed_per_window"`           // minimum number of blocks the node must sign per window
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration" yaml:"downtime_jail_duration"`         // minimum amount of time node must spend in jail after missing blocks
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"` // the factor of which a node is slashed for a double sign
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`       // the factor of which a node is slashed for missing blocks
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
		{Key: KeySessionBlock, Value: &p.SessionBlockFrequency},
		{Key: KeyDAOAllocation, Value: &p.DAOAllocation},
		{Key: KeyProposerAllocation, Value: &p.ProposerAllocation},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		UnstakingTime:           DefaultUnstakingTime,
		MaxValidators:           DefaultMaxValidators,
		StakeMinimum:            DefaultMinStake,
		StakeDenom:              sdk.DefaultStakeDenom,
		MaxEvidenceAge:          DefaultMaxEvidenceAge,
		SignedBlocksWindow:      DefaultSignedBlocksWindow,
		MinSignedPerWindow:      DefaultMinSignedPerWindow,
		DowntimeJailDuration:    DefaultDowntimeJailDuration,
		SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
		SlashFractionDowntime:   DefaultSlashFractionDowntime,
		SessionBlockFrequency:   DefaultSessionBlocktime,
		DAOAllocation:           DefaultDAOAllocation,
		ProposerAllocation:      DefaultProposerAllocation,
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
	if p.SessionBlockFrequency < 2 {
		return fmt.Errorf("session block must be greater than 1")
	}
	if p.DAOAllocation < 0 {
		return fmt.Errorf("the dao allocation must not be negative")
	}
	if p.ProposerAllocation < 0 {
		return fmt.Errorf("the proposer allication must not be negative")
	}
	if p.ProposerAllocation+p.DAOAllocation > 100 {
		return fmt.Errorf("the combo of proposer allocation and dao allocation mnust not be greater than 100")
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
  Unstaking Time:          %s
  Max Validators:          %d
  Stake Coin Denom:        %s
  Minimum Stake:     	   %d
  MaxEvidenceAge:          %s
  SignedBlocksWindow:      %d
  MinSignedPerWindow:      %s
  DowntimeJailDuration:    %s
  SlashFractionDoubleSign: %s
  SlashFractionDowntime:   %s
  SessionBlockFrequency    %d
  Proposer Allocation      %d
  DAO allocation           %d`,
		p.UnstakingTime,
		p.MaxValidators,
		p.StakeDenom,
		p.StakeMinimum,
		p.MaxEvidenceAge,
		p.SignedBlocksWindow,
		p.MinSignedPerWindow,
		p.DowntimeJailDuration,
		p.SlashFractionDoubleSign,
		p.SlashFractionDowntime,
		p.SessionBlockFrequency,
		p.ProposerAllocation,
		p.DAOAllocation)
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
