package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustGetValidator(t *testing.T) {
	boundedValidator := getBondedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
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
			args:     args{validator: boundedValidator},
			expected: expected{validator: boundedValidator},
		},
		{
			name:     "panics if no validator",
			panics:   true,
			args:     args{validator: boundedValidator},
			expected: expected{message: fmt.Sprintf("validator record not found for address: %X\n", boundedValidator.Address)},
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
				_ = keeper.mustGetValidator(context, test.args.validator.Address)
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator := keeper.mustGetValidator(context, test.args.validator.Address)
				assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
			}
		})
	}

}

func TestMustGetValidatorByConsAddr(t *testing.T) {
	boundedValidator := getBondedValidator()
	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
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
			args:     args{validator: boundedValidator},
			expected: expected{validator: boundedValidator},
		},
		{
			name:     "panics if no validator",
			panics:   true,
			args:     args{validator: boundedValidator},
			expected: expected{message: fmt.Sprintf("validator with consensus-Address %s not found", boundedValidator.ConsAddress())},
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
				_ = keeper.mustGetValidatorByConsAddr(context, test.args.validator.ConsAddress())
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetValidatorByConsAddr(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator := keeper.mustGetValidatorByConsAddr(context, test.args.validator.ConsAddress())
				assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
			}
		})
	}

}

func TestValidatorByConsAddr(t *testing.T) {
	boundedValidator := getBondedValidator()

	type args struct {
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
		null      bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "gets validator",
			args:     args{validator: boundedValidator},
			expected: expected{validator: boundedValidator, null: false},
		},
		{
			name:     "nil if not found",
			args:     args{validator: boundedValidator},
			expected: expected{null: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch test.expected.null {
			case true:
				validator := keeper.validatorByConsAddr(context, test.args.validator.ConsAddress())
				assert.Nil(t, validator)
			default:
				keeper.SetValidator(context, test.args.validator)
				keeper.SetValidatorByConsAddr(context, test.args.validator)
				keeper.SetStakedValidator(context, test.args.validator)
				validator := keeper.validatorByConsAddr(context, test.args.validator.ConsAddress())
				assert.Equal(t, validator, test.expected.validator, "validator does not match")
			}
		})
	}
}

func TestValidatorCaching(t *testing.T) {
	boundedValidator := getBondedValidator()

	type args struct {
		bz        []byte
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
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
			args:     args{validator: boundedValidator},
			expected: expected{validator: boundedValidator},
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
	boundedValidator := getBondedValidator()

	type args struct {
		bz        []byte
		validator types.Validator
	}
	type expected struct {
		validator types.Validator
		message   string
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
			args:     args{validator: boundedValidator},
			expected: expected{validator: boundedValidator, length: 1},
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
