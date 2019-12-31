package types

import (
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
				UnstakingTime:              DefaultUnstakingTime,
				MaxApplications:            DefaultMaxApplications,
				AppStakeMin:                DefaultMinStake,
				RelayCoefficientPercentage: DefaultRelayCoefficient,
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
		MaxApplications            uint64        `json:"max_applications" yaml:"max_applications"`   // maximum number of applications
		AppStakeMin                int64         `json:"app_stake_minimum" yaml:"app_stake_minimum"` // minimum amount needed to stake
		RelayCoefficientPercentage uint8
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
			RelayCoefficientPercentage: 0,
		}, args{Params{
			UnstakingTime:              0,
			MaxApplications:            0,
			AppStakeMin:                0,
			RelayCoefficientPercentage: 0,
		}}, true},
		{"Default Test False", fields{
			UnstakingTime:              0,
			MaxApplications:            0,
			AppStakeMin:                0,
			RelayCoefficientPercentage: 0,
		}, args{Params{
			UnstakingTime:              0,
			MaxApplications:            1,
			AppStakeMin:                0,
			RelayCoefficientPercentage: 0,
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:              tt.fields.UnstakingTime,
				MaxApplications:            tt.fields.MaxApplications,
				AppStakeMin:                tt.fields.AppStakeMin,
				RelayCoefficientPercentage: tt.fields.RelayCoefficientPercentage,
			}
			if got := p.Equal(tt.args.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	type fields struct {
		UnstakingTime              time.Duration `json:"unstaking_time" yaml:"unstaking_time"`       // duration of unstaking
		MaxApplications            uint64        `json:"max_applications" yaml:"max_applications"`   // maximum number of applications
		AppStakeMin                int64         `json:"app_stake_minimum" yaml:"app_stake_minimum"` // minimum amount needed to stake
		RelayCoefficientPercentage uint8
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Default Validation Test / Wrong All Parameters", fields{
			UnstakingTime:              0,
			MaxApplications:            0,
			AppStakeMin:                0,
			RelayCoefficientPercentage: 0,
		}, true},
		{"Default Validation Test / Valid", fields{
			UnstakingTime:              10000,
			MaxApplications:            2,
			AppStakeMin:                1,
			RelayCoefficientPercentage: 90,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:              tt.fields.UnstakingTime,
				MaxApplications:            tt.fields.MaxApplications,
				AppStakeMin:                tt.fields.AppStakeMin,
				RelayCoefficientPercentage: tt.fields.RelayCoefficientPercentage,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
