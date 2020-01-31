package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestInvariants_RegisterInvariants(t *testing.T) {
	tests := []struct {
		name      string
		invariant *InvariantRegistry
	}{
		{
			name:      "register invariant",
			invariant: new(InvariantRegistry),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, keeper := createTestInput(t, true)

			tt.invariant.On("RegisterRoute", types.ModuleName, "module-accounts", mock.AnythingOfType("types.Invariant")).Return(mock.Anything)
			tt.invariant.On("RegisterRoute", types.ModuleName, "nonnegative-power", mock.AnythingOfType("types.Invariant")).Return(mock.Anything)
			RegisterInvariants(tt.invariant, keeper)
		})
	}
}
func TestInvariants_ModuleAccountInvariants(t *testing.T) {
	context, _, keeper := createTestInput(t, true)

	stakedPool := keeper.GetStakedTokens(context)
	staked := sdk.ZeroInt()
	notStaked := sdk.ZeroInt()
	notStakedPool := keeper.GetUnstakedTokens(context)

	invariantMsg := sdk.FormatInvariant(types.ModuleName, "staked and not staked module account coins", fmt.Sprintf(
		"\tPool's staked tokens: %v\n"+
			"\tsum of staked tokens: %v\n"+
			"not staked token invariance:\n"+
			"\tPool's not staked tokens: %v\n"+
			"\tsum of not staked tokens: %v\n"+
			"module accounts total (staked + not staked):\n"+
			"\tModule Accounts' tokens: %v\n"+
			"\tsum tokens:              %v\n",
		stakedPool, staked, notStakedPool, notStaked, stakedPool.Add(notStakedPool), staked.Add(notStaked)))

	tests := []struct {
		name string
		want string
	}{
		{
			name: "register invariant",
			want: invariantMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invariant := ModuleAccountInvariants(keeper)
			if got, _ := invariant(context); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleAccountInvariants() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvariants_NonNegativePowerInvariant(t *testing.T) {
	context, _, keeper := createTestInput(t, true)
	application := getStakedApplication()

	keeper.SetApplication(context, application)
	keeper.SetStakedApplication(context, application)

	var msg string
	invariantMsg := sdk.FormatInvariant(types.ModuleName, "nonnegative power", fmt.Sprintf("found invalid application powers\n%s", msg))

	tests := []struct {
		name string
		want string
	}{
		{
			name: "register invariant",
			want: invariantMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invariant := NonNegativePowerInvariant(keeper)
			if got, _ := invariant(context); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModuleAccountInvariants() = got %v, want %v", got, tt.want)
			}
		})
	}
}
