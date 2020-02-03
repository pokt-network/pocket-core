package keeper

import (
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinsFromDAOToValidator(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name     string
		amount   sdk.Int
		expected string
		address  sdk.Address
		panics   bool
	}{
		{
			name:     "sends coin to account",
			amount:   sdk.NewInt(90),
			expected: fmt.Sprintf("was successfully minted to %s", validatorAddress.String()),
			address:  validatorAddress,
			panics:   false,
		},
		{
			name:     "errors invalid ammount of coins",
			amount:   sdk.NewInt(-1),
			expected: fmt.Sprintf("negative coin amount: -1"),
			address:  validatorAddress,
			panics:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			switch test.panics {
			case true:
				defer func() {
					err := recover().(error)
					assert.Contains(t, err.Error(), test.expected, "error does not match")
				}()
				keeper.coinsFromDAOToValidator(context, types.Validator{Address: test.address}, test.amount)
			default:
				addMintedCoinsToModule(t, context, &keeper, types.DAOPoolName)
				keeper.coinsFromDAOToValidator(context, types.Validator{Address: test.address}, test.amount)
				coins := keeper.coinKeeper.GetCoins(context, sdk.Address(test.address))
				assert.True(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)).IsEqual(coins), "coins should match")
			}
		})
	}
}
