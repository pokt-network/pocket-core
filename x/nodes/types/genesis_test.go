package types

import (
	"github.com/pokt-network/posmint/types"
	"reflect"
	"testing"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want GenesisState
	}{{"defaultState", GenesisState{
		Params:       DefaultParams(),
		SigningInfos: make(map[string]ValidatorSigningInfo),
		MissedBlocks: make(map[string][]MissedBlock),
		DAO:          DAOPool(NewPool(types.ZeroInt())),
	}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultGenesisState(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGenesisState(t *testing.T) {
	type args struct {
		params           Params
		validators       []Validator
		dao              DAOPool
		previousProposer types.Address
		signingInfos     map[string]ValidatorSigningInfo
		missedBlocks     map[string][]MissedBlock
	}
	ca, _ := types.AddressFromHex("ABCDEFGHIJK")

	tests := []struct {
		name string
		args args
		want GenesisState
	}{
		{"Default Change State Test", args{
			params:           DefaultParams(),
			validators:       nil,
			dao:              DAOPool(NewPool(types.ZeroInt())),
			previousProposer: ca,
			signingInfos:     make(map[string]ValidatorSigningInfo),
			missedBlocks:     make(map[string][]MissedBlock)},
			GenesisState{
				Params:           DefaultParams(),
				Validators:       nil,
				DAO:              DAOPool(NewPool(types.ZeroInt())),
				SigningInfos:     make(map[string]ValidatorSigningInfo),
				MissedBlocks:     make(map[string][]MissedBlock),
				PreviousProposer: ca,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenesisState(tt.args.params, tt.args.validators, tt.args.dao, tt.args.previousProposer, tt.args.signingInfos, tt.args.missedBlocks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}
