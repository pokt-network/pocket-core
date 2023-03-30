package keeper

import (
	"reflect"
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	coreTypes "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

func TestApplication_SetAndGetApplication(t *testing.T) {
	application := getStakedApplication()

	tests := []struct {
		name        string
		application types.Application
		want        bool
	}{
		{
			name:        "get and set application",
			application: application,
			want:        true,
		},
		{
			name:        "not found",
			application: application,
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			if tt.want {
				keeper.SetApplication(context, tt.application)
			}

			if _, found := keeper.GetApplication(context, tt.application.Address); found != tt.want {
				t.Errorf("Applicaiton.GetApplication() = got %v, want %v", found, tt.want)
			}
		})
	}
}

func TestApplication_CalculateAppRelays(t *testing.T) {
	tests := []struct {
		testName            string
		appStake            sdk.BigInt // uPOKT
		appChains           []string
		stabilityAdjustment int64
		baseRelaysPerPOKT   int64
		participationRateOn bool
		sessionNodeCount    int64
		wantAppRelays       sdk.BigInt // 1 * (baseRelaysPerPOKT / 100) * (appStake / 1000000) + 0
		wantSessionRelays   sdk.BigInt // wantAppRelays / numAppChains / sessionNodeCount
	}{
		{
			testName:          "Calculate App relays - default",
			appStake:          getStakedApplication().StakedTokens, // default
			sessionNodeCount:  1,
			wantAppRelays:     sdk.NewInt(100000),
			wantSessionRelays: sdk.NewInt(100000),
		},
		{
			testName:            "Calculate App relays - param values at height=90074",
			appStake:            sdk.NewInt(2228350000000),
			appChains:           []string{"0021"},
			stabilityAdjustment: 0,
			baseRelaysPerPOKT:   200000,
			participationRateOn: false,
			sessionNodeCount:    24,
			wantAppRelays:       sdk.NewInt(4456700000),
			wantSessionRelays:   sdk.NewInt(185695833),
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctx, _, keeper := createTestInput(t, true)
			if tt.stabilityAdjustment != 0 {
				keeper.Paramstore.Set(ctx, types.StabilityAdjustment, tt.stabilityAdjustment)
			}
			if tt.baseRelaysPerPOKT != 0 {
				keeper.Paramstore.Set(ctx, types.BaseRelaysPerPOKT, tt.baseRelaysPerPOKT)
			}
			if tt.participationRateOn {
				keeper.Paramstore.Set(ctx, types.ParticipationRateOn, tt.participationRateOn)
			}

			application := getStakedApplication()
			if tt.appStake.IsInt64() {
				application.StakedTokens = tt.appStake
			}
			if tt.appChains != nil {
				application.Chains = tt.appChains
			}

			gotAppRelays := keeper.CalculateAppRelays(ctx, application)
			if !gotAppRelays.Equal(tt.wantAppRelays) {
				t.Errorf("Application.CalculateAppRelays() = got %v, want %v", gotAppRelays, tt.wantAppRelays)
			}
			application.MaxRelays = gotAppRelays

			gotMaxPossibleRelays := coreTypes.MaxPossibleRelays(application, tt.sessionNodeCount)
			if !gotMaxPossibleRelays.Equal(tt.wantSessionRelays) {
				t.Errorf("Application.MaxPossibleRelays() = got %v, want %v", gotMaxPossibleRelays, tt.wantSessionRelays)
			}
		})
	}
}

func TestApplication_GetAllAplications(t *testing.T) {
	application := getStakedApplication()

	tests := []struct {
		name        string
		application types.Application
		want        types.Applications
	}{
		{
			name:        "gets all applications",
			application: application,
			want:        types.Applications([]types.Application{application}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.SetApplication(context, tt.application)

			if got := keeper.GetAllApplications(context); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Applicaiton.GetAllApplications() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_GetAplications(t *testing.T) {
	application := getStakedApplication()

	tests := []struct {
		name        string
		application types.Application
		maxRetrieve uint16
		want        types.Applications
	}{
		{
			name:        "gets all applications",
			application: application,
			maxRetrieve: 2,
			want:        types.Applications([]types.Application{application}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.SetApplication(context, tt.application)

			if got := keeper.GetApplications(context, tt.maxRetrieve); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Applicaiton.GetAllApplications() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplication_IterateAndExecuteOverApps(t *testing.T) {
	application := getStakedApplication()
	secondApp := getStakedApplication()

	tests := []struct {
		name              string
		application       types.Application
		secondApplication types.Application
		want              int
	}{
		{
			name:              "iterates over all applications",
			application:       application,
			secondApplication: secondApp,
			want:              2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.SetApplication(context, tt.application)
			keeper.SetApplication(context, tt.secondApplication)
			got := 0
			fn := modifyFn(&got)
			keeper.IterateAndExecuteOverApps(context, fn)
			if got != tt.want {
				t.Errorf("Application.IterateAndExecuteOverApps() = got %v, want %v", got, tt.want)
			}
		})
	}
}
