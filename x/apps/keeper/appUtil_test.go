package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"reflect"
	"strings"
	"testing"
)

func TestAppUtil_MustGetApplication(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
	}
	type want struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name   string
		panics bool
		args
		want
	}{
		{
			name:   "gets application",
			panics: false,
			args:   args{application: stakedApplication},
			want:   want{application: stakedApplication},
		},
		{
			name:   "errors if no application",
			panics: true,
			args:   args{application: stakedApplication},
			want:   want{message: fmt.Sprintf("application record not found for address: %X\n", stakedApplication.Address)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.panics {
			case true:
				defer func() {
					if got := recover(); !strings.Contains(got.(string), tt.want.message) {
						t.Errorf("keeperAppUtil.MustGetApplication()= %v, want %v", got, tt.want.application)
					}
				}()
				_ = keeper.mustGetApplication(context, tt.args.application.Address)
			default:
				keeper.SetApplication(context, tt.args.application)
				keeper.SetStakedApplication(context, tt.args.application)
				if got := keeper.mustGetApplication(context, tt.args.application.Address); !got.Equals(tt.want.application) {
					t.Errorf("keeperAppUtil.MustGetApplication()= %v, want %v", got, tt.want.application)
				}
			}
		})
	}

}

func TestAppUtil_Application(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
	}
	type want struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name string
		find bool
		args
		want
	}{
		{
			name: "gets application",
			find: false,
			args: args{application: stakedApplication},
			want: want{application: stakedApplication},
		},
		{
			name: "errors if no application",
			find: true,
			args: args{application: stakedApplication},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.find {
			case true:
				if got := keeper.Application(context, tt.args.application.Address); got != nil {
					t.Errorf("keeperAppUtil.Application()= %v, want nil", got)
				}
			default:
				keeper.SetApplication(context, tt.args.application)
				keeper.SetStakedApplication(context, tt.args.application)
				if got := keeper.Application(context, tt.args.application.Address); !reflect.DeepEqual(got, tt.want.application) {
					t.Errorf("keeperAppUtil.Application()= %v, want %v", got, tt.want.application)
				}
			}
		})
	}

}

func TestAppUtil_MustGetApplicationByConsAddr(t *testing.T) {
	stakedApplication := getStakedApplication()
	type args struct {
		application types.Application
	}
	type want struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name   string
		panics bool
		args
		want
	}{
		{
			name:   "gets application",
			panics: false,
			args:   args{application: stakedApplication},
			want:   want{application: stakedApplication},
		},
		{
			name:   "errors if no application",
			panics: true,
			args:   args{application: stakedApplication},
			want:   want{message: fmt.Sprintf("application with consensus-Address %s not found", stakedApplication.GetAddress())},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.panics {
			case true:
				defer func() {
					if err := recover().(error); !reflect.DeepEqual(err.Error(), tt.want.message) {
						t.Errorf("keeperAppUtil.MustGetApplicationByConsAddr()= %v, want %v", err, tt.want.application)
					}
				}()
				_ = keeper.mustGetApplicationByConsAddr(context, tt.args.application.GetAddress())
			default:
				keeper.SetApplication(context, tt.args.application)
				keeper.SetStakedApplication(context, tt.args.application)
				if got := keeper.mustGetApplicationByConsAddr(context, tt.args.application.GetAddress()); !got.Equals(tt.want.application) {
					t.Errorf("keeperAppUtil.MustGetApplicationByConsAddr()= %v, want %v", got, tt.want.application)
				}
			}
		})
	}

}

func TestAppUtil_ApplicationByConsAddr(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
	}
	type want struct {
		application types.Application
		message     string
		null        bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		want interface{}
	}{
		{
			name: "gets application",
			args: args{application: stakedApplication},
			want: stakedApplication,
		},
		{
			name: "nil if not found",
			args: args{application: stakedApplication},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.want {
			case nil:
				if got := keeper.applicationByConsAddr(context, tt.args.application.GetAddress()); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("keeperAppUtil.ApplicationByConsAddr()= %v, want %v", got, tt.want)
				}

			default:
				keeper.SetApplication(context, tt.args.application)
				keeper.SetStakedApplication(context, tt.args.application)
				if got := keeper.applicationByConsAddr(context, tt.args.application.GetAddress()); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("keeperAppUtil.ApplicationByConsAddr()= %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAppUtil_ApplicationCaching(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		bz             []byte
		application    types.Application
		aminoCacheSize int
	}
	tests := []struct {
		name   string
		panics bool
		args
		want types.Application
	}{
		{
			name: "gets application",
			args: args{application: stakedApplication},
			want: stakedApplication,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.args.application)
			keeper.SetStakedApplication(context, tt.args.application)
			store := context.KVStore(keeper.storeKey)
			bz := store.Get(types.KeyForAppByAllApps(tt.args.application.Address))
			if got := keeper.appCaching(bz, tt.args.application.Address); !got.Equals(tt.want) {
				t.Errorf("keeperAppUtil.ApplicationCaching()= %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppUtil_AllApplications(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected []exported.ApplicationI
	}{
		{
			name:     "gets application",
			panics:   false,
			args:     args{application: stakedApplication},
			expected: []exported.ApplicationI{stakedApplication},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.args.application)
			keeper.SetStakedApplication(context, tt.args.application)
			if got := keeper.AllApplications(context); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("keeperAppUtil.AllApplications()= %v, want %v", got, tt.expected)
			}
		})
	}
}

//func TestNewApplicationCaching(t *testing.T) { todo
//	stakedApplication := getStakedApplication()
//
//	type args struct {
//		bz        []byte
//		application types.Application
//	}
//	type expected struct {
//		application types.Application
//		message   string
//		length    int
//	}
//	tests := []struct {
//		name   string
//		errors bool
//		args
//		expected
//	}{
//		{
//			name:     "getPrevStatePowerMap",
//			errors:   false,
//			args:     args{application: stakedApplication},
//			expected: expected{application: stakedApplication, length: 1},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			context, _, keeper := createTestInput(t, true)
//			keeper.SetApplication(context, test.args.application)
//			keeper.SetStakedApplication(context, test.args.application)
//			store := context.KVStore(keeper.storeKey)
//			key := types.KeyForApplicationPrevStateStateByPower(test.args.application.Address)
//			store.Set(key, test.args.application.Address)
//			powermap := keeper.getPrevStatePowerMap(context)
//			assert.Len(t, powermap, test.expected.length, "does not have correct length")
//			var valAddr [sdk.AddrLen]byte
//			copy(valAddr[:], key[1:])
//
//			for mapKey, value := range powermap {
//				assert.Equal(t, valAddr, mapKey, "key is not correct")
//				bz := make([]byte, len(test.args.application.Address))
//				copy(bz, test.args.application.Address)
//				assert.Equal(t, bz, value, "key is not correct")
//			}
//		})
//	}
//}
