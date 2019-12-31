package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetStakedApplication(t *testing.T) {
	boundedApplication := getBondedApplication()
	unboundedApplication := getUnbondedApplication()

	type expected struct {
		applications []types.Application
		length       int
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected     expected
	}{
		{
			name:         "gets applications",
			applications: []types.Application{boundedApplication},
			expected:     expected{applications: []types.Application{boundedApplication}, length: 1},
		},
		{
			name:         "gets emtpy slice of applications",
			applications: []types.Application{unboundedApplication},
			expected:     expected{applications: []types.Application{}, length: 0},
		},
		{
			name:         "only gets bounded applications",
			applications: []types.Application{boundedApplication, unboundedApplication},
			expected:     expected{applications: []types.Application{boundedApplication}, length: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.applications {
				keeper.SetApplication(context, application)
				keeper.SetStakedApplication(context, application)
			}
			applications := keeper.getStakedApplications(context)

			if equal := assert.ObjectsAreEqualValues(applications, test.expected.applications); !equal { // note ObjectsAreEqualValues does not assert, manual verification is required
				t.FailNow()
			}
			assert.Equalf(t, len(applications), test.expected.length, "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestRemoveStakedApplicationTokens(t *testing.T) {
	boundedApplication := getBondedApplication()

	type expected struct {
		tokens       sdk.Int
		applications []types.Application
		errorMessage string
	}
	tests := []struct {
		name        string
		application types.Application
		panics      bool
		amount      sdk.Int
		expected
	}{
		{
			name:        "removes tokens from application applications",
			application: boundedApplication,
			amount:      sdk.NewInt(5),
			panics:      false,
			expected:    expected{tokens: sdk.NewInt(99999999995), applications: []types.Application{}},
		},
		{
			name:        "removes tokens from application applications",
			application: boundedApplication,
			amount:      sdk.NewInt(-5),
			panics:      true,
			expected:    expected{tokens: sdk.NewInt(99999999995), applications: []types.Application{}, errorMessage: "trying to remove negative tokens"},
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
					assert.Contains(t, err, test.expected.errorMessage)
				}()
				_ = keeper.removeApplicationTokens(context, test.application, test.amount)
			default:
				application := keeper.removeApplicationTokens(context, test.application, test.amount)
				assert.True(t, application.StakedTokens.Equal(test.expected.tokens), "application staked tokens is not as expected")

				store := context.KVStore(keeper.storeKey)
				assert.NotNil(t, store.Get(types.KeyForAppInStakingSet(application)))
			}
		})
	}
}

func TestRemoveDeleteFromStakingSet(t *testing.T) {
	boundedApplication := getBondedApplication()
	unboundedApplication := getUnbondedApplication()

	tests := []struct {
		name         string
		applications []types.Application
		panics       bool
		amount       sdk.Int
	}{
		{
			name:         "removes applications from set",
			applications: []types.Application{boundedApplication, unboundedApplication},
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
	boundedApplication := getBondedApplication()
	unboundedApplication := getUnbondedApplication()

	tests := []struct {
		name         string
		applications []types.Application
		panics       bool
		amount       sdk.Int
	}{
		{
			name:         "recieves a valid iterator",
			applications: []types.Application{boundedApplication, unboundedApplication},
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
