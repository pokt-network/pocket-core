package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/supply/exported"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCoinsFromUnstakedToStaked(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name      string
		expected  string
		validator types.Validator
		amount    sdk.Int
		errors    bool
	}{
		{
			name:      "stake coins on pool",
			validator: types.Validator{Address: validatorAddress},
			amount:    sdk.NewInt(10),
			errors:    true,
		},
		{
			name:      "error if negative ammount",
			validator: types.Validator{Address: validatorAddress},
			amount:    sdk.NewInt(-1),
			expected:  fmt.Sprintf("negative coin amount: -1"),
			errors:    true,
		},
		{name: "error if no supply is set",
			validator: types.Validator{Address: validatorAddress},
			expected:  fmt.Sprintf("insufficient account funds"),
			amount:    sdk.NewInt(10),
			errors:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.errors {
			case true:
				if strings.Contains(test.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100000000000))
				}
				err := keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				assert.NotNil(t, err)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100000000000))
				err := keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				assert.Nil(t, err)
				staked := keeper.GetStakedTokens(context)
				assert.True(t, test.amount.Add(sdk.NewInt(100000000000)).Equal(staked), "values do not match")
			}
		})
	}
}

func TestCoinsFromStakedToUnstaked(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name      string
		amount    sdk.Int
		expected  string
		validator types.Validator
		panics    bool
	}{
		{
			name:      "unstake coins from pool",
			validator: types.Validator{Address: validatorAddress, StakedTokens: sdk.NewInt(10)},
			amount:    sdk.NewInt(110),
			panics:    false,
		},
		{
			name:      "errors if negative ammount",
			validator: types.Validator{Address: validatorAddress, StakedTokens: sdk.NewInt(-1)},
			amount:    sdk.NewInt(-1),
			expected:  fmt.Sprintf("negative coin amount: -1"),
			panics:    true,
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
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100))
				}
				keeper.coinsFromStakedToUnstaked(context, test.validator)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, sdk.NewInt(100))
				keeper.coinsFromStakedToUnstaked(context, test.validator)
				unstaked := keeper.GetUnstakedTokens(context)
				assert.True(t, test.amount.Equal(unstaked), "values do not match")
			}
		})
	}
}

func TestBurnStakedTokens(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	supplySize := sdk.NewInt(100000000000)
	tests := []struct {
		name       string
		expected   string
		validator  types.Validator
		burnAmount sdk.Int
		amount     sdk.Int
		errs       bool
	}{
		{
			name:       "burn coins from pool",
			validator:  types.Validator{Address: validatorAddress},
			burnAmount: sdk.NewInt(5),
			amount:     sdk.NewInt(10),
			errs:       false,
		},
		{
			name:       "errs trying to burn from pool",
			validator:  types.Validator{Address: validatorAddress},
			burnAmount: sdk.NewInt(-1),
			amount:     sdk.NewInt(10),
			errs:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.errs {
			case true:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, supplySize)
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
				err := keeper.burnStakedTokens(context, test.burnAmount)
				assert.Nil(t, err, "error is not nil")
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, test.validator.Address, supplySize)
				keeper.coinsFromUnstakedToStaked(context, test.validator, test.amount)
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

func TestPool_GetFeePool(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"gets fee pool",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			got := keeper.getFeePool(context)

			if _, ok := got.(exported.ModuleAccountI); !ok {
				t.Errorf("KeeperPool.getFeePool()= %v", ok)
			}
		})
	}
}

func TestPool_StakedRatio(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name    string
		amount  sdk.Dec
		address sdk.Address
	}{
		{"return 0 if stake supply is lower than 0", sdk.ZeroDec(), validatorAddress},
		{"return supply", sdk.NewDec(1), validatorAddress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			if !tt.amount.Equal(sdk.ZeroDec()) {
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			}

			if got := keeper.StakedRatio(context); !got.Equal(tt.amount) {
				t.Errorf("KeeperPool.StakedRatio()= %v, %v", got.String(), tt.amount.String())
			}
		})
	}
}
