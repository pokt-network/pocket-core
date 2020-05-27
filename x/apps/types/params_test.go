package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestDefaultParams(t *testing.T) {
	tests := []struct {
		name string
		want Params
	}{
		{"Default Test",
			Params{
				UnstakingTime:       DefaultUnstakingTime,
				MaxApplications:     DefaultMaxApplications,
				AppStakeMin:         DefaultMinStake,
				BaseRelaysPerPOKT:   DefaultBaseRelaysPerPOKT,
				StabilityAdjustment: DefaultStabilityAdjustment,
				ParticipationRateOn: DefaultParticipationRateOn,
				MaxChains:           DefaultMaxChains,
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Equal(t *testing.T) {
	type fields struct {
		UnstakingTime              time.Duration `json:"unstaking_time" yaml:"unstaking_time"`       // duration of unstaking
		MaxApplications            int64         `json:"max_applications" yaml:"max_applications"`   // maximum number of applications
		AppStakeMin                int64         `json:"app_stake_minimum" yaml:"app_stake_minimum"` // minimum amount needed to stake
		BaslineThroughputStakeRate int64
		StabilityAdjustment        int64
		ParticipationRateOn        bool
	}
	type args struct {
		p2 Params
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Default Test Equal", fields{
			UnstakingTime:              0,
			MaxApplications:            0,
			AppStakeMin:                0,
			BaslineThroughputStakeRate: 0,
			StabilityAdjustment:        0,
			ParticipationRateOn:        false,
		}, args{Params{
			UnstakingTime:       0,
			MaxApplications:     0,
			AppStakeMin:         0,
			BaseRelaysPerPOKT:   0,
			StabilityAdjustment: 0,
			ParticipationRateOn: false,
		}}, true},
		{"Default Test False", fields{
			UnstakingTime:              0,
			MaxApplications:            0,
			AppStakeMin:                0,
			BaslineThroughputStakeRate: 0,
			StabilityAdjustment:        0,
			ParticipationRateOn:        false,
		}, args{Params{
			UnstakingTime:       0,
			MaxApplications:     1,
			AppStakeMin:         0,
			BaseRelaysPerPOKT:   0,
			StabilityAdjustment: 0,
			ParticipationRateOn: false,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:       tt.fields.UnstakingTime,
				MaxApplications:     tt.fields.MaxApplications,
				AppStakeMin:         tt.fields.AppStakeMin,
				BaseRelaysPerPOKT:   tt.fields.BaslineThroughputStakeRate,
				StabilityAdjustment: tt.fields.StabilityAdjustment,
				ParticipationRateOn: tt.fields.ParticipationRateOn,
			}
			if got := p.Equal(tt.args.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	type fields struct {
		UnstakingTime               time.Duration `json:"unstaking_time" yaml:"unstaking_time"`       // duration of unstaking
		MaxApplications             int64         `json:"max_applications" yaml:"max_applications"`   // maximum number of applications
		AppStakeMin                 int64         `json:"app_stake_minimum" yaml:"app_stake_minimum"` // minimum amount needed to stake
		BaselineThrouhgputStakeRate int64         `json:"baseline_throughput_stake_rate" yaml:"baseline_throughput_stake_rate"`
		StabilityAdjustment         int64         `json:"staking_adjustment" yaml:"staking_adjustment"`
		ParticipationRateOn         bool          `json:"participation_rate_on" yaml:"participation_rate_on"`
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Default Validation Test / Wrong All Parameters", fields{
			UnstakingTime:               0,
			MaxApplications:             0,
			AppStakeMin:                 0,
			BaselineThrouhgputStakeRate: 1,
			StabilityAdjustment:         0,
			ParticipationRateOn:         false,
		}, true},
		{"Default Validation Test / Wrong Appstake", fields{
			UnstakingTime:               0,
			MaxApplications:             2,
			AppStakeMin:                 0,
			BaselineThrouhgputStakeRate: 0,
		}, true},
		{"Default Validation Test / Wrong BaselinethroughputStakeRate", fields{
			UnstakingTime:               10000,
			MaxApplications:             2,
			AppStakeMin:                 1000000,
			BaselineThrouhgputStakeRate: -1,
		}, true},
		{"Default Validation Test / Valid", fields{
			UnstakingTime:               10000,
			MaxApplications:             2,
			AppStakeMin:                 1000000,
			BaselineThrouhgputStakeRate: 90,
			StabilityAdjustment:         100,
			ParticipationRateOn:         false,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:       tt.fields.UnstakingTime,
				MaxApplications:     tt.fields.MaxApplications,
				AppStakeMin:         tt.fields.AppStakeMin,
				BaseRelaysPerPOKT:   tt.fields.BaselineThrouhgputStakeRate,
				StabilityAdjustment: tt.fields.StabilityAdjustment,
				ParticipationRateOn: tt.fields.ParticipationRateOn,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Default Test",
			fmt.Sprintf(`Params:
  Unstaking Time:              %s
  Max Applications:            %d
  Minimum Stake:     	       %d
  BaseRelaysPerPOKT            %d
  Stability Adjustment         %d
  Participation Rate On        %v
  Maximum Chains               %d,`,
				DefaultUnstakingTime,
				DefaultMaxApplications,
				DefaultMinStake,
				DefaultBaseRelaysPerPOKT,
				DefaultStabilityAdjustment,
				DefaultParticipationRateOn,
				DefaultMaxChains),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
