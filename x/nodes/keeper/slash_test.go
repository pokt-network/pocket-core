package keeper

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
)

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
		missedBlocksCounter int64
		message             string
		jail                bool
	}
	tests := []struct {
		name     string
		hasError bool
		args
		expected
	}{
		{
			name:     "handles a signature",
			hasError: false,
			args:     args{validator: stakedValidator, power: int64(10), signed: false},
			expected: expected{validator: stakedValidator, missedBlocksCounter: int64(1)},
		},
		{
			name:     "previously signed signature",
			hasError: false,
			args:     args{validator: stakedValidator, power: int64(10), signed: true},
			expected: expected{validator: stakedValidator, missedBlocksCounter: int64(0)},
		},
		{
			name:     "jails if signature with overflown minHeight and maxHeight",
			hasError: false,
			args:     args{validator: stakedValidator, power: int64(10), signed: true, increasedContext: 51, maxMissed: 51},
			expected: expected{validator: stakedValidator, missedBlocksCounter: int64(0)},
		},
		{
			name:     "errors if no signed info",
			hasError: true,
			args:     args{validator: stakedValidator, power: int64(10), signed: false},
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
			switch test.hasError {
			case true:
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
				signedBlocksWindow := keeper.SignedBlocksWindow(context)
				minSignedPerWindow := keeper.MinSignedPerWindow(context)
				downtimeJailDuration := keeper.DowntimeJailDuration(context)
				slashFractionDowntime := keeper.SlashFractionDowntime(context)
				keeper.handleValidatorSignature(context, sdk.Address(cryptoAddr), test.args.power, test.args.signed, signedBlocksWindow, minSignedPerWindow, downtimeJailDuration, slashFractionDowntime)
				signedInfo, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
				if !found {
					t.FailNow()
				}
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
		validator types.Validator
	}
	type expected struct {
		validator      types.Validator
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
			},
		},
		{
			name:   "ignores double signature on tombstoned validator",
			panics: true,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: true,
				message:        "ERROR:\nCodespace: pos\nCode: 113\nMessage: \"Warning: validator is already tombstoned\"\n",
			},
		},
		{
			name:   "ignores double signature on tombstoned validator",
			panics: true,
			args:   args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				pubKeyRelation: false,
				message:        "ERROR:\nCodespace: pos\nCode: 114\nMessage: \"Warning: the DS evidence is unable to be handled\"\n",
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
		validator types.Validator
		power     int64
	}
	type expected struct {
		validator      types.Validator
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
			infractionHeight := context.BlockHeight()
			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			keeper.handleDoubleSign(context, cryptoAddr, infractionHeight, time.Unix(0, 0), test.args.power)

			_, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if found != test.expected.found {
				t.FailNow()
			}
		})
	}
}

func TestValidateSlash(t *testing.T) {
	stakedValidator := getStakedValidator()
	unstakedValidator := getUnstakedValidator()
	supplySize := sdk.NewInt(100)

	type args struct {
		validator     types.Validator
		power         int64
		slashFraction sdk.BigDec
	}
	type expected struct {
		validator      types.Validator
		message        string
		pubKeyRelation bool
		fraction       bool
		customHeight   bool
		found          bool
	}
	tests := []struct {
		name     string
		hasError bool
		args
		expected
	}{
		{
			name:     "validates slash",
			hasError: false,
			args:     args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
			},
		},
		{
			name:     "empty validator if not found",
			hasError: false,
			args:     args{validator: stakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
			},
		},
		{
			name:     "errors if unstakedValidator",
			hasError: true,
			args:     args{validator: unstakedValidator},
			expected: expected{
				validator:      unstakedValidator,
				found:          true,
				pubKeyRelation: true,
				fraction:       false,
				message:        fmt.Sprintf("should not be slashing unstaked validator: %s", unstakedValidator.Address),
			},
		},
		{
			name:     "errors with invalid slashFactor",
			hasError: true,
			args:     args{validator: unstakedValidator, slashFraction: sdk.NewDec(-10)},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
				fraction:       true,
				message:        fmt.Sprintf("attempted to slash with a negative slash factor: %v", sdk.NewDec(-10)),
			},
		},
		{
			name:     "errors with wrong infraction height",
			hasError: true,
			args:     args{validator: unstakedValidator},
			expected: expected{
				validator:      stakedValidator,
				found:          true,
				pubKeyRelation: true,
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
			infractionHeight := context.BlockHeight()

			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			_, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if !found {
				t.FailNow()
			}
			var fraction sdk.BigDec
			if test.expected.fraction {
				fraction = test.args.slashFraction
			} else {
				fraction = keeper.SlashFractionDoubleSign(context)
			}
			switch test.hasError {
			case true:
				if test.expected.customHeight {
					updatedContext := context.WithBlockHeight(100)
					infractionHeight = updatedContext.BlockHeight()
				}
				val := keeper.validateSlash(context, sdk.Address(cryptoAddr), infractionHeight, test.args.power, fraction)
				assert.Equal(t, types.Validator{}, val)
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
		validator     types.Validator
		power         int64
		slashFraction sdk.BigDec
	}
	type expected struct {
		validator      types.Validator
		pubKeyRelation bool
		fraction       bool
		found          bool
		stakedTokens   sdk.BigInt
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
			infractionHeight := context.BlockHeight()

			keeper.SetValidatorSigningInfo(context, sdk.Address(cryptoAddr), signingInfo)
			_, found := keeper.GetValidatorSigningInfo(context, sdk.Address(cryptoAddr))
			if !found {
				t.FailNow()
			}
			var fraction sdk.BigDec
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
