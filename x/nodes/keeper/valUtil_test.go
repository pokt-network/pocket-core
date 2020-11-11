package keeper

import (
	types2 "github.com/pokt-network/pocket-core/types"
	"testing"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
)

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
			validator, _ := keeper.validatorCache.Get(test.args.validator.Address.String())
			assert.True(t, validator.(types.Validator).Equals(test.expected.validator), "validator does not match")
		})
	}

}

func TestValidatorCachingAfterUpdate(t *testing.T) {
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
			validator, _ := keeper.validatorCache.Get(test.args.validator.Address.String())
			assert.True(t, validator.(types.Validator).Equals(test.expected.validator), "validator does not match")
			modifiedVal := test.args.validator
			modifiedVal.Chains = []string{"00", "01", "03"}
			modifiedVal.UpdateStatus(types2.Unstaking)
			keeper.SetValidator(context, test.args.validator)
			validator2, _ := keeper.validatorCache.Get(test.args.validator.Address.String())
			assert.True(t, validator2.(types.Validator).Equals(modifiedVal), "validator does not match")
		})
	}

}
