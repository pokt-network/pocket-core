package keeper

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
)

func TestGetAndSetValidatorBurn(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		amount    sdk.Int
		validator types.Validator
	}
	type expected struct {
		amount sdk.Int
		found  bool
	}
	tests := []struct {
		name string
		args
		expected
	}{
		{
			name:     "can get and set validator burn",
			args:     args{amount: sdk.NewInt(10), validator: stakedValidator},
			expected: expected{amount: sdk.NewInt(10), found: true},
		},
		{
			name:     "returns no coins if not set",
			args:     args{amount: sdk.NewInt(10), validator: stakedValidator},
			expected: expected{amount: sdk.NewInt(0), found: false},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			if test.expected.found {
				keeper.setValidatorBurn(context, test.args.amount, test.args.validator.Address)
			}
			coins, found := keeper.getValidatorBurn(context, test.args.validator.Address)
			assert.Equal(t, test.expected.found, found, "found does not match expected")
			if test.expected.found {
				assert.True(t, test.expected.amount.Equal(coins), "received coins are not the expected coins")
			}
		})
	}
}

func TestDeleteValidatorBurn(t *testing.T) {
	stakedValidator := getStakedValidator()
	var emptyCoins sdk.Int

	type args struct {
		amount    sdk.Int
		validator types.Validator
	}
	type expected struct {
		amount  sdk.Int
		found   bool
		message string
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "deletes validator burn",
			panics:   false,
			args:     args{amount: sdk.NewInt(10), validator: stakedValidator},
			expected: expected{amount: emptyCoins, found: false},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.setValidatorBurn(context, test.args.amount, test.args.validator.Address)
			keeper.deleteValidatorBurn(context, test.args.validator.Address)
			coins, found := keeper.getValidatorBurn(context, test.args.validator.Address)
			assert.Equal(t, test.expected.found, found, "found does not match expected")
			if found {
				assert.True(t, coins.IsZero(), "received coins are not the expected coins")
			}
		})
	}
}

func TestHandleValidatorSignature(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator        types.Validator
		power            int64
		signed           bool
		increasedContext int64
		maxMissed        int64
	}
	type expected struct {
		validator           types.Validator
		tombstoned          bool
		missedBlocksCounter int64
		message             string
		signedInfo          bool
		jail                bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:     "handles a signature",
			panics:   false,
			args:     args{validator: stakedValidator, power: int64(10), signed: false},
			expected: expected{validator: stakedValidator, tombstoned: false, missedBlocksCounter: int64(1)},
		},
		{
			name:     "previously signed signature",
			panics:   false,
			args:     args{validator: stakedValidator, power: int64(10), signed: true},
			expected: expected{validator: stakedValidator, tombstoned: false, missedBlocksCounter: int64(0)},
		},
		{
			name:     "jails if signature with overflown minHeight and maxHeight",
			panics:   false,
			args:     args{validator: stakedValidator, power: int64(10), signed: true, increasedContext: 51, maxMissed: 51},
			expected: expected{validator: stakedValidator, tombstoned: false, missedBlocksCounter: int64(0)},
		},
		{
			name:   "errors if no signed info",
			panics: true,
			args:   args{validator: stakedValidator, power: int64(10), signed: false},
			expected: expected{
				message: fmt.Sprintf("Expected signing info for validator %s but not found", sdk.Address(stakedValidator.GetPublicKey().Address())),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.validator.GetPublicKey().Address()
			keeper.SetValidator(context, test.args.validator)
			switch test.panics {
			case true:
				defer func() {
					err := recover()
					assert.Contains(t, test.expected.message, err, "does not contain error ")
				}()
				if test.expected.signedInfo {
					keeper.handleValidatorSignature(context, cryptoAddr, test.args.power, test.args.signed)
				}
				keeper.handleValidatorSignature(context, cryptoAddr, test.args.power, test.args.signed)
			default:
				signingInfo := types.ValidatorSigningInfo{
					Address:     test.args.validator.GetAddress(),
					StartHeight: context.BlockHeight(),
					JailedUntil: time.Unix(0, 0),
				}
				if test.expected.jail {
					context.WithBlockHeight(101)
					signingInfo.MissedBlocksCounter = test.args.maxMissed
				}
				keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
				keeper.handleValidatorSignature(context, cryptoAddr, test.args.power, test.args.signed)
				signedInfo, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
				if !found {
					t.FailNow()
				}
				assert.Equal(t, test.expected.tombstoned, signedInfo.Tombstoned)
				assert.Equal(t, test.expected.missedBlocksCounter, signedInfo.MissedBlocksCounter)
				if test.expected.jail {
					validator, found := keeper.GetValidator(context, sdk.Address(cryptoAddr))
					if !found {
						t.FailNow()
					}
					assert.True(t, validator.Jailed)
				}

			}
		})
	}
}

