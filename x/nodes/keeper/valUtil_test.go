package keeper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
)

func TestMustGetValidator(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
	}
	tests := []struct {
		name     string
		hasError bool
		args
		expected
	}{
		{
			name:     "gets validator",
			hasError: false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator},
		},
		{
			name:     "errors if no validator",
			hasError: true,
			args:     args{validator: stakedValidator},
			expected: expected{message: fmt.Sprintf("validator record not found for address: %X\n", stakedValidator.Address)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.hasError {
			case true:
				_, _ = keeper.GetValidator(context, test.args.validator.Address)
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator, _ := keeper.GetValidator(context, test.args.validator.Address)
				assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
			}
		})
	}

}

func TestValidatorCaching(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets validator",
			panics:   false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.args.validator)
			keeper.SetStakedValidator(context, test.args.validator)
			store := context.KVStore(keeper.storeKey)
			bz := store.Get(types.KeyForValByAllVals(test.args.validator.Address))
			validator := keeper.validatorCaching(bz, test.args.validator.Address)
			assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
		})
	}

}

func TestNewValidatorCaching(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		length    int
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "getPrevStatePowerMap",
			panics:   false,
			args:     args{validator: stakedValidator},
			expected: expected{validator: stakedValidator, length: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.args.validator)
			keeper.SetStakedValidator(context, test.args.validator)
			store := context.KVStore(keeper.storeKey)
			key := types.KeyForValidatorPrevStateStateByPower(test.args.validator.Address)
			store.Set(key, test.args.validator.Address)
			powermap := keeper.getPrevStatePowerMap(context)
			assert.Len(t, powermap, test.expected.length, "does not have correct length")
			var valAddr [sdk.AddrLen]byte
			copy(valAddr[:], key[1:])

			for mapKey, value := range powermap {
				assert.Equal(t, valAddr, mapKey, "key is not correct")
				bz := make([]byte, len(test.args.validator.Address))
				copy(bz, test.args.validator.Address)
				assert.Equal(t, bz, value, "key is not correct")
			}
		})
	}
}

func Test_sortNoLongerStakedValidators(t *testing.T) {
	type args struct {
		prevState valPowerMap
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortNoLongerStakedValidators(tt.args.prevState); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortNoLongerStakedValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}
