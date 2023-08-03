package keeper

import (
	"fmt"
	"reflect"
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
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

// There are three versions of structs to represent a validator.
//   - LegacyValidator - the original version
//   - LegacyValidator8 - LegacyValidator + OutputAddress (introduced in 0.8)
//   - Validator - LegacyValidator8 + Delegators (introduced in 0.11)
//
// The following two tests verify marshaling/unmarshaling has backward/forward
// compatibility, meaning marshaled bytes can be unmarshaled as a newer version
// or an older version.
func TestValidator_Amino_MarshalingCompatibility(t *testing.T) {
	_, _, k := createTestInput(t, false)
	Marshal := k.Cdc.AminoCodec().MarshalBinaryLengthPrefixed
	Unmarshal := k.Cdc.AminoCodec().UnmarshalBinaryLengthPrefixed

	// Amino cannot handle type.Validator because map is not supported.
	// We don't have to test type.Validator because it didn't exist while
	// we were using Amino (before UpgradeCodecHeight).
	var (
		val_1, val_2   types.LegacyValidator8
		valL_1, valL_2 types.LegacyValidator
		marshaled      []byte
		err            error
	)

	val_1 = getStakedValidator().ToLegacy8()
	val_1.OutputAddress = getRandomValidatorAddress()
	valL_1 = val_1.ToLegacy()

	// Validator --> []byte --> Validator
	marshaled, err = Marshal(&val_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val_2.Reset()
	err = Unmarshal(marshaled, &val_2)
	assert.Nil(t, err)
	assert.True(t, val_2.ToLegacy().ExactEqualsTo(val_1.ToLegacy()))
	assert.True(t, val_2.OutputAddress.Equals(val_1.OutputAddress))

	// Validator --> []byte --> LegacyValidator
	marshaled, err = Marshal(&val_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	valL_2.Reset()
	err = Unmarshal(marshaled, &valL_2)
	assert.Nil(t, err)
	assert.True(t, valL_2.ExactEqualsTo(val_1.ToLegacy()))

	// LegacyValidator --> []byte --> Validator
	marshaled, err = Marshal(&valL_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val_2.Reset()
	err = Unmarshal(marshaled, &val_2)
	assert.Nil(t, err)
	assert.True(t, val_2.ToLegacy().ExactEqualsTo(valL_1))
	assert.Nil(t, val_2.OutputAddress)
}

func TestValidator_Proto_MarshalingCompatibility(t *testing.T) {
	_, _, k := createTestInput(t, false)
	Marshal := k.Cdc.ProtoCodec().MarshalBinaryLengthPrefixed
	Unmarshal := k.Cdc.ProtoCodec().UnmarshalBinaryLengthPrefixed

	var (
		val_1, val_2   types.Validator
		val8_1, val8_2 types.LegacyValidator8
		valL_1, valL_2 types.LegacyValidator
		marshaled      []byte
		err            error
	)

	val_1 = getStakedValidator()
	val_1.OutputAddress = getRandomValidatorAddress()
	val_1.Delegators = map[string]uint32{}
	val_1.Delegators[getRandomValidatorAddress().String()] = 10
	val_1.Delegators[getRandomValidatorAddress().String()] = 20
	val8_1 = val_1.ToLegacy8()
	valL_1 = val_1.ToLegacy()

	// Validator --> []byte --> Validator
	marshaled, err = Marshal(&val_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val_2.Reset()
	err = Unmarshal(marshaled, &val_2)
	assert.Nil(t, err)
	assert.True(t, val_2.ToLegacy().ExactEqualsTo(val_1.ToLegacy()))
	assert.True(t, val_2.OutputAddress.Equals(val_1.OutputAddress))
	assert.True(t, types.CompareStringMaps(val_2.Delegators, val_1.Delegators))

	// Validator --> []byte --> LegacyValidator8
	marshaled, err = Marshal(&val_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val8_2.Reset()
	err = Unmarshal(marshaled, &val8_2)
	assert.Nil(t, err)
	assert.True(t, val8_2.ToLegacy().ExactEqualsTo(val_1.ToLegacy()))
	assert.True(t, val8_2.OutputAddress.Equals(val_1.OutputAddress))

	// Validator --> []byte --> LegacyValidator
	marshaled, err = Marshal(&val_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	valL_2.Reset()
	err = Unmarshal(marshaled, &valL_2)
	assert.Nil(t, err)
	assert.True(t, valL_2.ExactEqualsTo(val_1.ToLegacy()))

	// LegacyValidator8 --> []byte --> Validator
	marshaled, err = Marshal(&val8_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val_2.Reset()
	err = Unmarshal(marshaled, &val_2)
	assert.Nil(t, err)
	assert.True(t, val_2.ToLegacy().ExactEqualsTo(val8_1.ToLegacy()))
	assert.True(t, val_2.OutputAddress.Equals(val8_1.OutputAddress))
	assert.Nil(t, val_2.Delegators)

	// LegacyValidator --> []byte --> Validator
	marshaled, err = Marshal(&valL_1)
	assert.Nil(t, err)
	assert.NotNil(t, marshaled)
	val_2.Reset()
	err = Unmarshal(marshaled, &val_2)
	assert.Nil(t, err)
	assert.True(t, val_2.ToLegacy().ExactEqualsTo(valL_1))
	assert.Nil(t, val_2.OutputAddress)
	assert.Nil(t, val_2.Delegators)
}
