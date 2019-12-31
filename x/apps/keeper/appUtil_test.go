package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustGetApplication(t *testing.T) {
	boundedApplication := getBondedApplication()

	type args struct {
		application types.Application
	}
	type expected struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets application",
			panics:   false,
			args:     args{application: boundedApplication},
			expected: expected{application: boundedApplication},
		},
		{
			name:     "panics if no application",
			panics:   true,
			args:     args{application: boundedApplication},
			expected: expected{message: fmt.Sprintf("application record not found for address: %X\n", boundedApplication.Address)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.panics {
			case true:
				defer func() {
					err := recover()
					assert.Contains(t, test.expected.message, err, "does not cointain error message")
				}()
				_ = keeper.mustGetApplication(context, test.args.application.Address)
			default:
				keeper.SetApplication(context, test.args.application)
				keeper.SetStakedApplication(context, test.args.application)
				application := keeper.mustGetApplication(context, test.args.application.Address)
				assert.True(t, application.Equals(test.expected.application), "application does not match")
			}
		})
	}

}

func TestMustGetApplicationByConsAddr(t *testing.T) {
	boundedApplication := getBondedApplication()
	type args struct {
		application types.Application
	}
	type expected struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets application",
			panics:   false,
			args:     args{application: boundedApplication},
			expected: expected{application: boundedApplication},
		},
		{
			name:     "panics if no application",
			panics:   true,
			args:     args{application: boundedApplication},
			expected: expected{message: fmt.Sprintf("application with consensus-Address %s not found", boundedApplication.ConsAddress())},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Equal(t, test.expected.message, err.Error(), "messages don't match")
				}()
				_ = keeper.mustGetApplicationByConsAddr(context, test.args.application.ConsAddress())
			default:
				keeper.SetApplication(context, test.args.application)
				keeper.SetAppByConsAddr(context, test.args.application)
				keeper.SetStakedApplication(context, test.args.application)
				application := keeper.mustGetApplicationByConsAddr(context, test.args.application.ConsAddress())
				assert.True(t, application.Equals(test.expected.application), "application does not match")
			}
		})
	}

}

func TestApplicationByConsAddr(t *testing.T) {
	boundedApplication := getBondedApplication()

	type args struct {
		application types.Application
	}
	type expected struct {
		application types.Application
		message     string
		null        bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets application",
			args:     args{application: boundedApplication},
			expected: expected{application: boundedApplication, null: false},
		},
		{
			name:     "nil if not found",
			args:     args{application: boundedApplication},
			expected: expected{null: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.expected.null {
			case true:
				application := keeper.applicationByConsAddr(context, test.args.application.ConsAddress())
				assert.Nil(t, application)
			default:
				keeper.SetApplication(context, test.args.application)
				keeper.SetAppByConsAddr(context, test.args.application)
				keeper.SetStakedApplication(context, test.args.application)
				application := keeper.applicationByConsAddr(context, test.args.application.ConsAddress())
				assert.Equal(t, application, test.expected.application, "application does not match")
			}
		})
	}
}

func TestApplicationCaching(t *testing.T) {
	boundedApplication := getBondedApplication()

	type args struct {
		bz          []byte
		application types.Application
	}
	type expected struct {
		application types.Application
		message     string
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets application",
			panics:   false,
			args:     args{application: boundedApplication},
			expected: expected{application: boundedApplication},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, test.args.application)
			keeper.SetStakedApplication(context, test.args.application)
			store := context.KVStore(keeper.storeKey)
			bz := store.Get(types.KeyForAppByAllApps(test.args.application.Address))
			application := keeper.appCaching(bz, test.args.application.Address)
			assert.True(t, application.Equals(test.expected.application), "application does not match")
		})
	}

}

//func TestNewApplicationCaching(t *testing.T) { todo
//	boundedApplication := getBondedApplication()
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
//		panics bool
//		args
//		expected
//	}{
//		{
//			name:     "getPrevStatePowerMap",
//			panics:   false,
//			args:     args{application: boundedApplication},
//			expected: expected{application: boundedApplication, length: 1},
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
