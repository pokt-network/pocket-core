package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetStakedApplication(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstakedApplication := getUnstakedApplication()
	jailedApp := getStakedApplication()
	jailedApp.Jailed = true

	type want struct {
		applications []types.Application
		length       int
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		want         want
	}{
		{
			name:         "gets applications",
			applications: []types.Application{stakedApplication},
			want:         want{applications: []types.Application{stakedApplication}, length: 1},
		},
		{
			name:         "gets emtpy slice of applications",
			applications: []types.Application{unstakedApplication},
			want:         want{applications: []types.Application{}, length: 0},
		},
		{
			name:         "gets emtpy slice of applications",
			applications: []types.Application{jailedApp},
			want:         want{applications: []types.Application{}, length: 0},
		},
		{
			name:         "only gets staked applications",
			applications: []types.Application{stakedApplication, unstakedApplication},
			want:         want{applications: []types.Application{stakedApplication}, length: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.applications {
				keeper.SetApplication(context, application)
				if application.IsStaked() {
					keeper.SetStakedApplication(context, application)
				}
			}
			applications := keeper.getStakedApplications(context)
			if equal := assert.ObjectsAreEqualValues(applications, test.want.applications); !equal { // note ObjectsAreEqualValues does not assert, manual verification is required
				t.FailNow()
			}
			assert.Equalf(t, len(applications), test.want.length, "length of the applications does not match want on %v", test.name)
		})
	}
}

func TestRemoveStakedApplicationTokens(t *testing.T) {
	stakedApplication := getStakedApplication()

	type want struct {
		tokens       sdk.Int
		applications []types.Application
		errorMessage string
	}
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		amount      sdk.Int
		want
	}{
		{
			name:        "removes tokens from application applications",
			application: stakedApplication,
			amount:      sdk.NewInt(5),
			panics:      false,
			want:        want{tokens: sdk.NewInt(99999999995), applications: []types.Application{}},
		},
		{
			name:        "removes tokens from application applications",
			application: stakedApplication,
			amount:      sdk.NewInt(-5),
			panics:      true,
			want:        want{tokens: sdk.NewInt(99999999995), applications: []types.Application{}, errorMessage: "trying to remove negative tokens"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, test.application)
			keeper.SetStakedApplication(context, test.application)
			switch test.panics {
			case true:
				defer func() {
					err := recover()
					assert.Contains(t, err, test.want.errorMessage)
				}()
				_ = keeper.removeApplicationTokens(context, test.application, test.amount)
			default:
				application := keeper.removeApplicationTokens(context, test.application, test.amount)
				assert.True(t, application.StakedTokens.Equal(test.want.tokens), "application staked tokens is not as want")

				store := context.KVStore(keeper.storeKey)
				assert.NotNil(t, store.Get(types.KeyForAppInStakingSet(application)))
			}
		})
	}
}

func TestRemoveDeleteFromStakingSet(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstakedApplication := getUnstakedApplication()

	tests := []struct {
		name         string
		applications []types.Application
		panics       bool
		amount       sdk.Int
	}{
		{
			name:         "removes applications from set",
			applications: []types.Application{stakedApplication, unstakedApplication},
			panics:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.applications {
				keeper.SetApplication(context, application)
				keeper.SetStakedApplication(context, application)
			}
			for _, application := range test.applications {
				keeper.deleteApplicationFromStakingSet(context, application)
			}

			applications := keeper.getStakedApplications(context)
			assert.Empty(t, applications, "there should not be any applications in the set")
		})
	}
}

func TestGetValsIterator(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstakedApplication := getUnstakedApplication()

	tests := []struct {
		name         string
		applications []types.Application
		panics       bool
		amount       sdk.Int
	}{
		{
			name:         "recieves a valid iterator",
			applications: []types.Application{stakedApplication, unstakedApplication},
			panics:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.applications {
				keeper.SetApplication(context, application)
				keeper.SetStakedApplication(context, application)
			}

			it := keeper.stakedAppsIterator(context)
			assert.Implements(t, (*sdk.Iterator)(nil), it, "does not implement interface")
		})
	}
}

func TestApplicationStaked_IterateAndExecuteOverStakedApps(t *testing.T) {
	stakedApplication := getStakedApplication()
	secondStakedApplication := getStakedApplication()

	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		want         int
	}{
		{
			name:         "iterates over applications",
			applications: []types.Application{stakedApplication, secondStakedApplication},
			want:         2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.applications {
				keeper.SetApplication(context, application)
				keeper.SetStakedApplication(context, application)
			}
			got := 0
			fn := modifyFn(&got)

			keeper.IterateAndExecuteOverStakedApps(context, fn)

			if got != tt.want {
				t.Errorf("appStaked.IterateAndExecuteOverApps() = got %v, want %v", got, tt.want)
			}
		})
	}
}
func TestApplicationStaked_RemoveApplicationRelays(t *testing.T) {
	stakedApplication := getStakedApplication()

	tests := []struct {
		name        string
		application types.Application
		remove      sdk.Int
		want        sdk.Int
	}{
		{
			name:        "iterates over applications",
			application: stakedApplication,
			remove:      sdk.NewInt(100000000000 / 2),
			want:        sdk.NewInt(100000000000 / 2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.application)
			keeper.SetStakedApplication(context, tt.application)

			if got := keeper.removeApplicationRelays(context, tt.application, tt.remove); !got.MaxRelays.Equal(tt.want) {
				t.Errorf("appStaked.RemoveApplicationRelays() = got %v, want %v", got.MaxRelays, tt.want)
			}
		})
	}
}
