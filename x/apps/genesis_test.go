package pos

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"reflect"
	"testing"
)

func TestPos_InitGenesis(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "set init genesis"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, keeper, supplyKeeper, posKeeper := createTestInput(t, true)
			state := types.DefaultGenesisState()
			InitGenesis(context, keeper, supplyKeeper, posKeeper, state)
			if got := keeper.GetParams(context); got != state.Params {
				t.Errorf("InitGenesis()= got %v, want %v", got, state.Params)
			}
		})
	}
}
func TestPos_ExportGenesis(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "get genesis from app"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, keeper, supplyKeeper, posKeeper := createTestInput(t, true)
			state := types.DefaultGenesisState()
			InitGenesis(context, keeper, supplyKeeper, posKeeper, state)
			state.Exported = true // Export genesis returns an exported state
			if got := ExportGenesis(context, keeper); !reflect.DeepEqual(got, state) {
				t.Errorf("\nExportGenesis()=\nGot-> %v\nWant-> %v", got, state.Params)
			}
		})
	}
}
func TestPos_ValidateGeneis(t *testing.T) {
	application := getApplication()

	jailedApp := getApplication()
	jailedApp.Jailed = true

	zeroStakeApp := getApplication()
	zeroStakeApp.StakedTokens = sdk.NewInt(0)

	singleTokenApp := getApplication()
	singleTokenApp.StakedTokens = sdk.NewInt(1)
	tests := []struct {
		name   string
		state  types.GenesisState
		apps   types.Applications
		params bool
		want   interface{}
	}{
		{
			name: "valdiates genesis for application",
			apps: types.Applications{application},
			want: nil,
		},
		{
			name:   "errs if invalid params",
			apps:   types.Applications{application},
			params: true,
			want:   fmt.Errorf("staking parameter StakeMimimum must be a positive integer"),
		},
		{
			name: "errs if dupplicate application in geneiss state",
			apps: types.Applications{application, application},
			want: fmt.Errorf("duplicate application in genesis state: address %v", application.GetAddress()),
		},
		{
			name: "errs if jailed app staked",
			apps: types.Applications{jailedApp},
			want: fmt.Errorf("application is staked and jailed in genesis state: address %v", jailedApp.GetAddress()),
		},
		{
			name: "errs if staked with zero tokens",
			apps: types.Applications{zeroStakeApp},
			want: fmt.Errorf("staked/unstaked genesis application cannot have zero stake, application: %v", zeroStakeApp),
		},
		{
			name: "errs if lower or equal than minimum stake ",
			apps: types.Applications{singleTokenApp},
			want: fmt.Errorf("application has less than minimum stake: %v", singleTokenApp),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := types.DefaultGenesisState()
			state.Applications = tt.apps
			if tt.params {
				state.Params.AppStakeMin = 0
			}
			if got := ValidateGenesis(state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateGenesis()= got %v, want %v", got, tt.want)
			}
		})
	}
}
