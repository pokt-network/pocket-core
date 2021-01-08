package keeper

import (
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/stretchr/testify/assert"
)

func TestAppUnstaked_GetAndSetlUnstaking(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstaking := getUnstakingApplication()

	type want struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		applications      []types.Application
		stakedApplication types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		want
		args
	}{
		{
			name: "gets applications",
			args: args{applications: []types.Application{unstaking}},
			want: want{applications: []types.Application{unstaking}, length: 1, stakedApplications: false},
		},
		{
			name: "gets emtpy slice of applications",
			want: want{length: 0, stakedApplications: true},
			args: args{stakedApplication: stakedApplication},
		},
		{
			name:         "only gets unstaking applications",
			applications: []types.Application{stakedApplication, unstaking},
			want:         want{length: 1, stakedApplications: true},
			args:         args{stakedApplication: stakedApplication, applications: []types.Application{unstaking}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.args.applications {
				keeper.SetApplication(context, application)
			}
			if tt.want.stakedApplications {
				keeper.SetApplication(context, tt.args.stakedApplication)
			}
			applications := keeper.getAllUnstakingApplications(context)
			if len(applications) != tt.want.length {
				t.Errorf("appUnstaked.GetApplications() = %v, want %v", len(applications), tt.want.length)
			}
		})
	}
}

func TestAppUnstaked_DeleteUnstakingApplication(t *testing.T) {
	stakedApplication := getStakedApplication()
	secondStakedApp := getStakedApplication()

	type want struct {
		stakedApplications bool
		length             int
	}
	type args struct {
		applications []types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		sets         bool
		want
		args
	}{
		{
			name: "deletes",
			args: args{applications: []types.Application{stakedApplication, secondStakedApp}},
			sets: false,
			want: want{length: 1, stakedApplications: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.args.applications {
				keeper.SetApplication(context, application)
			}
			keeper.SetUnstakingApplication(context, tt.args.applications[0])
			_ = keeper.getAllUnstakingApplications(context)

			keeper.deleteUnstakingApplication(context, tt.args.applications[1])

			if got := keeper.getAllUnstakingApplications(context); len(got) != tt.want.length {
				t.Errorf("KeeperCoins.BurnStakedTokens()= %v, want %v", len(got), tt.want.length)
			}
		})
	}
}

func TestAppUnstaked_DeleteUnstakingApplications(t *testing.T) {
	stakedApplication := getStakedApplication()
	secondaryStakedApplication := getStakedApplication()

	type want struct {
		stakedApplications bool
		length             int
	}
	type args struct {
		applications []types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		want
		args
	}{
		{
			name: "deletes all unstaking application",
			args: args{applications: []types.Application{stakedApplication, secondaryStakedApplication}},
			want: want{length: 0, stakedApplications: false},
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

			assert.Equalf(t, test.want.length, len(applications), "length of the applications does not match want on %v", test.name)
		})
	}
}

func TestAppUnstaked_GetAllMatureApplications(t *testing.T) {
	stakingApplication := getUnstakingApplication()

	type want struct {
		applications       []types.Application
		stakedApplications bool
		length             int
	}
	type args struct {
		applications []types.Application
	}
	tests := []struct {
		name         string
		application  types.Application
		applications []types.Application
		want
		args
	}{
		{
			name: "gets all mature applications",
			args: args{applications: []types.Application{stakingApplication}},
			want: want{applications: []types.Application{stakingApplication}, length: 1, stakedApplications: false},
		},
		{
			name: "gets empty slice if no mature applications",
			args: args{applications: []types.Application{}},
			want: want{applications: []types.Application{stakingApplication}, length: 0, stakedApplications: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.args.applications {
				keeper.SetApplication(context, application)
			}
			if got := keeper.getMatureApplications(context); len(got) != tt.want.length {
				t.Errorf("appUnstaked.unstakeAllMatureApplications()= %v, want %v", len(got), tt.want.length)
			}
		})
	}
}

//func TestAppUnstaked_UnstakeAllMatureApplications(t *testing.T) {
//	stakingApplication := getUnstakingApplication()
//
//	type want struct {
//		applications       []types.Application
//		stakedApplications bool
//		length             int
//	}
//	type args struct {
//		stakedVal         types.Application
//		applications      []types.Application
//		stakedApplication types.Application
//	}
//	tests := []struct {
//		name         string
//		application  types.Application
//		applications []types.Application
//		want
//		args
//	}{
//		{
//			name: "unstake mature applications",
//			args: args{applications: []types.Application{stakingApplication}},
//			want: want{applications: []types.Application{stakingApplication}, length: 0, stakedApplications: false},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			context, _, keeper := createTestInput(t, true)
//			for _, application := range tt.args.applications {
//				keeper.SetApplication(context, application)
//				keeper.SetUnstakingApplication(context, application)
//			}
//			keeper.unstakeAllMatureApplications(context)
//			if got := keeper.getAllUnstakingApplications(context); len(got) != tt.want.length {
//				t.Errorf("appUnstaked.unstakeAllMatureApplications()= %v, want %v", len(got), tt.want.length)
//			}
//		})
//	}
//}

func TestAppUnstaked_UnstakingApplicationsIterator(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstakedApplication := getUnstakedApplication()

	tests := []struct {
		name         string
		applications []types.Application
		panics       bool
		amount       sdk.BigInt
	}{
		{
			name:         "recieves a valid iterator",
			applications: []types.Application{stakedApplication, unstakedApplication},
			panics:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.applications {
				keeper.SetApplication(context, application)
				keeper.SetStakedApplication(context, application)
			}

			it, _ := keeper.unstakingApplicationsIterator(context, context.BlockHeader().Time)
			if v, ok := it.(sdk.Iterator); !ok {
				t.Errorf("appUnstaked.UnstakingApplicationsIterator()= %v does not implement sdk.Iterator", v)
			}
		})
	}
}
