package keeper

import (
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/stretchr/testify/require"
)

const initialPower = int64(100)

var (
	// holderAcc     = types.NewEmptyModuleAccount(holder)
	// burnerAcc     = types.NewEmptyModuleAccount(types.Burner, types.Burner)
	// minterAcc     = types.NewEmptyModuleAccount(types.Minter, types.Minter)
	// multiPermAcc  = types.NewEmptyModuleAccount(multiPerm, types.Burner, types.Minter, types.Staking)
	// randomPermAcc = types.NewEmptyModuleAccount(randomPerm, "random")

	initTokens = sdk.TokensFromConsensusPower(initialPower)
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, initTokens))
)

func getCoinsByName(ctx sdk.Ctx, k Keeper, moduleName string) sdk.Coins {
	moduleAddress := k.GetModuleAddress(moduleName)
	macc := k.GetAccount(ctx, moduleAddress)
	if macc == nil {
		return sdk.Coins(nil)
	}
	return macc.GetCoins()
}

func TestSendCoins(t *testing.T) {
	nAccs := int64(4)
	ctx, keeper := createTestInput(t, false, initialPower, nAccs)
	baseAcc, _ := keeper.NewAccountWithAddress(ctx, types.NewModuleAddress("baseAcc"))
	err := holderAcc.SetCoins(initCoins)
	require.NoError(t, err)
	keeper.SetModuleAccount(ctx, holderAcc)
	keeper.SetModuleAccount(ctx, burnerAcc)
	keeper.SetAccount(ctx, baseAcc)
	err = keeper.SendCoinsFromModuleToModule(ctx, "", holderAcc.GetName(), initCoins)
	require.Error(t, err)
	err = keeper.SendCoinsFromModuleToModule(ctx, types.Burner, "", initCoins)
	require.Error(t, err)
	err = keeper.SendCoinsFromModuleToAccount(ctx, "", baseAcc.GetAddress(), initCoins)
	require.Error(t, err)
	err = keeper.SendCoinsFromModuleToAccount(ctx, holderAcc.GetName(), baseAcc.GetAddress(), initCoins.Add(initCoins))
	require.Error(t, err)
	err = keeper.SendCoinsFromModuleToModule(ctx, holderAcc.GetName(), types.Burner, initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(nil), getCoinsByName(ctx, keeper, holderAcc.GetName()))
	require.Equal(t, initCoins, getCoinsByName(ctx, keeper, types.Burner))
	err = keeper.SendCoinsFromModuleToAccount(ctx, types.Burner, baseAcc.GetAddress(), initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(nil), getCoinsByName(ctx, keeper, types.Burner))
	require.Equal(t, initCoins, keeper.GetAccount(ctx, baseAcc.GetAddress()).GetCoins())
	err = keeper.SendCoinsFromAccountToModule(ctx, baseAcc.GetAddress(), types.Burner, initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(nil), keeper.GetAccount(ctx, baseAcc.GetAddress()).GetCoins())
	require.Equal(t, initCoins, getCoinsByName(ctx, keeper, types.Burner))
}

func TestMintCoins(t *testing.T) {
	nAccs := int64(4)
	ctx, keeper := createTestInput(t, false, initialPower, nAccs)

	keeper.SetModuleAccount(ctx, burnerAcc)
	keeper.SetModuleAccount(ctx, minterAcc)
	keeper.SetModuleAccount(ctx, multiPermAcc)
	keeper.SetModuleAccount(ctx, randomPermAcc)

	initialSupply := keeper.GetSupply(ctx)

	require.Error(t, keeper.MintCoins(ctx, "", initCoins), "no module account")
	require.Error(t, keeper.MintCoins(ctx, types.Burner, initCoins), "invalid permission")
	require.Error(t, keeper.MintCoins(ctx, types.Minter, sdk.Coins{sdk.Coin{Denom: "denom", Amount: sdk.NewInt(-10)}}), "insufficient coins") //nolint

	require.Error(t, keeper.MintCoins(ctx, randomPerm, initCoins))

	err := keeper.MintCoins(ctx, types.Minter, initCoins)
	require.NoError(t, err)
	require.Equal(t, initCoins, getCoinsByName(ctx, keeper, types.Minter))
	got := keeper.GetSupply(ctx).GetTotal()
	want := initialSupply.GetTotal().Add(initCoins)
	require.Equal(t, want, got)

	// test same functionality on module account with multiple permissions
	initialSupply = keeper.GetSupply(ctx)

	err = keeper.MintCoins(ctx, multiPermAcc.GetName(), initCoins)
	require.NoError(t, err)
	require.Equal(t, initCoins, getCoinsByName(ctx, keeper, multiPermAcc.GetName()))
	require.Equal(t, initialSupply.GetTotal().Add(initCoins), keeper.GetSupply(ctx).GetTotal())

	require.Error(t, keeper.MintCoins(ctx, types.Burner, initCoins))
}

func TestBurnCoins(t *testing.T) {
	nAccs := int64(4)
	ctx, keeper := createTestInput(t, false, initialPower, nAccs)

	require.NoError(t, burnerAcc.SetCoins(initCoins))
	keeper.SetModuleAccount(ctx, burnerAcc)

	initialSupply := keeper.GetSupply(ctx)
	initialSupply = initialSupply.Inflate(initCoins)
	keeper.SetSupply(ctx, initialSupply)

	require.Error(t, keeper.BurnCoins(ctx, "", initCoins), "no module account")
	require.Error(t, keeper.BurnCoins(ctx, types.Minter, initCoins), "invalid permission")
	require.Error(t, keeper.BurnCoins(ctx, randomPerm, initialSupply.GetTotal()), "random permission")
	require.Error(t, keeper.BurnCoins(ctx, types.Burner, initialSupply.GetTotal()), "insufficient coins")

	err := keeper.BurnCoins(ctx, types.Burner, initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(nil), getCoinsByName(ctx, keeper, types.Burner))
	require.Equal(t, initialSupply.GetTotal().Sub(initCoins), keeper.GetSupply(ctx).GetTotal())

	// test same functionality on module account with multiple permissions
	initialSupply = keeper.GetSupply(ctx)
	initialSupply = initialSupply.Inflate(initCoins)
	keeper.SetSupply(ctx, initialSupply)

	require.NoError(t, multiPermAcc.SetCoins(initCoins))
	keeper.SetModuleAccount(ctx, multiPermAcc)

	err = keeper.BurnCoins(ctx, multiPermAcc.GetName(), initCoins)
	require.NoError(t, err)
	require.Equal(t, sdk.Coins(nil), getCoinsByName(ctx, keeper, multiPermAcc.GetName()))
	require.Equal(t, initialSupply.GetTotal().Sub(initCoins), keeper.GetSupply(ctx).GetTotal())
}
