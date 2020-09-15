package keeper

import (
	"testing"

	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	holderAcc     = types.NewEmptyModuleAccount(holder)
	burnerAcc     = types.NewEmptyModuleAccount(types.Burner, types.Burner)
	minterAcc     = types.NewEmptyModuleAccount(types.Minter, types.Minter)
	multiPermAcc  = types.NewEmptyModuleAccount(multiPerm, types.Burner, types.Minter, types.Staking)
	randomPermAcc = types.NewEmptyModuleAccount(randomPerm, "random")
)

func TestSetAndGetAccounts(t *testing.T) {
	nAccs := int64(4)
	ctx, keeper := createTestInput(t, false, initialPower, nAccs)
	baseAcc, _ := keeper.NewAccountWithAddress(ctx, types.NewModuleAddress("baseAcc"))
	err := holderAcc.SetCoins(initCoins)
	require.NoError(t, err)
	err = baseAcc.SetCoins(initCoins)
	require.NoError(t, err)
	keeper.SetModuleAccount(ctx, holderAcc)
	keeper.SetAccount(ctx, baseAcc)

	gotHold := keeper.GetModuleAccount(ctx, holderAcc.GetName())
	assert.Equal(t, holderAcc, gotHold)
	assert.Equal(t, initCoins, gotHold.GetCoins())

	gotAcc := keeper.GetAccount(ctx, baseAcc.GetAddress())
	assert.Equal(t, baseAcc, gotAcc)
}
