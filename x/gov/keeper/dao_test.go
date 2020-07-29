package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestKeeper_GetDAOAccountAndTokens(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	acc := k.GetDAOAccount(ctx)
	assert.NotNil(t, acc)
	err := acc.SetCoins(sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000))))
	if err != nil {
		panic(fmt.Sprintf("unable to set dao tokens: %s", err.Error()))
	}
	k.AuthKeeper.SetModuleAccount(ctx, acc)
	acc = k.GetDAOAccount(ctx)
	daoActualCoins := acc.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	assert.Equal(t, daoActualCoins.Int64(), int64(1000))
	assert.Equal(t, k.GetDAOTokens(ctx).Int64(), daoActualCoins.Int64())
}

func TestKeeper_GetDAOBurn(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	acc := k.GetDAOAccount(ctx)
	assert.NotNil(t, acc)
	err := k.AuthKeeper.MintCoins(ctx, types.DAOAccountName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000))))
	acc = k.GetDAOAccount(ctx)
	daoActualCoins := acc.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	assert.Nil(t, err)
	assert.Equal(t, daoActualCoins.Int64(), int64(1000))
	assert.Equal(t, k.GetDAOTokens(ctx).Int64(), daoActualCoins.Int64())
	daoBurnResult := k.DAOBurn(ctx, k.GetDAOOwner(ctx), sdk.OneInt())
	acc = k.GetDAOAccount(ctx)
	daoActualCoins = acc.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	assert.Equal(t, daoActualCoins.Int64(), int64(999))
	// Test "message.sender" event emission
	assert.Equal(
		t,
		true,
		ContainsEvent(
			daoBurnResult.Events,
			abci.Event(
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(sdk.AttributeKeySender, k.GetDAOOwner(ctx).String()),
				),
			),
		),
	)
}

func TestKeeper_GetDAOTransfer(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	acc := k.GetDAOAccount(ctx)
	assert.NotNil(t, acc)
	err := k.AuthKeeper.MintCoins(ctx, types.DAOAccountName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(1000))))
	acc = k.GetDAOAccount(ctx)
	daoActualCoins := acc.GetCoins().AmountOf(sdk.DefaultStakeDenom)
	assert.Nil(t, err)
	assert.Equal(t, daoActualCoins.Int64(), int64(1000))
	assert.Equal(t, k.GetDAOTokens(ctx).Int64(), daoActualCoins.Int64())
	daoTransferResult := k.DAOTransferFrom(ctx, k.GetDAOOwner(ctx), k.AuthKeeper.GetModuleAddress("FAKE"), sdk.NewInt(1))
	acc = k.AuthKeeper.GetModuleAccount(ctx, "FAKE")
	assert.Equal(t, acc.GetCoins().AmountOf(sdk.DefaultStakeDenom).Int64(), sdk.OneInt().Int64())
	assert.Equal(t, k.GetDAOTokens(ctx).Int64(), int64(999))
	// Test "message.sender" event emission
	assert.Equal(
		t,
		true,
		ContainsEvent(
			daoTransferResult.Events,
			abci.Event(
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeySender, k.GetDAOAccount(ctx).GetAddress().String()),
				),
			),
		),
	)
}
