package types

import (
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/types"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestDefaultGenesisState(t *testing.T) {
	tests := []struct {
		name string
		want GenesisState
	}{{"defaultState", GenesisState{
		Params:       DefaultParams(),
		Applications: make(Applications, 0),
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
		params       Params
		applications []Application
	}

	var pub crypto.Ed25519PublicKey
	rand.Read(pub[:])
	defaultApp := Application{
		Address:                 types.Address(pub.Address()),
		PublicKey:               pub,
		Jailed:                  false,
		Status:                  0,
		Chains:                  []string{"b60d7bdd334cd3768d43f14a05c7fe7e886ba5bcb77e1064530052fed1a3f145"},
		StakedTokens:            types.NewInt(int64(10000000)),
		MaxRelays:               types.NewInt(int64(10000000)),
		UnstakingCompletionTime: time.Unix(0, 0).UTC()}

	tests := []struct {
		name string
		args args
		want GenesisState
	}{
		{"Default Change State Test", args{
			params:       DefaultParams(),
			applications: []Application{defaultApp},
		},
			GenesisState{
				Params:       DefaultParams(),
				Applications: []Application{defaultApp},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenesisState(tt.args.params, tt.args.applications); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenesisState() = %v, want %v", got, tt.want)
			}
		})
	}
}
