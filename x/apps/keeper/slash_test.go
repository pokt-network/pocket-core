package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAndSetApplicationBurn(t *testing.T) {
	stakedApplication := getStakedApplication()

	type args struct {
		amount      sdk.Int
		application types.Application
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
			name:     "can get and set application burn",
			args:     args{amount: sdk.NewInt(10), application: stakedApplication},
			expected: expected{amount: sdk.NewInt(10), found: true},
		},
		{
			name:     "returns no coins if not set",
			args:     args{amount: sdk.NewInt(10), application: stakedApplication},
			expected: expected{amount: sdk.NewInt(0), found: false},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			if test.expected.found {
				keeper.setApplicationBurn(context, test.args.amount, test.args.application.Address)
			}
			coins, found := keeper.getApplicationBurn(context, test.args.application.Address)
			assert.Equal(t, test.expected.found, found, "found does not match expected")
			if test.expected.found {
				assert.True(t, test.expected.amount.Equal(coins), "received coins are not the expected coins")
			}
		})
	}
}

func TestDeleteApplicationBurn(t *testing.T) {
	stakedApplication := getStakedApplication()
	var emptyCoins sdk.Int

	type args struct {
		amount      sdk.Int
		application types.Application
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
			name:     "deletes application burn",
			panics:   false,
			args:     args{amount: sdk.NewInt(10), application: stakedApplication},
			expected: expected{amount: emptyCoins, found: false},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.setApplicationBurn(context, test.args.amount, test.args.application.Address)
			keeper.deleteApplicationBurn(context, test.args.application.Address)
			_, found := keeper.getApplicationBurn(context, test.args.application.Address)
			assert.Equal(t, test.expected.found, found, "found does not match expected")
		})
	}
}

func TestValidateSlash(t *testing.T) {
	stakedApplication := getStakedApplication()
	unstakedApplication := getUnstakedApplication()
	supplySize := sdk.NewInt(100)

	type args struct {
		application      types.Application
		power            int64
		increasedContext int64
		slashFraction    sdk.Int
		maxMissed        int64
	}
	type expected struct {
		application    types.Application
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
			name:   "validates simpleSlash",
			panics: false,
			args:   args{application: stakedApplication, slashFraction: sdk.NewInt(90)},
			expected: expected{
				application:    stakedApplication,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
			},
		},
		{
			name:   "empty application if not found",
			panics: false,
			args:   args{application: stakedApplication, slashFraction: sdk.NewInt(90)},
			expected: expected{
				application:    stakedApplication,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
			},
		},
		{
			name:   "errors if unstakedApplication",
			panics: true,
			args:   args{application: unstakedApplication, slashFraction: sdk.NewInt(90)},
			expected: expected{
				application:    stakedApplication,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				fraction:       false,
				message:        fmt.Sprintf("should not be slashing unstaked application: %s", unstakedApplication.Address),
			},
		},
		{
			name:   "errors with invalid slashFactor",
			panics: true,
			args:   args{application: unstakedApplication, slashFraction: sdk.NewInt(-10)},
			expected: expected{
				application:    stakedApplication,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				fraction:       true,
				message:        fmt.Sprintf("attempted to simpleSlash with a negative simpleSlash factor: %v", sdk.NewInt(-10)),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.application.GetPublicKey().Address()
			if test.expected.found {
				keeper.SetApplication(context, test.args.application)
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.application.Address, supplySize)
			}
			fraction := test.args.slashFraction

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Equal(t, test.expected.message, err.Error(), "message error does not match")
				}()
				_ = keeper.validateSimpleSlash(context, sdk.Address(cryptoAddr), fraction)
			default:
				val := keeper.validateSimpleSlash(context, sdk.Address(cryptoAddr), fraction)
				if test.expected.found {
					assert.Equal(t, test.expected.application, val)
				} else {
					assert.Equal(t, types.Application{}, val)
				}
			}
		})
	}
}

func TestSlash(t *testing.T) {
	stakedApplication := getStakedApplication()
	supplySize := sdk.NewInt(50001)

	type args struct {
		application      types.Application
		power            int64
		increasedContext int64
		slashFraction    sdk.Int
		maxMissed        int64
	}
	type expected struct {
		application    types.Application
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
			name:   "simpleSlash application coins",
			panics: false,
			args:   args{application: stakedApplication, power: int64(1), slashFraction: sdk.NewInt(100000000000)},
			expected: expected{
				application:    stakedApplication,
				found:          true,
				pubKeyRelation: true,
				tombstoned:     false,
				stakedTokens:   sdk.ZeroInt(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			cryptoAddr := test.args.application.GetPublicKey().Address()
			if test.expected.found {
				keeper.SetApplication(context, test.args.application)
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.application.Address, supplySize)
			}
			fraction := test.args.slashFraction
			keeper.simpleSlash(context, sdk.Address(cryptoAddr), fraction)
			application, found := keeper.GetApplication(context, sdk.Address(cryptoAddr))
			if !found {
				t.Fail()
			}
			assert.True(t, application.StakedTokens.Equal(test.expected.stakedTokens), "tokens were not slashed")
		})
	}
}

func TestBurnApplications(t *testing.T) {
	primaryStakedApplication := getStakedApplication()

	type args struct {
		amount      sdk.Int
		application types.Application
	}
	type expected struct {
		amount      sdk.Int
		found       bool
		application types.Application
	}
	tests := []struct {
		name string
		args
		expected
	}{
		{
			name: "can get and set application burn",
			args: args{
				amount:      sdk.NewInt(1000000000000),
				application: primaryStakedApplication,
			},
			expected: expected{
				amount:      sdk.ZeroInt(),
				found:       true,
				application: primaryStakedApplication,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, test.args.application)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.args.application.Address, test.args.application.StakedTokens)
			keeper.setApplicationBurn(context, test.args.amount, test.args.application.Address)
			keeper.burnApplications(context)

			primaryCryptoAddr := test.args.application.GetAddress()

			primaryApplication, found := keeper.GetApplication(context, primaryCryptoAddr)
			if !found {
				t.Fail()
			}
			assert.True(t, test.expected.amount.Equal(primaryApplication.StakedTokens))
		})
	}
}
