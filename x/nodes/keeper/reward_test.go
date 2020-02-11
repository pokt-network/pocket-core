package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type args struct {
	amount      sdk.Int
	address     sdk.Address
	consAddress sdk.Address
}

func TestSetandGetValidatorAward(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		expectedFind  bool
	}{
		{
			name:          "can set award",
			expectedCoins: sdk.NewInt(1),
			expectedFind:  true,
			args:          args{amount: sdk.NewInt(int64(1)), address: validatorAddress},
		},
		{
			name:          "can get award",
			expectedCoins: sdk.NewInt(2),
			expectedFind:  true,
			args:          args{amount: sdk.NewInt(int64(2)), address: validatorAddress},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.setValidatorAward(context, test.args.amount, test.args.address)
			coins, found := keeper.getValidatorAward(context, test.args.address)
			fmt.Println(coins, test.expectedCoins)
			assert.True(t, test.expectedCoins.Equal(coins), "coins don't match")
			assert.Equal(t, test.expectedFind, found, "finds don't match")

		})
	}
}

func TestSetAndGetProposer(t *testing.T) {
	validator := getStakedValidator()
	consAddress := validator.GetAddress()

	tests := []struct {
		name            string
		args            args
		expectedAddress sdk.Address
	}{
		{
			name:            "can set the preivous proposer",
			args:            args{consAddress: consAddress},
			expectedAddress: consAddress,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.SetPreviousProposer(context, test.args.consAddress)
			receivedAddress := keeper.GetPreviousProposer(context)
			assert.True(t, test.expectedAddress.Equals(receivedAddress), "addresses do not match ")
		})
	}
}

func TestDeleteValidatorAward(t *testing.T) {
	validator := getStakedValidator()
	validatorAddress := validator.Address

	tests := []struct {
		name          string
		args          args
		expectedCoins sdk.Int
		expectedFind  bool
	}{
		{
			name:          "can delete award",
			expectedCoins: sdk.NewInt(0),
			expectedFind:  false,
			args:          args{amount: sdk.NewInt(int64(1)), address: validatorAddress},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)

			keeper.setValidatorAward(context, test.args.amount, test.args.address)
			keeper.deleteValidatorAward(context, test.args.address)
			_, found := keeper.getValidatorAward(context, test.args.address)
			assert.Equal(t, test.expectedFind, found, "finds do not match")

		})
	}
}

func TestGetProposerAllocation(t *testing.T) {
	tests := []struct {
		name               string
		expectedPercentage sdk.Int
	}{
		{
			name:               "get reward percentage",
			expectedPercentage: sdk.NewInt(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			percentage := keeper.getProposerAllocaiton(context)
			assert.True(t, test.expectedPercentage.Equal(percentage), "percentages do not match")
		})
	}
}

func TestMint(t *testing.T) {
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
			name:     "mints a coin",
			amount:   sdk.NewInt(90),
			expected: fmt.Sprintf("a custom reward of "),
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
				_ = keeper.mint(context, test.amount, test.address)
			default:
				result := keeper.mint(context, test.amount, test.address)
				assert.Contains(t, result.Log, test.expected, "does not contain message")
				coins := keeper.coinKeeper.GetCoins(context, sdk.Address(test.address))
				assert.True(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), test.amount)).IsEqual(coins), "coins should match")
			}
		})
	}
}

func TestMintValidatorAwards(t *testing.T) {
	validatorAddress := getRandomValidatorAddress()
	tests := []struct {
		name     string
		amount   sdk.Int
		expected string
		address  sdk.Address
		panics   bool
	}{
		{
			name:     "mints a coin",
			amount:   sdk.NewInt(90),
			expected: fmt.Sprintf("was successfully minted to %s", validatorAddress.String()),
			address:  validatorAddress,
			panics:   false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.setValidatorAward(context, test.amount, test.address)

			keeper.mintNodeRelayRewards(context)
			coins := keeper.coinKeeper.GetCoins(context, sdk.Address(test.address))
			expected := keeper.NodeCutOfReward(context).Mul(test.amount).Quo(sdk.NewInt(100))
			assert.True(t, sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), expected)).IsEqual(coins), "coins should match")
		})
	}
}

func TestKeeper_GetTotalCustomValidatorAwards(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx sdk.Context
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   sdk.Int
	}{
		{"Test GetTotalCustomValidatorAwards", fields{keeper: keeper},
			args{ctx: context}, sdk.ZeroInt()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper
			if got := k.GetTotalCustomValidatorAwards(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTotalCustomValidatorAwards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_AwardCoinsTo(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx     sdk.Context
		relays  sdk.Int
		address sdk.Address
	}
	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test AwardCoinsTo", fields{keeper: keeper},
			args{
				ctx:     context,
				relays:  sdk.ZeroInt(),
				address: getRandomValidatorAddress(),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper

			k.setValidatorAward(tt.args.ctx, sdk.ZeroInt(), tt.args.address)
			k.AwardCoinsTo(tt.args.ctx, tt.args.relays, tt.args.address)

		})
	}
}

func TestKeeper_rewardFromFees(t *testing.T) {
	type fields struct {
		keeper Keeper
	}
	type args struct {
		ctx              sdk.Context
		previousProposer sdk.Address
	}
	stakedValidator := getStakedValidator()

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test rewardFromFees", fields{keeper: keeper},
			args{
				ctx:              context,
				previousProposer: stakedValidator.GetAddress(),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.keeper

			k.rewardFromFees(tt.args.ctx, tt.args.previousProposer)

		})
	}
}
