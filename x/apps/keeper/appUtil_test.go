package keeper

import (
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/x/apps/exported"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

func TestAppUtil_Application(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
	}
	type want struct {
		application types.Application
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

func TestAppUtil_ApplicationCaching(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		application types.Application
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
			if got, _ := keeper.ApplicationCache.Get(tt.args.application.Address.String()); !got.(types.Application).Equals(tt.want) {
				t.Errorf("keeperAppUtil.ApplicationCaching()= %v, want %v", got, tt.want)
			}
		})
	}
}
