package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"reflect"
	"testing"
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
	application := getStakedApplication()

	tests := []struct {
		name        string
		application types.Application
		want        sdk.BigInt
	}{
		{
			name:        "calculates App relays",
			application: application,
			want:        sdk.NewInt(100000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			if got := keeper.CalculateAppRelays(context, tt.application); !got.Equal(tt.want) {
				t.Errorf("Applicaiton.CalculateAppRelays() = got %v, want %v", got, tt.want)
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
