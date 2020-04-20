package keeper

import (
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_StakeDenom(t *testing.T) {
	ctx, _, _, _, k, _ := createTestInput(t, false)
	stakeDenom := types.DefaultStakeDenom
	assert.Equal(t, stakeDenom, k.posKeeper.StakeDenom(ctx))
}

func TestKeeper_GetNodesStakedTokens(t *testing.T) {
	ctx, vals, _, _, k, _ := createTestInput(t, false)
	assert.NotZero(t, len(vals))
	tokens := vals[0].StakedTokens
	assert.Equal(t, k.GetNodesStakedTokens(ctx), tokens.Mul(types.NewInt(int64(len(vals)))))
}

func TestKeeper_GetAppsStakedTokens(t *testing.T) {
	ctx, _, apps, _, k, _ := createTestInput(t, false)
	assert.NotZero(t, len(apps))
	tokens := apps[0].StakedTokens
	assert.Equal(t, k.GetAppStakedTokens(ctx).Int64(), tokens.Mul(types.NewInt(int64(len(apps)))).Int64())
}

func TestKeeper_GetTotalStakedTokens(t *testing.T) {
	ctx, vals, apps, _, k, _ := createTestInput(t, false)
	assert.NotZero(t, len(apps))
	appToken := apps[0].StakedTokens
	appTokens := appToken.Mul(types.NewInt(int64(len(apps))))
	valToken := vals[0].StakedTokens
	valTokens := valToken.Mul(types.NewInt(int64(len(vals))))
	assert.Equal(t, k.GetTotalStakedTokens(ctx), appTokens.Add(valTokens))
}
