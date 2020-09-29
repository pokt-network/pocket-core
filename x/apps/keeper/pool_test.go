package keeper

import (
	"strings"
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
	"github.com/stretchr/testify/assert"
)

func TestPool_CoinsFromUnstakedToStaked(t *testing.T) {
	application := getStakedApplication()
	applicationAddress := application.Address

	tests := []struct {
		name        string
		want        string
		application types.Application
		amount      sdk.BigInt
		errors      bool
	}{
		{
			name:        "stake coins on pool",
			application: types.Application{Address: applicationAddress},
			amount:      sdk.NewInt(10),
			errors:      false,
		},
		{
			name:        "errors if negative ammount",
			application: types.Application{Address: applicationAddress},
			amount:      sdk.NewInt(-1),
			want:        "negative coin amount: -1",
			errors:      true,
		},
		{name: "errors if no supply is set",
			application: types.Application{Address: applicationAddress},
			want:        "insufficient account funds",
			amount:      sdk.NewInt(10),
			errors:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch tt.errors {
			case true:
				if strings.Contains(tt.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
				}
				err := keeper.coinsFromUnstakedToStaked(context, tt.application, tt.amount)
				assert.NotNil(t, err)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
				err := keeper.coinsFromUnstakedToStaked(context, tt.application, tt.amount)
				assert.Nil(t, err)
				if got := keeper.GetStakedTokens(context); !tt.amount.Add(sdk.NewInt(100000000000)).Equal(got) {
					t.Errorf("KeeperCoins.FromUnstakedToStaked()= %v, want %v", got, tt.amount.Add(sdk.NewInt(100000000000)))
				}
			}
		})
	}
}

func TestPool_CoinsFromStakedToUnstaked(t *testing.T) {
	application := getStakedApplication()
	applicationAddress := application.Address

	tests := []struct {
		name        string
		amount      sdk.BigInt
		want        string
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
			name:        "errors if negative ammount",
			application: types.Application{Address: applicationAddress, StakedTokens: sdk.NewInt(-1)},
			amount:      sdk.NewInt(-1),
			want:        "negative coin amount: -1",
			panics:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch tt.panics {
			case true:
				defer func() {
					if err := recover().(error); !strings.Contains(err.Error(), tt.want) {
						t.Errorf("KeeperCoins.FromStakedToUnstaked()= %v, want %v", err.Error(), tt.want)
					}
				}()
				if strings.Contains(tt.name, "setup") {
					addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
					sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100))
				}
				_ = keeper.coinsFromStakedToUnstaked(context, tt.application)
			default:
			}
		})
	}
}

func TestPool_BurnStakedTokens(t *testing.T) {
	application := getStakedApplication()
	applicationAddress := application.Address

	supplySize := sdk.NewInt(100000000000)
	tests := []struct {
		name        string
		expected    string
		application types.Application
		burnAmount  sdk.BigInt
		amount      sdk.BigInt
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch tt.errs {
			case true:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, supplySize)
				_ = keeper.coinsFromUnstakedToStaked(context, tt.application, tt.amount)
				if err := keeper.burnStakedTokens(context, tt.burnAmount); err != nil {
					t.Errorf("KeeperCoins.BurnStakedTokens()= %v, want nil", err)
				}
			default:
				addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
				sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, supplySize)
				_ = keeper.coinsFromUnstakedToStaked(context, tt.application, tt.amount)
				err := keeper.burnStakedTokens(context, tt.burnAmount)
				if err != nil {
					t.Fail()
				}
				if got := keeper.GetStakedTokens(context); !tt.amount.Sub(tt.burnAmount).Add(supplySize).Equal(got) {
					t.Errorf("KeeperCoins.BurnStakedTokens()= %v, want %v", got, tt.amount.Sub(tt.burnAmount).Add(supplySize))
				}
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
	application := getStakedApplication()
	applicationAddress := application.Address

	tests := []struct {
		name    string
		amount  sdk.BigDec
		address sdk.Address
	}{
		{"return 0 if stake supply is lower than 0", sdk.ZeroDec(), applicationAddress},
		{"return supply", sdk.NewDec(1), applicationAddress},
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