func TestValidateDoubleSign(t *testing.T) {
	stakedValidator := getStakedValidator()

	type args struct {
		validator        types.Validator
		increasedContext int64
		maxMissed        int64
	}
	type expected struct {
		validator      types.Validator
		tombstoned     bool
		message        string
		pubKeyRelation bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:   "handles double signature",
			panics: false,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: true,
				tombstoned:     false,
			},
		},
		{
			name:   "ignores double signature on tombstoned validator",
			panics: true,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: true,
				tombstoned:     true,
				message:        fmt.Sprintf("ERROR:\nCodespace: pos\nCode: 113\nMessage: \"Warning: validator is already tombstoned\"\n"),
			},
		},
		{
			name:   "ignores double signature on tombstoned validator",
			panics: true,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: false,
				tombstoned:     false,
				message:        fmt.Sprintf("ERROR:\nCodespace: pos\nCode: 114\nMessage: \"Warning: the DS evidence is unable to be handled\"\n"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.validator.GetAddress()
			keeper.SetValidator(context, test.args.validator)
			signingInfo := types.ValidatorSigningInfo{
				Address:     test.args.validator.GetAddress(),
				StartHeight: context.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}
			if test.expected.tombstoned {
				signingInfo.Tombstoned = test.expected.tombstoned
			}
			infractionHeight := context.BlockHeight()
			keeper.SetValidatorSigningInfo(context, cryptoAddr, signingInfo)
			signingInfo, found := keeper.GetValidatorSigningInfo(context, cryptoAddr)
			if !found {
				t.FailNow()
			}
			consAddr, signedInfo, validator, err := keeper.validateDoubleSign(context, crypto.Address(cryptoAddr), infractionHeight, time.Unix(0, 0))
			if err != nil {
				assert.Equal(t, test.expected.message, err.Error())
			} else {
				assert.Equal(t, cryptoAddr, consAddr, "addresses do not match")
				assert.Equal(t, signedInfo, signingInfo, "signed Info do not match")
				assert.Equal(t, test.expected.validator, validator, "validators do not match")
			}
		})
	}
}

func TestHandleDoubleSign(t *testing.T) {
	stakedValidator := getStakedValidator()
	supplySize := sdk.NewInt(100)

	type args struct {
		validator        types.Validator
		increasedContext int64
		maxMissed        int64
		power            int64
	}
	type expected struct {
		validator      types.Validator
		tombstoned     bool
		message        string
		found          bool
		pubKeyRelation bool
	}
	tests := []struct {
		name string
		args
		expected
	}{
		{
			name: "handles double signature",
			args: args{validator: stakedValidator, power: int64(10)},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: true,
				found:          true,
				tombstoned:     false,
			},
		},
		{
			name: "ignores double signature on tombstoned validator",
			args: args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: true,
				tombstoned:     true,
				found:          true,
				message:        fmt.Sprintf("ERROR:\nCodespace: pos\nCode: 113\nMessage: \"Warning: validator is already tombstoned\"\n"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.validator.GetPublicKey().Address()
			keeper.SetValidator(context, test.args.validator)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.validator.Address, supplySize)
			signingInfo := types.ValidatorSigningInfo{
				Address:     test.args.validator.GetAddress(),
				StartHeight: context.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}

			if test.expected.tombstoned {
				signingInfo.Tombstoned = test.expected.tombstoned
			}
			infractionHeight := context.BlockHeight()
			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			keeper.handleDoubleSign(context, cryptoAddr, infractionHeight, time.Unix(0, 0), test.args.power)

			signingInfo, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if found != test.expected.found {
				t.FailNow()
			}

			assert.Equal(t, test.expected.tombstoned, signingInfo.Tombstoned)
		})
	}
}

