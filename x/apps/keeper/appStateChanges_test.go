package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"reflect"
	"strings"
	"testing"
)

func TestAppStateChange_ValidateApplicaitonBeginUnstaking(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		want        interface{}
	}{
		{
			name:        "validates application",
			application: getStakedApplication(),
			want:        nil,
		},
		{
			name:        "errors if application not staked",
			application: getUnstakedApplication(),
			want:        types.ErrApplicationStatus("apps"),
		},
		{
			name:        "validates application",
			application: getStakedApplication(),
			panics:      true,
			want:        "should not happen: application trying to begin unstaking has less than the minimum stake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.panics {
			case true:
				defer func() {
					if err := recover(); err.(string) != tt.want {
						t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", err, tt.want)
					}
				}()
				tt.application.StakedTokens = sdk.NewInt(-1)
				_ = keeper.ValidateApplicationBeginUnstaking(context, tt.application)
			default:
				if got := keeper.ValidateApplicationBeginUnstaking(context, tt.application); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAppStateChange_ValidateApplicaitonStaking(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		amount      sdk.Int
		want        interface{}
	}{
		{
			name:        "validates application",
			application: getUnstakedApplication(),
			amount:      sdk.NewInt(1000000),
			want:        nil,
		},
		{
			name:        "errors if application staked",
			application: getStakedApplication(),
			amount:      sdk.NewInt(100),
			want:        types.ErrApplicationStatus("apps"),
		},
		{
			name:        "errors if application staked",
			application: getUnstakedApplication(),
			amount:      sdk.NewInt(0),
			want:        types.ErrMinimumStake("apps"),
		},
		{
			name:        "errors bank does not have enough coins",
			application: getUnstakedApplication(),
			amount:      sdk.NewInt(1000000000000000000),
			want:        types.ErrNotEnoughCoins("apps"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			if got := keeper.ValidateApplicationStaking(context, tt.application, tt.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppStateChange_JailApplication(t *testing.T) {
	jailedApp := getStakedApplication()
	jailedApp.Jailed = true
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		want        interface{}
	}{
		{
			name:        "jails application",
			application: getStakedApplication(),
			want:        true,
		},
		{
			name:        "already jailed app ",
			application: jailedApp,
			panics:      true,
			want:        fmt.Sprint("cannot jail already jailed application, application:"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.application)
			keeper.SetStakedApplication(context, tt.application)

			switch tt.panics {
			case true:
				defer func() {
					if err := recover().(string); !strings.Contains(err, tt.want.(string)) {
						t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, does not contain %v", err, tt.want)
					}
				}()
				keeper.JailApplication(context, tt.application.GetAddress())
			default:
				keeper.JailApplication(context, tt.application.GetAddress())
				if got, _ := keeper.GetApplication(context, tt.application.GetAddress()); got.Jailed != tt.want {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", tt.application.Jailed, tt.want)
				}
			}

		})
	}
}

func TestAppStateChange_UnjailApplication(t *testing.T) {
	jailedApp := getStakedApplication()
	jailedApp.Jailed = true
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		want        interface{}
	}{
		{
			name:        "unjails application",
			application: jailedApp,
			want:        false,
		},
		{
			name:        "already jailed app ",
			application: getStakedApplication(),
			panics:      true,
			want:        fmt.Sprint("cannot unjail already unjailed application, application:"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.application)
			keeper.SetStakedApplication(context, tt.application)

			switch tt.panics {
			case true:
				defer func() {
					if err := recover().(string); !strings.Contains(err, tt.want.(string)) {
						t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, does not contain %v", err, tt.want)
					}
				}()
				keeper.UnjailApplication(context, tt.application.GetAddress())
			default:
				keeper.UnjailApplication(context, tt.application.GetAddress())
				if got, _ := keeper.GetApplication(context, tt.application.GetAddress()); got.Jailed != tt.want {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", tt.application.Jailed, tt.want)
				}
			}

		})
	}
}

func TestAppStateChange_RegisterApplication(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
	}{
		{
			name:        "name registers apps",
			application: getStakedApplication(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.RegisterApplication(context, tt.application)

			_, found := keeper.GetApplication(context, tt.application.Address)
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app")
			}

			_, found = keeper.GetApplication(context, tt.application.GetAddress())
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app by ConsAddr")
			}
		})

	}
}

func TestAppStateChange_StakeApplication(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		amount      sdk.Int
	}{
		{
			name:        "name registers apps",
			application: getUnstakedApplication(),
			amount:      sdk.NewInt(100000000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			keeper.StakeApplication(context, tt.application, tt.amount)
			got, found := keeper.GetApplication(context, tt.application.Address)
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app")
			}
			if !got.StakedTokens.Equal(tt.amount) {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app %v", got.StakedTokens)
			}

		})

	}
}

func TestAppStateChange_BeginUnstakingApplication(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		want        sdk.StakeStatus
	}{
		{
			name:        "name registers apps",
			application: getStakedApplication(),
			want:        sdk.Unstaking,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			keeper.BeginUnstakingApplication(context, tt.application)
			got, found := keeper.GetApplication(context, tt.application.Address)
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app")
			}
			if got.Status != tt.want {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app %v", got.StakedTokens)
			}
		})
	}
}
