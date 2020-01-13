package types

import (
	"fmt"
	"github.com/pokt-network/posmint/codec"
	"github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/params"
	"github.com/tendermint/go-amino"
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
				UnstakingTime:            DefaultUnstakingTime,
				MaxValidators:            DefaultMaxValidators,
				StakeMinimum:             DefaultMinStake,
				StakeDenom:               types.DefaultBondDenom,
				ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
				MaxEvidenceAge:           DefaultMaxEvidenceAge,
				SignedBlocksWindow:       DefaultSignedBlocksWindow,
				MinSignedPerWindow:       DefaultMinSignedPerWindow,
				DowntimeJailDuration:     DefaultDowntimeJailDuration,
				SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
				SlashFractionDowntime:    DefaultSlashFractionDowntime,
				SessionBlockFrequency:    DefaultSessionBlocktime,
				RelaysToTokens:           DefaultRelaysToTokens,
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
		MaxValidators            uint64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.Dec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.Dec
		SlashFractionDowntime    types.Dec
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
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, args{Params{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{}}}, true},
		{"Default Test False", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, args{Params{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			MaxEvidenceAge:           1,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if got := p.Equal(tt.args.p2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_Validate(t *testing.T) {
	type fields struct {
		UnstakingTime            time.Duration `json:"unstaking_time" yaml:"unstaking_time"`           // duration of unstaking
		MaxValidators            uint64        `json:"max_validators" yaml:"max_validators"`           // maximum number of validators
		StakeDenom               string        `json:"stake_denom" yaml:"stake_denom"`                 // bondable coin denomination
		StakeMinimum             int64         `json:"stake_minimum" yaml:"stake_minimum"`             // minimum amount needed to stake
		ProposerRewardPercentage int8          `json:"base_proposer_award" yaml:"base_proposer_award"` // minimum award for the proposer
		SessionBlock             int64         `json:"session_block" yaml:"session_block"`
		RelaysToTokens           types.Dec     `json:"relays_to_tokens" yaml:"relays_to_tokens"`
		// slashing params
		MaxEvidenceAge          time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`
		SignedBlocksWindow      int64         `json:"signed_blocks_window" yaml:"signed_blocks_window"`
		MinSignedPerWindow      types.Dec     `json:"min_signed_per_window" yaml:"min_signed_per_window"`
		DowntimeJailDuration    time.Duration `json:"downtime_jail_duration" yaml:"downtime_jail_duration"`
		SlashFractionDoubleSign types.Dec     `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"`
		SlashFractionDowntime   types.Dec     `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Default Validation Test / Wrong All Parameters", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			SessionBlock:             0,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, true},
		{"Default Validation Test / Wrong StakeDenom", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "",
			StakeMinimum:             2,
			ProposerRewardPercentage: 0,
			SessionBlock:             1,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong sessionblock", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "3",
			StakeMinimum:             2,
			ProposerRewardPercentage: 0,
			SessionBlock:             0,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong max val", fields{
			UnstakingTime:            0,
			MaxValidators:            0,
			StakeDenom:               "3",
			StakeMinimum:             2,
			ProposerRewardPercentage: 0,
			SessionBlock:             1,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong stake minimun", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "3",
			StakeMinimum:             0,
			ProposerRewardPercentage: 0,
			SessionBlock:             1,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong reward percentage above", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "3",
			StakeMinimum:             1,
			ProposerRewardPercentage: 101,
			SessionBlock:             1,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong reward percentage below", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "3",
			StakeMinimum:             1,
			ProposerRewardPercentage: -2,
			SessionBlock:             1,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Wrong relays to token", fields{
			UnstakingTime:            0,
			MaxValidators:            2,
			StakeDenom:               "3",
			StakeMinimum:             1,
			ProposerRewardPercentage: 100,
			SessionBlock:             1,
			RelaysToTokens:           types.NewDec(2),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.ZeroDec(),
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.ZeroDec(),
			SlashFractionDowntime:    types.ZeroDec(),
		}, true},
		{"Default Validation Test / Valid", fields{
			UnstakingTime:            0,
			MaxValidators:            1000,
			StakeDenom:               "3",
			StakeMinimum:             1,
			ProposerRewardPercentage: 100,
			SessionBlock:             30,
			RelaysToTokens:           types.OneDec(),
			MaxEvidenceAge:           0,
			SignedBlocksWindow:       0,
			MinSignedPerWindow:       types.Dec{},
			DowntimeJailDuration:     0,
			SlashFractionDoubleSign:  types.Dec{},
			SlashFractionDowntime:    types.Dec{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				SessionBlockFrequency:    tt.fields.SessionBlock,
				RelaysToTokens:           tt.fields.RelaysToTokens,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParams_ParamSetPairs(t *testing.T) {
	type fields struct {
		UnstakingTime            time.Duration
		MaxValidators            uint64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		SessionBlockFrequency    int64
		RelaysToTokens           types.Dec
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.Dec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.Dec
		SlashFractionDowntime    types.Dec
	}

	defParam := Params{
		UnstakingTime:            DefaultUnstakingTime,
		MaxValidators:            DefaultMaxValidators,
		StakeMinimum:             DefaultMinStake,
		StakeDenom:               types.DefaultBondDenom,
		ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
		MaxEvidenceAge:           DefaultMaxEvidenceAge,
		SignedBlocksWindow:       DefaultSignedBlocksWindow,
		MinSignedPerWindow:       DefaultMinSignedPerWindow,
		DowntimeJailDuration:     DefaultDowntimeJailDuration,
		SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
		SlashFractionDowntime:    DefaultSlashFractionDowntime,
		SessionBlockFrequency:    DefaultSessionBlocktime,
		RelaysToTokens:           DefaultRelaysToTokens,
	}
	defparamPairs := defParam.ParamSetPairs()

	tests := []struct {
		name   string
		fields fields
		want   params.ParamSetPairs
	}{
		{"Test Set Pairs", fields{
			UnstakingTime:            DefaultUnstakingTime,
			MaxValidators:            DefaultMaxValidators,
			StakeMinimum:             DefaultMinStake,
			StakeDenom:               types.DefaultBondDenom,
			ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
			MaxEvidenceAge:           DefaultMaxEvidenceAge,
			SignedBlocksWindow:       DefaultSignedBlocksWindow,
			MinSignedPerWindow:       DefaultMinSignedPerWindow,
			DowntimeJailDuration:     DefaultDowntimeJailDuration,
			SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
			SlashFractionDowntime:    DefaultSlashFractionDowntime,
			SessionBlockFrequency:    DefaultSessionBlocktime,
			RelaysToTokens:           DefaultRelaysToTokens,
		}, defparamPairs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				SessionBlockFrequency:    tt.fields.SessionBlockFrequency,
				RelaysToTokens:           tt.fields.RelaysToTokens,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if got := p.ParamSetPairs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParamSetPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParams_String(t *testing.T) {
	type fields struct {
		UnstakingTime            time.Duration
		MaxValidators            uint64
		StakeDenom               string
		StakeMinimum             int64
		ProposerRewardPercentage int8
		SessionBlockFrequency    int64
		RelaysToTokens           types.Dec
		MaxEvidenceAge           time.Duration
		SignedBlocksWindow       int64
		MinSignedPerWindow       types.Dec
		DowntimeJailDuration     time.Duration
		SlashFractionDoubleSign  types.Dec
		SlashFractionDowntime    types.Dec
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"String Test", fields{
			UnstakingTime:            DefaultUnstakingTime,
			MaxValidators:            DefaultMaxValidators,
			StakeMinimum:             DefaultMinStake,
			StakeDenom:               types.DefaultBondDenom,
			ProposerRewardPercentage: DefaultBaseProposerAwardPercentage,
			MaxEvidenceAge:           DefaultMaxEvidenceAge,
			SignedBlocksWindow:       DefaultSignedBlocksWindow,
			MinSignedPerWindow:       DefaultMinSignedPerWindow,
			DowntimeJailDuration:     DefaultDowntimeJailDuration,
			SlashFractionDoubleSign:  DefaultSlashFractionDoubleSign,
			SlashFractionDowntime:    DefaultSlashFractionDowntime,
			SessionBlockFrequency:    DefaultSessionBlocktime,
			RelaysToTokens:           DefaultRelaysToTokens,
		}, fmt.Sprintf(`Params:
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
  SlashFractionDowntime:   %s
  SessionBlockFrequency    %d`,
			DefaultUnstakingTime,
			DefaultMaxValidators,
			types.DefaultBondDenom,
			DefaultMinStake,
			DefaultBaseProposerAwardPercentage,
			DefaultMaxEvidenceAge,
			DefaultSignedBlocksWindow,
			DefaultMinSignedPerWindow,
			DefaultDowntimeJailDuration,
			DefaultSlashFractionDoubleSign,
			DefaultSlashFractionDowntime,
			DefaultSessionBlocktime)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Params{
				UnstakingTime:            tt.fields.UnstakingTime,
				MaxValidators:            tt.fields.MaxValidators,
				StakeDenom:               tt.fields.StakeDenom,
				StakeMinimum:             tt.fields.StakeMinimum,
				ProposerRewardPercentage: tt.fields.ProposerRewardPercentage,
				SessionBlockFrequency:    tt.fields.SessionBlockFrequency,
				RelaysToTokens:           tt.fields.RelaysToTokens,
				MaxEvidenceAge:           tt.fields.MaxEvidenceAge,
				SignedBlocksWindow:       tt.fields.SignedBlocksWindow,
				MinSignedPerWindow:       tt.fields.MinSignedPerWindow,
				DowntimeJailDuration:     tt.fields.DowntimeJailDuration,
				SlashFractionDoubleSign:  tt.fields.SlashFractionDoubleSign,
				SlashFractionDowntime:    tt.fields.SlashFractionDowntime,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalParams(t *testing.T) {
	type args struct {
		cdc   *codec.Codec
		value []byte
	}

	defaultParams := DefaultParams()
	value, _ := amino.MarshalBinaryLengthPrefixed(DefaultParams())

	tests := []struct {
		name       string
		args       args
		wantParams Params
		wantErr    bool
	}{
		{"Unmarshall Test", args{
			cdc:   codec.New(),
			value: value,
		}, defaultParams, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, err := UnmarshalParams(tt.args.cdc, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("UnmarshalParams() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}

func TestMustUnmarshalParams(t *testing.T) {
	type args struct {
		cdc   *codec.Codec
		value []byte
	}

	defaultParams := DefaultParams()
	value, _ := amino.MarshalBinaryLengthPrefixed(DefaultParams())

	tests := []struct {
		name string
		args args
		want Params
	}{
		{"Must Unmarshall Test", args{
			cdc:   codec.New(),
			value: value,
		}, defaultParams},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustUnmarshalParams(tt.args.cdc, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustUnmarshalParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