func TestValidateSlash(t *testing.T) {
	stakedValidator := getStakedValidator()
	unstakedValidator := getUnstakedValidator()
	supplySize := sdk.NewInt(100)

	type args struct {
		validator        types.Validator
		power            int64
		increasedContext int64
		slashFraction    sdk.Dec
		maxMissed        int64
	}
	type expected struct {
		validator      types.Validator
		tombstoned     bool
		message        string
		pubKeyRelation bool
		fraction       bool
		customHeight   bool
		found          bool
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:   "validates slash",
			panics: false,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
			},
		},
		{
			name:   "empty validator if not found",
			panics: false,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
			},
		},
		{
			name:   "errors if unstakedValidator",
			panics: false,
			args:   args{validator: unstakedValidator},
			expected: expected{
				validator:      unstakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				fraction:       false,
				message:        fmt.Sprintf("should not be slashing unstaked validator: %s", unstakedValidator.Address),
			},
		},
		{
			name:   "errors with invalid slashFactor",
			panics: true,
			args:   args{validator: unstakedValidator, slashFraction: sdk.NewDec(-10)},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				fraction:       true,
				message:        fmt.Sprintf("attempted to slash with a negative slash factor: %v", sdk.NewDec(-10)),
			},
		},
		{
			name:   "errors with wrong infraction height",
			panics: true,
			args:   args{validator: unstakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				fraction:       false,
				customHeight:   true,
				message:        fmt.Sprintf("impossible attempt to slash future infraction at height %d but we are at height %d", 100, 0),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.validator.GetPublicKey().Address()
			if test.expected.found {
				keeper.SetValidator(context, test.args.validator)
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.validator.Address, supplySize)
			}
			signingInfo := types.ValidatorSigningInfo{
				Address:     test.args.validator.GetAddress(),
				StartHeight: context.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}
			if test.expected.tombstoned {
				signingInfo.Tombstoned = test.expected.tombstoned
			}
			infractionHeight := context.BlockHeight()

			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			signingInfo, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if !found {
				t.FailNow()
			}
			var fraction sdk.Dec
			if test.expected.fraction {
				fraction = test.args.slashFraction
			} else {
				fraction = keeper.SlashFractionDoubleSign(context)
			}

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Equal(t, test.expected.message, err.Error(), "message error does not match")
				}()
				if test.expected.customHeight {
					updatedContext := context.WithBlockHeight(100)
					infractionHeight = updatedContext.BlockHeight()
				}
				_ = keeper.validateSlash(context, sdk.Address(cryptoAddr), infractionHeight, test.args.power, fraction)
			default:
				val := keeper.validateSlash(context, sdk.Address(cryptoAddr), infractionHeight, test.args.power, fraction)
				if test.expected.found {
					assert.Equal(t, test.expected.validator, val)
				} else {
					assert.Equal(t, types.Validator{}, val)
				}
			}
		})
	}
}

