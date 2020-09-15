package keeper

import (
	"testing"

	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_GetParams(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	d := types.DefaultGenesisState()
	d.Params.ACL = createTestACL()
	assert.Equal(t, k.GetParams(ctx).String(), d.Params.String())
}

func TestKeeper_GetACL(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	d := types.DefaultGenesisState()
	d.Params.ACL = createTestACL()
	assert.Equal(t, k.GetACL(ctx).String(), d.Params.ACL.String())
}

func TestKeeper_SetParamsAndGetDAOOwner(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	d := types.DefaultParams()
	d.DAOOwner = getRandomValidatorAddress()
	k.SetParams(ctx, d)
	assert.Equal(t, k.GetDAOOwner(ctx).String(), d.DAOOwner.String())
}
