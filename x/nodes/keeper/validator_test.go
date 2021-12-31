package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestKeeper_GetValidators(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx         sdk.Context
		maxRetrieve uint16
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantValidators []types.Validator
	}{
		{"Test GetValidators 0", fields{keeper: keeper}, args{
			ctx:         context,
			maxRetrieve: 0,
		}, []types.Validator{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper

			if gotValidators := k.GetValidators(tt.args.ctx, tt.args.maxRetrieve); !reflect.DeepEqual(gotValidators, tt.wantValidators) {
				t.Errorf("GetValidators() = %v, want %v", gotValidators, tt.wantValidators)
			}
		})
	}
}

func TestKeeper_GetValidatorOutputAddress(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		v   types.Validator
	}
	validator := getStakedValidator()
	validator.OutputAddress = validator.Address
	validatorNoOuptut := getStakedValidator()
	validatorNoOuptut.OutputAddress = nil
	context, _, keeper := createTestInput(t, true)
	keeper.SetValidator(context, validator)
	keeper.SetValidator(context, validatorNoOuptut)
	tests := []struct {
		name string
		args args
		want sdk.Address
	}{
		{"Test GetValidatorOutput With Output Address", args{
			ctx: context,
			k:   keeper,
			v:   validator,
		}, validator.OutputAddress},
		{"Test GetValidatorOutput Without Output Address", args{
			ctx: context,
			k:   keeper,
			v:   validatorNoOuptut,
		}, validatorNoOuptut.Address},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := tt.args.k.GetValidatorOutputAddress(tt.args.ctx, tt.args.v.Address)
			if !assert.True(t, len(got) == len(tt.want)) {
				t.Errorf("GetValidatorOutputAddress() = %v, want %v", got, tt.want)
			}
			assert.True(t, found)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
				validator, _ := keeper.GetValidator(context, test.args.validator.Address)
				assert.True(t, validator.Equals(test.expected.validator), "validator does not match")
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