func TestSlash(t *testing.T) {
	stakedValidator := getStakedValidator()
	supplySize := sdk.NewInt(50001)

	type args struct {
		validator        types.Validator
		power            int64
		increasedContext int64
		slashFraction    sdk.Dec
		maxMissed        int64
	}
	type expected struct {
		validator      types.Validator
		tombstoned     bool
		message        string
		pubKeyRelation bool
		fraction       bool
		customHeight   bool
		found          bool
		stakedTokens   sdk.Int
	}
	tests := []struct {
		name   string
		panics bool
		args
		expected
	}{
		{
			name:   "slash all validaor coins",
			panics: false,
			args:   args{validator: stakedValidator, power: int64(1)},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				stakedTokens:   stakedValidator.StakedTokens.Sub(sdk.NewInt(50000)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.validator.GetPublicKey().Address()
			if test.expected.found {
				keeper.SetValidator(context, test.args.validator)
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.validator.Address, supplySize)
				v, found := keeper.GetValidator(context, sdk.Address(cryptoAddr))
				if !found {
					t.FailNow()
				}

				fmt.Println(v)
			}
			signingInfo := types.ValidatorSigningInfo{
				Address:     test.args.validator.GetAddress(),
				StartHeight: context.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}
			if test.expected.tombstoned {
				signingInfo.Tombstoned = test.expected.tombstoned
			}
			infractionHeight := context.BlockHeight()

			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			signingInfo, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if !found {
				t.FailNow()
			}
			var fraction sdk.Dec
			if test.expected.fraction {
				fraction = test.args.slashFraction
			} else {
				fraction = keeper.SlashFractionDoubleSign(context)
			}

			keeper.slash(context, sdk.Address(cryptoAddr), infractionHeight, test.args.power, fraction)
			validator, found := keeper.GetValidator(context, sdk.Address(cryptoAddr))
			if !found {
				t.Fail()
			}
			assert.True(t, validator.StakedTokens.Equal(test.expected.stakedTokens), "tokens were not slashed")
		})
	}
}

func TestBurnValidators(t *testing.T) {
	primaryStakedValidator := getStakedValidator()

	type args struct {
		amount    sdk.Int
		validator types.Validator
	}
	type expected struct {
		amount    sdk.Int
		found     bool
		validator types.Validator
	}
	tests := []struct {
		name string
		args
		expected
	}{
		{
			name: "can get and set validator burn",
			args: args{
				amount:    sdk.NewInt(100000000000),
				validator: primaryStakedValidator,
			},
			expected: expected{
				amount:    sdk.ZeroInt(),
				found:     true,
				validator: primaryStakedValidator,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidator(context, test.args.validator)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.validator.Address, test.args.validator.StakedTokens)
			keeper.setValidatorBurn(context, test.args.amount, test.args.validator.Address)
			keeper.burnValidators(context)
			primaryCryptoAddr := test.args.validator.GetAddress()
			primaryValidator, found := keeper.GetValidator(context, primaryCryptoAddr)
			if !found {
				t.Fail()
			}
			fmt.Println(test.expected.amount)
			fmt.Println(primaryValidator.StakedTokens)
			assert.True(t, test.expected.amount.Equal(primaryValidator.StakedTokens))
		})
	}
}

func TestKeeper_getBurnFromSeverity(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}

	context, _, keeper := createTestInput(t, true)
	primaryStakedValidator := getStakedValidator()

	type args struct {
		ctx                sdk.Context
		address            sdk.Address
		severityPercentage sdk.Dec
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Int
	}{
		{"Test getBurnFromSeverity", fields{keeper},
			args{
				ctx:                context,
				address:            primaryStakedValidator.Address,
				severityPercentage: sdk.OneDec(),
			}, sdk.NewInt(100000000000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper
			k.SetValidator(context, primaryStakedValidator)

			if got := k.getBurnFromSeverity(tt.args.ctx, tt.args.address, tt.args.severityPercentage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBurnFromSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_BurnValidator(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}

	primaryStakedValidator := getStakedValidator()

	context, _, keeper := createTestInput(t, true)

	type args struct {
		ctx     sdk.Context
		address sdk.Address
		amount  sdk.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test BurnValidator", fields{Keeper: keeper},
			args{
				ctx:     context,
				address: primaryStakedValidator.Address,
				amount:  sdk.ZeroInt(),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper
			k.SetValidator(context, primaryStakedValidator)
			store := tt.args.ctx.KVStore(k.storeKey)
			store.Set(types.KeyForValidatorBurn(tt.args.address), k.cdc.MustMarshalBinaryBare(sdk.NewDec(1)))
			k.BurnValidator(tt.args.ctx, tt.args.address, tt.args.amount)
			burn, found := keeper.getValidatorBurn(tt.args.ctx, tt.args.address)
			assert.True(t, found)
			assert.NotNil(t, burn)
		})
	}
}
