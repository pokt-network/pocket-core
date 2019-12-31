package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/apps/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCoinsFromUnstakedToStaked(t *testing.T) {
	application := getBondedApplication()
	applicationAddress := application.Address

	tests := []struct {
		name        string
		expected    string
		application types.Application
		amount      sdk.Int
		panics      bool
	}{
		{
			name:        "stake coins on pool",
			application: types.Application{Address: applicationAddress},
			amount:      sdk.NewInt(10),
			panics:      false,
		},
		{
			name:        "panics if negative ammount",
			application: types.Application{Address: applicationAddress},
			amount:      sdk.NewInt(-1),
			expected:    fmt.Sprintf("negative coin amount: -1"),
			panics:      true,
		},
		{name: "panics if no supply is set",
			application: types.Application{Address: applicationAddress},
			expected:    fmt.Sprintf("insufficient account funds"),
			amount:      sdk.NewInt(10),
			panics:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected)
				}()
				if strings.Contains(test.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, sdk.NewInt(100000000000))
				}
				keeper.coinsFromUnstakedToStaked(context, test.application, test.amount)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, sdk.NewInt(100000000000))
				keeper.coinsFromUnstakedToStaked(context, test.application, test.amount)
				staked := keeper.GetStakedTokens(context)
				assert.True(t, test.amount.Add(sdk.NewInt(100000000000)).Equal(staked), "values do not match")
			}
		})
	}
}

func TestCoinsFromStakedToUnstaked(t *testing.T) {
	application := getBondedApplication()
	applicationAddress := application.Address

	tests := []struct {
		name        string
		amount      sdk.Int
		expected    string
		application types.Application
		panics      bool
	}{
		{
			name:        "unstake coins from pool",
			application: types.Application{Address: applicationAddress, StakedTokens: sdk.NewInt(10)},
			amount:      sdk.NewInt(110),
			panics:      false,
		},
		{
			name:        "panics if negative ammount",
			application: types.Application{Address: applicationAddress, StakedTokens: sdk.NewInt(-1)},
			amount:      sdk.NewInt(-1),
			expected:    fmt.Sprintf("negative coin amount: -1"),
			panics:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected)
				}()
				if strings.Contains(test.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, sdk.NewInt(100))
				}
				keeper.coinsFromStakedToUnstaked(context, test.application)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, sdk.NewInt(100))
				keeper.coinsFromStakedToUnstaked(context, test.application)
				unstaked := keeper.GetUnstakedTokens(context)
				assert.True(t, test.amount.Equal(unstaked), "values do not match")
			}
		})
	}
}

func TestBurnStakedTokens(t *testing.T) {
	application := getBondedApplication()
	applicationAddress := application.Address

	supplySize := sdk.NewInt(100000000000)
	tests := []struct {
		name        string
		expected    string
		application types.Application
		burnAmount  sdk.Int
		amount      sdk.Int
		errs        bool
	}{
		{
			name:        "burn coins from pool",
			application: types.Application{Address: applicationAddress},
			burnAmount:  sdk.NewInt(5),
			amount:      sdk.NewInt(10),
			errs:        false,
		},
		{
			name:        "errs trying to burn from pool",
			application: types.Application{Address: applicationAddress},
			burnAmount:  sdk.NewInt(-1),
			amount:      sdk.NewInt(10),
			errs:        true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.errs {
			case true:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, supplySize)
				keeper.coinsFromUnstakedToStaked(context, test.application, test.amount)
				err := keeper.burnStakedTokens(context, test.burnAmount)
				assert.Nil(t, err, "error is not nil")
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.application.Address, supplySize)
				keeper.coinsFromUnstakedToStaked(context, test.application, test.amount)
				err := keeper.burnStakedTokens(context, test.burnAmount)
				if err != nil {
					t.Fail()
				}
				staked := keeper.GetStakedTokens(context)
				assert.True(t, test.amount.Sub(test.burnAmount).Add(supplySize).Equal(staked), "values do not match")
			}
		})
	}
}
