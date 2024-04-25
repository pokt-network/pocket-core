package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/auth/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
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

// There are a few orphaned addresses that are owned by the protocol to manage staked pool funds.
// For example, when one stakes, their funds are moved to the pool. Once an unstake occurs, the funds are moved back from the pool to the owner.
// These are owned by the protocol and their supply can be checked via `pocket query supply`
func TestGetPoolAddresses(t *testing.T) {
	tests := []struct {
		name    string
		pool    string
		address string
	}{
		{
			name:    "Staked nodes pool address",
			pool:    nodeTypes.StakedPoolName,
			address: "8ef97b488e66a2b2e89a3b4999549816768910fb",
		},
		{
			name:    "App nodes pool address",
			pool:    appTypes.StakedPoolName,
			address: "63533fb8f43b4883a1f37265f1561ce7b1c6c307",
		},
		{
			name:    "DAO owned address",
			pool:    govTypes.DAOAccountName,
			address: "6386713deb27b609daad5e2e32ee6591753e5f4e",
		},
	}
	ctx, keeper := createTestInput(t, false, initialPower, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseAcc, _ := keeper.NewAccountWithAddress(ctx, types.NewModuleAddress(tt.pool))
			assert.Equal(t, tt.address, baseAcc.Address.String())
		})
	}
}
