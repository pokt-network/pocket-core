package keeper

import (
	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/pokt-network/pocket-core/types"
)

func TestSupply(t *testing.T) {
	initialPower := int64(100)
	initTokens := sdk.TokensFromConsensusPower(initialPower)
	nAccs := int64(4)

	ctx, keeper := createTestInput(t, false, initialPower, nAccs)

	total := keeper.GetSupply(ctx).GetTotal()
	expectedTotal := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, initTokens.MulRaw(nAccs)))

	require.Equal(t, expectedTotal, total)
}

func TestValidatePermissions(t *testing.T) {
	nAccs := int64(0)
	initialPower := int64(100)
	_, keeper := createTestInput(t, false, initialPower, nAccs)

	err := keeper.ValidatePermissions(multiPermAcc)
	require.NoError(t, err)

	err = keeper.ValidatePermissions(randomPermAcc)
	require.NoError(t, err)

	// unregistered permissions
	otherAcc := types.NewEmptyModuleAccount("other", "other")
	err = keeper.ValidatePermissions(otherAcc)
	require.Error(t, err)
}

func TestKeeper(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	addr := sdk.Address([]byte("addr1"))
	addr2 := sdk.Address([]byte("addr2"))
	acc, _ := input.Keeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	input.Keeper.SetAccount(ctx, acc)
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins()))

	_ = input.Keeper.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))

	// Test HasCoins
	require.True(t, input.Keeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, input.Keeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))
	require.False(t, input.Keeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15))))
	require.False(t, input.Keeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 5))))

	// Test AddCoins
	_, _ = input.Keeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 25))))

	_, _ = input.Keeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 15)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 15), sdk.NewInt64Coin("foocoin", 25))))

	// Test SubtractCoins
	_, _ = input.Keeper.SubtractCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10)))
	_, _ = input.Keeper.SubtractCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 5)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15))))

	_, _ = input.Keeper.SubtractCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 11)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15))))

	_, _ = input.Keeper.SubtractCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15))))
	require.False(t, input.Keeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 1))))

	// Test SendCoins
	_ = input.Keeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, input.Keeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))

	err2 := input.Keeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 50)))
	require.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, input.Keeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))

	_, _ = input.Keeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 30)))
	_ = input.Keeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)))
	require.True(t, input.Keeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5))))
	require.True(t, input.Keeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10))))
}

func TestSendKeeper(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx
	sendKeeper := input.Keeper
	addr := sdk.Address([]byte("addr1"))
	addr2 := sdk.Address([]byte("addr2"))
	acc, _ := input.Keeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	input.Keeper.SetAccount(ctx, acc)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins()))

	_ = input.Keeper.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10)))
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))

	// Test HasCoins
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15))))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 5))))

	_ = input.Keeper.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15)))

	// Test SendCoins
	_ = sendKeeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5)))
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))

	err := sendKeeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 50)))
	require.Implements(t, (*sdk.Error)(nil), err)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))

	_, _ = input.Keeper.AddCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 30)))
	_ = sendKeeper.SendCoins(ctx, addr, addr2, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)))
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5))))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10))))

	// validate coins with invalid denoms or negative values cannot be sent
	// NOTE: We must use the Coin literal as the constructor does not allow
	// negative values.
	err = sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.Coin{Denom: "FOOCOIN", Amount: sdk.NewInt(-5)}})
	require.Error(t, err)
}

func TestViewKeeper(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx
	//paramSpace := input.pk.Subspace(types.DefaultCodespace)
	viewKeeper := input.Keeper

	addr := sdk.Address([]byte("addr1"))
	acc, _ := input.Keeper.NewAccountWithAddress(ctx, addr)

	// Test GetCoins/SetCoins
	input.Keeper.SetAccount(ctx, acc)
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins()))

	_ = input.Keeper.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10)))
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))

	// Test HasCoins
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 10))))
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 5))))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("foocoin", 15))))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("barcoin", 5))))
}

func TestModuleUpgrade(t *testing.T) {
	nAccs := int64(4)
	ctx, keeper := createTestInput(t, false, initialPower, nAccs)
	accs := keeper.GetAllAccounts(ctx)
	p := keeper.GetParams(ctx)
	keeper.ConvertState(ctx)
	keeper.Cdc.SetUpgradeOverride(true)
	accs2 := keeper.GetAllAccounts(ctx)
	p2 := keeper.GetParams(ctx)
	assert.Equal(t, accs, accs2)
	assert.Equal(t, p, p2)
}
