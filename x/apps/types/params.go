package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"math"
	"reflect"
	"time"
)

// POS params default values
const (
	// DefaultParamspace for params keeper
	DefaultParamspace                = ModuleName
	DefaultUnstakingTime             = time.Hour * 24 * 7 * 3
	DefaultMaxApplications     int64 = math.MaxInt64
	DefaultMinStake            int64 = 1000000
	DefaultBaseRelaysPerPOKT   int64 = 100
	DefaultStabilityAdjustment int64 = 0
	DefaultParticipationRateOn bool  = false
	DefaultMaxChains           int64 = 15
)

// Keys for parameter access
var (
	KeyUnstakingTime       = []byte("AppUnstakingTime")
	KeyMaxApplications     = []byte("MaxApplications")
	KeyApplicationMinStake = []byte("ApplicationStakeMinimum")
	BaseRelaysPerPOKT      = []byte("BaseRelaysPerPOKT")
	StabilityAdjustment    = []byte("StabilityAdjustment")
	ParticipationRateOn    = []byte("ParticipationRateOn")
	KeyMaximumChains       = []byte("MaximumChains")
)

var _ types.ParamSet = (*Params)(nil)

// Params defines the high level settings for pos module
type Params struct {
	UnstakingTime       time.Duration `json:"unstaking_time" yaml:"unstaking_time"`               // duration of unstaking
	MaxApplications     int64         `json:"max_applications" yaml:"max_applications"`           // maximum number of applications
	AppStakeMin         int64         `json:"app_stake_minimum" yaml:"app_stake_minimum"`         // minimum amount needed to stake as an application
	BaseRelaysPerPOKT   int64         `json:"base_relays_per_pokt" yaml:"base_relays_per_pokt"`   // base relays per POKT coin staked
	StabilityAdjustment int64         `json:"stability_adjustment" yaml:"stability_adjustment"`   // the stability adjustment from the governance
	ParticipationRateOn bool          `json:"participation_rate_on" yaml:"participation_rate_on"` // the participation rate affects the amount minted based on staked ratio
	MaxChains           int64         `json:"maximum_chains" yaml:"maximum_chains"`               // the maximum number of chains an app can stake for
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{
		{Key: KeyUnstakingTime, Value: &p.UnstakingTime},
		{Key: KeyMaxApplications, Value: &p.MaxApplications},
		{Key: KeyApplicationMinStake, Value: &p.AppStakeMin},
		{Key: BaseRelaysPerPOKT, Value: &p.BaseRelaysPerPOKT},
		{Key: StabilityAdjustment, Value: &p.StabilityAdjustment},
		{Key: ParticipationRateOn, Value: &p.ParticipationRateOn},
		{Key: KeyMaximumChains, Value: &p.MaxChains},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		UnstakingTime:       DefaultUnstakingTime,
		MaxApplications:     DefaultMaxApplications,
		AppStakeMin:         DefaultMinStake,
		BaseRelaysPerPOKT:   DefaultBaseRelaysPerPOKT,
		StabilityAdjustment: DefaultStabilityAdjustment,
		ParticipationRateOn: DefaultParticipationRateOn,
		MaxChains:           DefaultMaxChains,
	}
}

// Validate a set of params
func (p Params) Validate() error {
	if p.MaxApplications == 0 {
		return fmt.Errorf("staking parameter MaxApplications must be a positive integer")
	}
	if p.AppStakeMin < DefaultMinStake {
		return fmt.Errorf("staking parameter StakeMimimum must be a positive integer")
	}
	if p.BaseRelaysPerPOKT < 0 {
		return fmt.Errorf("invalid baseline throughput stake rate, must be above 0")
	}
	// todo
	return nil
}

// Checks the equality of two param objects
func (p Params) Equal(p2 Params) bool {
	return reflect.DeepEqual(p, p2)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Unstaking Time:              %s
  Max Applications:            %d
  Minimum Stake:     	       %d
  BaseRelaysPerPOKT            %d
  Stability Adjustment         %d
  Participation Rate On        %v
  MaxChains                    %d,`,
		p.UnstakingTime,
		p.MaxApplications,
		p.AppStakeMin,
		p.BaseRelaysPerPOKT,
		p.StabilityAdjustment,
		p.ParticipationRateOn,
		p.MaxChains)
}
