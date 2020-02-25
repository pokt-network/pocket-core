package types

import (
	"reflect"
	"testing"
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
