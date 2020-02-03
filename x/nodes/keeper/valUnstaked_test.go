package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetlUnstaking(t *testing.T) {
	unstakingValidator := getUnstakingValidator()
	secondaryStakedValidator := getStakedValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		stakedVal       types.Validator
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
			name:     "gets validators",
			args:     args{validators: []types.Validator{unstakingValidator}},
			expected: expected{validators: []types.Validator{unstakingValidator}, length: 1, stakedValidators: false},
		},
		{
			name:     "gets emtpy slice of validators",
			expected: expected{length: 0, stakedValidators: true},
			args:     args{stakedValidator: unstakingValidator},
		},
		{
			name:       "only gets unstakedstaked validators",
			validators: []types.Validator{unstakingValidator, secondaryStakedValidator},
			expected:   expected{length: 1, stakedValidators: true},
			args:       args{stakedValidator: unstakingValidator, validators: []types.Validator{unstakingValidator}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.args.validators {
				keeper.SetValidator(context, validator)
				if validator.IsUnstaking() {
					keeper.SetUnstakingValidator(context, validator)
				}
				if validator.IsStaked() {
					keeper.SetStakedValidator(context, validator)
				}
			}
			validators := keeper.getAllUnstakingValidators(context)
			if !test.expected.stakedValidators {
				assert.Contains(t, validators, test.expected.validators[0])
			}
			assert.Equalf(t, test.expected.length, len(validators), "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestDeleteUnstakingValidator(t *testing.T) {
	stakedValidator := getStakedValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		stakedVal       types.Validator
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
			args:     args{validators: []types.Validator{stakedValidator}},
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
	stakedValidator := getStakedValidator()
	secondaryStakedValidator := getStakedValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		stakedVal       types.Validator
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
			args:     args{validators: []types.Validator{stakedValidator, secondaryStakedValidator}},
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
	stakingValidator := getUnstakingValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		stakedVal       types.Validator
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
			args:     args{validators: []types.Validator{stakingValidator}},
			expected: expected{validators: []types.Validator{stakingValidator}, length: 1, stakedValidators: false},
		},
		{
			name:     "gets empty slice if no mature validators",
			args:     args{validators: []types.Validator{}},
			expected: expected{validators: []types.Validator{stakingValidator}, length: 0, stakedValidators: false},
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
	stakingValidator := getUnstakingValidator()

	type expected struct {
		validators       []types.Validator
		stakedValidators bool
		length           int
	}
	type args struct {
		stakedVal       types.Validator
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
			args:     args{validators: []types.Validator{stakingValidator}},
			expected: expected{validators: []types.Validator{stakingValidator}, length: 0, stakedValidators: false},
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
	stakedValidator := getStakedValidator()
	unstakedValidator := getUnstakedValidator()

	tests := []struct {
		name       string
		validators []types.Validator
		panics     bool
		amount     sdk.Int
	}{
		{
			name:       "recieves a valid iterator",
			validators: []types.Validator{stakedValidator, unstakedValidator},
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
