package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/types"
)

func TestDefaultParams(t *testing.T) {
	tests := []struct {
		name string
		want Params
	}{
		{"Default Test",
			Params{
				UnstakingTime:            DefaultUnstakingTime,
				MaxValidators:            DefaultMaxValidators,
				StakeMinimum:             DefaultMinStake,
				StakeDenom:               types.DefaultStakeDenom,
				MaxEvidenceAge:           DefaultMaxEvidenceAge,
				SignedBlocksWindow:       DefaultSignedBlocksWindow,
				MinSignedPerWindow:       DefaultMinSignedPerWindow,
				DowntimeJailDuration:     DefaultDowntimeJailDuration,
				SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
				SlashFractionDowntime:    DefaultSlashFractionDowntime,
				SessionBlockFrequency:    DefaultSessionBlocktime,
				DAOAllocation:            DefaultDAOAllocation,
				ProposerAllocation:       DefaultProposerAllocation,
				RelaysToTokensMultiplier: DefaultRelaysToTokensMultiplier,
				MaximumChains:            DefaultMaxChains,
				MaxJailedBlocks:          DefaultMaxJailedBlocks,
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
		UnstakingTime            time.Duration
		MaxValidators            int64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.BigDec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.BigDec
		SlashFractionDowntime    types.BigDec
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
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "",
			StakeMinimum:            0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.BigDec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.BigDec{},
			SlashFractionDowntime:   types.BigDec{},
		}, args{Params{
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "",
			StakeMinimum:            0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.BigDec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.BigDec{},
			SlashFractionDowntime:   types.BigDec{}}}, true},
		{"Default Test False", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.BigDec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.BigDec{},
			SlashFractionDowntime:    types.BigDec{},
		}, args{Params{
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "",
			StakeMinimum:            0,
			MaxEvidenceAge:          1,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.BigDec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.BigDec{},
			SlashFractionDowntime:   types.BigDec{}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:           tt.fields.UnstakingTime,
				MaxValidators:           tt.fields.MaxValidators,
				StakeDenom:              tt.fields.StakeDenom,
				StakeMinimum:            tt.fields.StakeMinimum,
				MaxEvidenceAge:          tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:      tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:      tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:    tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign: tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:   tt.fields.SlashFractionDowntime,
			}
			if got := p.Equal(tt.args.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	type fields struct {
		UnstakingTime      time.Duration `json:"unstaking_time" yaml:"unstaking_time"` // duration of unstaking
		MaxValidators      int64         `json:"max_validators" yaml:"max_validators"` // maximum number of validators
		StakeDenom         string        `json:"stake_denom" yaml:"stake_denom"`       // bondable coin denomination
		StakeMinimum       int64         `json:"stake_minimum" yaml:"stake_minimum"`   // minimum amount needed to stake
		DaoAllocation      int64
		SessionBlock       int64 `json:"session_block" yaml:"session_block"`
		ProposerAllocation int64 `json:"relays_to_tokens" yaml:"relays_to_tokens"`
		// slashing params
		MaxEvidenceAge          time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`
		SignedBlocksWindow      int64         `json:"signed_blocks_window" yaml:"signed_blocks_window"`
		MinSignedPerWindow      types.BigDec  `json:"min_signed_per_window" yaml:"min_signed_per_window"`
		DowntimeJailDuration    time.Duration `json:"downtime_jail_duration" yaml:"downtime_jail_duration"`
		SlashFractionDoubleSign types.BigDec  `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"`
		SlashFractionDowntime   types.BigDec  `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Default Validation Test / Wrong All Parameters", fields{
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "",
			StakeMinimum:            0,
			SessionBlock:            0,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.BigDec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.BigDec{},
			SlashFractionDowntime:   types.BigDec{},
		}, true},
		{"Default Validation Test / Wrong StakeDenom", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "",
			StakeMinimum:            1000000,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong sessionblock", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			SessionBlock:            0,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong max val", fields{
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong stake minimun", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "3",
			StakeMinimum:            1,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong reward percentage above", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong dao allocation below", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			DaoAllocation:           -2,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong relays to token", fields{
			UnstakingTime:           0,
			MaxValidators:           2,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			SessionBlock:            1,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.ZeroDec(),
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.ZeroDec(),
			SlashFractionDowntime:   types.ZeroDec(),
		}, true},
		{"Default Validation Test / Valid", fields{
			UnstakingTime:           0,
			MaxValidators:           1000,
			StakeDenom:              "3",
			StakeMinimum:            1000000,
			SessionBlock:            30,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      types.BigDec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: types.BigDec{},
			SlashFractionDowntime:   types.BigDec{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:           tt.fields.UnstakingTime,
				MaxValidators:           tt.fields.MaxValidators,
				StakeDenom:              tt.fields.StakeDenom,
				StakeMinimum:            tt.fields.StakeMinimum,
				DAOAllocation:           tt.fields.DaoAllocation,
				SessionBlockFrequency:   tt.fields.SessionBlock,
				ProposerAllocation:      tt.fields.ProposerAllocation,
				MaxEvidenceAge:          tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:      tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:      tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:    tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign: tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:   tt.fields.SlashFractionDowntime,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_ParamSetPairs(t *testing.T) {
	type fields struct {
		UnstakingTime           time.Duration
		MaxValidators           int64
		StakeDenom              string
		StakeMinimum            int64
		ProposerAllocation      int64
		DAOAllocation           int64
		SessionBlockFrequency   int64
		MaxEvidenceAge          time.Duration
		SignedBlocksWindow      int64
		MinSignedPerWindow      types.BigDec
		DowntimeJailDuration    time.Duration
		SlashFractionDoubleSign types.BigDec
		SlashFractionDowntime   types.BigDec
		MaximumChains           int64
		MaxJailedBlocks         int64
	}

	defParam := Params{
		UnstakingTime:           DefaultUnstakingTime,
		MaxValidators:           DefaultMaxValidators,
		StakeMinimum:            DefaultMinStake,
		StakeDenom:              types.DefaultStakeDenom,
		ProposerAllocation:      DefaultProposerAllocation,
		MaxEvidenceAge:          DefaultMaxEvidenceAge,
		SignedBlocksWindow:      DefaultSignedBlocksWindow,
		MinSignedPerWindow:      DefaultMinSignedPerWindow,
		DowntimeJailDuration:    DefaultDowntimeJailDuration,
		SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
		SlashFractionDowntime:   DefaultSlashFractionDowntime,
		SessionBlockFrequency:   DefaultSessionBlocktime,
		DAOAllocation:           DefaultDAOAllocation,
		MaximumChains:           DefaultMaxChains,
		MaxJailedBlocks:         DefaultMaxJailedBlocks,
	}
	defparamPairs := defParam.ParamSetPairs()

	tests := []struct {
		name   string
		fields fields
		want   types.ParamSetPairs
	}{
		{"Test Set Pairs", fields{
			UnstakingTime:           DefaultUnstakingTime,
			MaxValidators:           DefaultMaxValidators,
			StakeMinimum:            DefaultMinStake,
			StakeDenom:              types.DefaultStakeDenom,
			ProposerAllocation:      DefaultProposerAllocation,
			MaxEvidenceAge:          DefaultMaxEvidenceAge,
			SignedBlocksWindow:      DefaultSignedBlocksWindow,
			MinSignedPerWindow:      DefaultMinSignedPerWindow,
			DowntimeJailDuration:    DefaultDowntimeJailDuration,
			SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
			SlashFractionDowntime:   DefaultSlashFractionDowntime,
			SessionBlockFrequency:   DefaultSessionBlocktime,
			DAOAllocation:           DefaultDAOAllocation,
			MaximumChains:           DefaultMaxChains,
			MaxJailedBlocks:         DefaultMaxJailedBlocks,
		}, defparamPairs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				UnstakingTime:           tt.fields.UnstakingTime,
				MaxValidators:           tt.fields.MaxValidators,
				StakeDenom:              tt.fields.StakeDenom,
				StakeMinimum:            tt.fields.StakeMinimum,
				ProposerAllocation:      tt.fields.ProposerAllocation,
				SessionBlockFrequency:   tt.fields.SessionBlockFrequency,
				DAOAllocation:           tt.fields.DAOAllocation,
				MaxEvidenceAge:          tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:      tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:      tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:    tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign: tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:   tt.fields.SlashFractionDowntime,
				MaximumChains:           tt.fields.MaximumChains,
				MaxJailedBlocks:         tt.fields.MaxJailedBlocks,
			}
			if got := p.ParamSetPairs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParamSetPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_String(t *testing.T) {
	type fields struct {
		UnstakingTime           time.Duration
		MaxValidators           int64
		StakeDenom              string
		StakeMinimum            int64
		ProposerAllocation      int64
		SessionBlockFrequency   int64
		DaoAllocation           int64
		MaxEvidenceAge          time.Duration
		SignedBlocksWindow      int64
		MinSignedPerWindow      types.BigDec
		DowntimeJailDuration    time.Duration
		SlashFractionDoubleSign types.BigDec
		SlashFractionDowntime   types.BigDec
		MaximumChains           int64
		MaxJailedBlocks         int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"String Test", fields{
			UnstakingTime:           DefaultUnstakingTime,
			MaxValidators:           DefaultMaxValidators,
			StakeMinimum:            DefaultMinStake,
			StakeDenom:              types.DefaultStakeDenom,
			ProposerAllocation:      DefaultProposerAllocation,
			MaxEvidenceAge:          DefaultMaxEvidenceAge,
			SignedBlocksWindow:      DefaultSignedBlocksWindow,
			MinSignedPerWindow:      DefaultMinSignedPerWindow,
			DowntimeJailDuration:    DefaultDowntimeJailDuration,
			SlashFractionDoubleSign: DefaultSlashFractionDoubleSign,
			SlashFractionDowntime:   DefaultSlashFractionDowntime,
			SessionBlockFrequency:   DefaultSessionBlocktime,
			DaoAllocation:           DefaultDAOAllocation,
			MaximumChains:           DefaultMaxChains,
			MaxJailedBlocks:         DefaultMaxJailedBlocks,
		}, fmt.Sprintf(`Params:
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
  BlocksPerSession         %d
  Proposer Allocation      %d
  DAO allocation           %d
  Maximum Chains           %d
  Max Jailed Blocks        %d`,
			DefaultUnstakingTime,
			DefaultMaxValidators,
			types.DefaultStakeDenom,
			DefaultMinStake,
			DefaultMaxEvidenceAge,
			DefaultSignedBlocksWindow,
			DefaultMinSignedPerWindow,
			DefaultDowntimeJailDuration,
			DefaultSlashFractionDoubleSign,
			DefaultSlashFractionDowntime,
			DefaultSessionBlocktime,
			DefaultProposerAllocation,
			DefaultDAOAllocation,
			DefaultMaxChains,
			DefaultMaxJailedBlocks)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:           tt.fields.UnstakingTime,
				MaxValidators:           tt.fields.MaxValidators,
				StakeDenom:              tt.fields.StakeDenom,
				StakeMinimum:            tt.fields.StakeMinimum,
				ProposerAllocation:      tt.fields.ProposerAllocation,
				SessionBlockFrequency:   tt.fields.SessionBlockFrequency,
				DAOAllocation:           tt.fields.DaoAllocation,
				MaxEvidenceAge:          tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:      tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:      tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:    tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign: tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:   tt.fields.SlashFractionDowntime,
				MaximumChains:           tt.fields.MaximumChains,
				MaxJailedBlocks:         tt.fields.MaxJailedBlocks,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
