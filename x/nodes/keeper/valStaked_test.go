package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetStakedValidator(t *testing.T) {
	stakedValidator := getStakedValidator()
	unstakedValidator := getUnstakedValidator()

	type expected struct {
		validators []types.Validator
		length     int
	}
	tests := []struct {
		name       string
		validator  types.Validator
		validators []types.Validator
		expected   expected
	}{
		{
			name:       "gets validators",
			validators: []types.Validator{stakedValidator},
			expected:   expected{validators: []types.Validator{stakedValidator}, length: 1},
		},
		{
			name:       "gets emtpy slice of validators",
			validators: []types.Validator{unstakedValidator},
			expected:   expected{validators: []types.Validator{}, length: 0},
		},
		{
			name:       "only gets staked validators",
			validators: []types.Validator{stakedValidator, unstakedValidator},
			expected:   expected{validators: []types.Validator{stakedValidator}, length: 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.validators {
				keeper.SetValidator(context, validator)
			}
			validators := keeper.getStakedValidators(context)

			if equal := assert.ObjectsAreEqualValues(validators, test.expected.validators); !equal { // note ObjectsAreEqualValues does not assert, manual verification is required
				t.FailNow()
			}
			assert.Equalf(t, len(validators), test.expected.length, "length of the validators does not match expected on %v", test.name)
		})
	}
}

func TestGetSetDeleteValidatorsByChain(t *testing.T) {
	stakedValidator := getStakedValidator()
	context, _, keeper := createTestInput(t, true)
	keeper.SetValidator(context, stakedValidator)
	keeper.SetStakedValidatorByChains(context, stakedValidator)
	vals := keeper.GetValidatorsByChain(context, stakedValidator.Chains[0])
	assert.Contains(t, vals, stakedValidator)
	vals = keeper.GetValidatorsByChain(context, stakedValidator.Chains[1])
	assert.Contains(t, vals, stakedValidator)
	vals = keeper.GetValidatorsByChain(context, stakedValidator.Chains[2])
	assert.Contains(t, vals, stakedValidator)
	keeper.deleteValidatorForChains(context, stakedValidator)
	vals = keeper.GetValidatorsByChain(context, stakedValidator.Chains[0])
	assert.NotContains(t, vals, stakedValidator)
	vals = keeper.GetValidatorsByChain(context, stakedValidator.Chains[1])
	assert.NotContains(t, vals, stakedValidator)
	vals = keeper.GetValidatorsByChain(context, stakedValidator.Chains[2])
	assert.NotContains(t, vals, stakedValidator)
}

func TestRemoveStakedValidatorTokens(t *testing.T) {
	stakedValidator := getStakedValidator()

	type expected struct {
		tokens       sdk.Int
		validators   []types.Validator
		errorMessage string
	}
	tests := []struct {
		name      string
		validator types.Validator
		hasError  bool
		amount    sdk.Int
		expected
	}{
		{
			name:      "removes tokens from validator validators",
			validator: stakedValidator,
			amount:    sdk.NewInt(5),
			hasError:  false,
			expected:  expected{tokens: sdk.NewInt(99999999995), validators: []types.Validator{}},
		},
		{
			name:      "removes tokens from validator validators",
			validator: stakedValidator,
			amount:    sdk.NewInt(-5),
			hasError:  true,
			expected:  expected{tokens: sdk.NewInt(99999999995), validators: []types.Validator{}, errorMessage: "trying to remove negative tokens"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.validator)
			switch test.hasError {
			case true:
				_, _ = keeper.removeValidatorTokens(context, test.validator, test.amount)
			default:
				validator, _ := keeper.removeValidatorTokens(context, test.validator, test.amount)
				assert.True(t, validator.StakedTokens.Equal(test.expected.tokens), "validator staked tokens is not as expected")

				store := context.KVStore(keeper.storeKey)
				assert.NotNil(t, store.Get(types.KeyForValidatorInStakingSet(validator)))
			}
		})
	}
}

func TestRemoveDeleteFromStakingSet(t *testing.T) {
	stakedValidator := getStakedValidator()
	unstakedValidator := getUnstakedValidator()

	tests := []struct {
		name       string
		validators []types.Validator
		panics     bool
		amount     sdk.Int
	}{
		{
			name:       "removes validators from set",
			validators: []types.Validator{stakedValidator, unstakedValidator},
			panics:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, validator := range test.validators {
				keeper.SetValidator(context, validator)
			}
			for _, validator := range test.validators {
				keeper.deleteValidatorFromStakingSet(context, validator)
			}

			validators := keeper.getStakedValidators(context)
			assert.Empty(t, validators, "there should not be any validators in the set")
		})
	}
}

func TestGetValsIterator(t *testing.T) {
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
			}

			it := keeper.stakedValsIterator(context)
			assert.Implements(t, (*sdk.Iterator)(nil), it, "does not implement interface")
		})
	}
}
func TestApplicationStaked_IterateAndExecuteOverStakedApps(t *testing.T) {
	stakedValidator := getStakedValidator()
	secondStakedValidator := getStakedValidator()
	tests := []struct {
		name         string
		application  types.Validator
		applications []types.Validator
		want         int
	}{
		{
			name:         "iterates over applications",
			applications: []types.Validator{stakedValidator, secondStakedValidator},
			want:         2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			for _, application := range tt.applications {
				keeper.SetValidator(context, application)
			}
			got := 0
			fn := modifyFn(&got)
			keeper.IterateAndExecuteOverStakedVals(context, fn)
			if got != tt.want {
				t.Errorf("appStaked.IterateAndExecuteOverApps() = got %v, want %v", got, tt.want)
			}
		})
	}
}
