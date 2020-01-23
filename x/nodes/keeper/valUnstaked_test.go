package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//
//func TestGetAndSetlUnstaking(t *testing.T) {
//	boundedValidator := getBondedValidator()
//	secondaryBoundedValidator := getBondedValidator()
//	stakedValidator := getBondedValidator()
//
//	type expected struct {
//		validators       []types.Validator
//		stakedValidators bool
//		length           int
//	}
//	type args struct {
//		boundedVal      types.Validator
//		validators      []types.Validator
//		stakedValidator types.Validator
//	}
//	tests := []struct {
//		name       string
//		validator  types.Validator
//		validators []types.Validator
//		expected
//		args
//	}{
//		{
//			name:     "gets validators",
//			args:     args{validators: []types.Validator{boundedValidator}},
//			expected: expected{validators: []types.Validator{boundedValidator}, length: 1, stakedValidators: false},
//		},
//		{
//			name:     "gets emtpy slice of validators",
//			expected: expected{length: 0, stakedValidators: true},
//			args:     args{stakedValidator: stakedValidator},
//		},
//		{
//			name:       "only gets unstakedbounded validators",
//			validators: []types.Validator{boundedValidator, secondaryBoundedValidator},
//			expected:   expected{length: 1, stakedValidators: true},
//			args:       args{stakedValidator: stakedValidator, validators: []types.Validator{boundedValidator}},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			context, _, keeper := createTestInput(t, true)
//			for _, validator := range test.args.validators {
//				keeper.SetValidator(context, validator)
//				keeper.SetUnstakingValidator(context, validator)
//			}
//			if test.expected.stakedValidators {
//				keeper.SetValidator(context, test.args.stakedValidator)
//				keeper.SetStakedValidator(context, test.args.stakedValidator)
//			}
//			validators := keeper.getAllUnstakingValidators(context)
//
//			for _, validator := range validators {
//				assert.True(t, validator.Status.Equal(sdk.Unbonded))
//			}
//			assert.Equalf(t, test.expected.length, len(validators), "length of the validators does not match expected on %v", test.name)
//		})
//	}
//}

func TestDeleteUnstakingValidator(t *testing.T) {
	boundedValidator := getBondedValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		boundedVal      types.Validator
		validators      []types.Validator
		stakedValidator types.Validator
	}
	tests := []struct {
		name       string
		validator  types.Validator
		validators []types.Validator
		expected
		args
	}{
		{
			name:     "deletes validator",
			args:     args{validators: []types.Validator{boundedValidator}},
			expected: expected{length: 0, stakedValidators: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.args.validators {
				keeper.SetValidator(context, validator)
				keeper.SetUnstakingValidator(context, validator)
				keeper.deleteUnstakingValidator(context, validator)
			}
			if test.expected.stakedValidators {
				keeper.SetValidator(context, test.args.stakedValidator)
				keeper.SetStakedValidator(context, test.args.stakedValidator)
			}

			validators := keeper.getAllUnstakingValidators(context)

			assert.Equalf(t, test.expected.length, len(validators), "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestDeleteUnstakingValidators(t *testing.T) {
	boundedValidator := getBondedValidator()
	secondaryBoundedValidator := getBondedValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		boundedVal      types.Validator
		validators      []types.Validator
		stakedValidator types.Validator
	}
	tests := []struct {
		name       string
		validator  types.Validator
		validators []types.Validator
		expected
		args
	}{
		{
			name:     "deletes all unstaking validator",
			args:     args{validators: []types.Validator{boundedValidator, secondaryBoundedValidator}},
			expected: expected{length: 0, stakedValidators: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.args.validators {
				keeper.SetValidator(context, validator)
				keeper.SetUnstakingValidator(context, validator)
				keeper.deleteUnstakingValidators(context, validator.UnstakingCompletionTime)
			}

			validators := keeper.getAllUnstakingValidators(context)

			assert.Equalf(t, test.expected.length, len(validators), "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestGetAllMatureValidators(t *testing.T) {
	unboundingValidator := getUnbondingValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		boundedVal      types.Validator
		validators      []types.Validator
		stakedValidator types.Validator
	}
	tests := []struct {
		name       string
		validator  types.Validator
		validators []types.Validator
		expected
		args
	}{
		{
			name:     "gets all mature validators",
			args:     args{validators: []types.Validator{unboundingValidator}},
			expected: expected{validators: []types.Validator{unboundingValidator}, length: 1, stakedValidators: false},
		},
		{
			name:     "gets empty slice if no mature validators",
			args:     args{validators: []types.Validator{}},
			expected: expected{validators: []types.Validator{unboundingValidator}, length: 0, stakedValidators: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.args.validators {
				keeper.SetValidator(context, validator)
				keeper.SetUnstakingValidator(context, validator)
			}
			keeper.UpdateTendermintValidators(context)
			matureValidators := keeper.getMatureValidators(context)

			assert.Equalf(t, test.expected.length, len(matureValidators), "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestUnstakeAllMatureValidators(t *testing.T) {
	unboundingValidator := getUnbondingValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		boundedVal      types.Validator
		validators      []types.Validator
		stakedValidator types.Validator
	}
	tests := []struct {
		name       string
		validator  types.Validator
		validators []types.Validator
		expected
		args
	}{
		{
			name:     "unstake mature validators",
			args:     args{validators: []types.Validator{unboundingValidator}},
			expected: expected{validators: []types.Validator{unboundingValidator}, length: 0, stakedValidators: false},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.args.validators {
				keeper.SetValidator(context, validator)
				keeper.SetUnstakingValidator(context, validator)
			}
			keeper.UpdateTendermintValidators(context)
			keeper.unstakeAllMatureValidators(context)
			validators := keeper.getAllUnstakingValidators(context)

			assert.Equalf(t, test.expected.length, len(validators), "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestUnstakingValidatorsIterator(t *testing.T) {
	boundedValidator := getBondedValidator()
	unboundedValidator := getUnbondedValidator()

	tests := []struct {
		name       string
		validators []types.Validator
		panics     bool
		amount     sdk.Int
	}{
		{
			name:       "recieves a valid iterator",
			validators: []types.Validator{boundedValidator, unboundedValidator},
			panics:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.validators {
				keeper.SetValidator(context, validator)
				keeper.SetStakedValidator(context, validator)
			}

			it := keeper.unstakingValidatorsIterator(context, context.BlockHeader().Time)
			assert.Implements(t, (*sdk.Iterator)(nil), it, "does not implement interface")
		})
	}
}
