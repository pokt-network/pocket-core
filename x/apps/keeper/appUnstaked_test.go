package keeper

import (
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetlUnstaking(t *testing.T) {
	boundedApplication := getBondedApplication()
	secondaryBoundedApplication := getBondedApplication()
	stakedApplication := getBondedApplication()

	type expected struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		boundedVal        types.Application
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected
		args
	}{
		{
			name:     "gets applications",
			args:     args{applications: []types.Application{boundedApplication}},
			expected: expected{applications: []types.Application{boundedApplication}, length: 1, stakedApplications: false},
		},
		{
			name:     "gets emtpy slice of applications",
			expected: expected{length: 0, stakedApplications: true},
			args:     args{stakedApplication: stakedApplication},
		},
		{
			name:         "only gets unstakedbounded applications",
			applications: []types.Application{boundedApplication, secondaryBoundedApplication},
			expected:     expected{length: 1, stakedApplications: true},
			args:         args{stakedApplication: stakedApplication, applications: []types.Application{boundedApplication}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.args.applications {
				keeper.SetApplication(context, application)
				keeper.SetUnstakingApplication(context, application)
			}
			if test.expected.stakedApplications {
				keeper.SetApplication(context, test.args.stakedApplication)
				keeper.SetStakedApplication(context, test.args.stakedApplication)
			}
			applications := keeper.getAllUnstakingApplications(context)

			for _, application := range applications {
				assert.True(t, application.Status.Equal(sdk.Unbonded))
			}
			assert.Equalf(t, test.expected.length, len(applications), "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestDeleteUnstakingApplication(t *testing.T) {
	boundedApplication := getBondedApplication()

	type expected struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		boundedVal        types.Application
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected
		args
	}{
		{
			name:     "deletes application",
			args:     args{applications: []types.Application{boundedApplication}},
			expected: expected{length: 0, stakedApplications: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.args.applications {
				keeper.SetApplication(context, application)
				keeper.SetUnstakingApplication(context, application)
				keeper.deleteUnstakingApplication(context, application)
			}
			if test.expected.stakedApplications {
				keeper.SetApplication(context, test.args.stakedApplication)
				keeper.SetStakedApplication(context, test.args.stakedApplication)
			}

			applications := keeper.getAllUnstakingApplications(context)

			assert.Equalf(t, test.expected.length, len(applications), "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestDeleteUnstakingApplications(t *testing.T) {
	boundedApplication := getBondedApplication()
	secondaryBoundedApplication := getBondedApplication()

	type expected struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		boundedVal        types.Application
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected
		args
	}{
		{
			name:     "deletes all unstaking application",
			args:     args{applications: []types.Application{boundedApplication, secondaryBoundedApplication}},
			expected: expected{length: 0, stakedApplications: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.args.applications {
				keeper.SetApplication(context, application)
				keeper.SetUnstakingApplication(context, application)
				keeper.deleteUnstakingApplications(context, application.UnstakingCompletionTime)
			}

			applications := keeper.getAllUnstakingApplications(context)

			assert.Equalf(t, test.expected.length, len(applications), "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestGetAllMatureApplications(t *testing.T) {
	unboundingApplication := getUnbondingApplication()

	type expected struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		boundedVal        types.Application
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected
		args
	}{
		{
			name:     "gets all mature applications",
			args:     args{applications: []types.Application{unboundingApplication}},
			expected: expected{applications: []types.Application{unboundingApplication}, length: 1, stakedApplications: false},
		},
		{
			name:     "gets empty slice if no mature applications",
			args:     args{applications: []types.Application{}},
			expected: expected{applications: []types.Application{unboundingApplication}, length: 0, stakedApplications: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.args.applications {
				keeper.SetApplication(context, application)
				keeper.SetUnstakingApplication(context, application)
			}
			matureApplications := keeper.getMatureApplications(context)
			assert.Equalf(t, test.expected.length, len(matureApplications), "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestUnstakeAllMatureApplications(t *testing.T) {
	unboundingApplication := getUnbondingApplication()

	type expected struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		boundedVal        types.Application
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		expected
		args
	}{
		{
			name:     "unstake mature applications",
			args:     args{applications: []types.Application{unboundingApplication}},
			expected: expected{applications: []types.Application{unboundingApplication}, length: 0, stakedApplications: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range test.args.applications {
				keeper.SetApplication(context, application)
				keeper.SetUnstakingApplication(context, application)
			}
			keeper.unstakeAllMatureApplications(context)
			applications := keeper.getAllUnstakingApplications(context)

			assert.Equalf(t, test.expected.length, len(applications), "length of the applications does not match expected on %v", test.name)
		})
	}
}

func TestUnstakingApplicationsIterator(t *testing.T) {
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

			it := keeper.unstakingApplicationsIterator(context, context.BlockHeader().Time)
			assert.Implements(t, (*sdk.Iterator)(nil), it, "does not implement interface")
		})
	}
}
