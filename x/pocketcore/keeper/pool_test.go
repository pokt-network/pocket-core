package keeper

import (
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeeper_StakeDenom(t *testing.T) {
	ctx, _, _, _, k := createTestInput(t, false)
	stakeDenom := types.DefaultStakeDenom
	assert.Equal(t, stakeDenom, k.posKeeper.StakeDenom(ctx))
}

func TestKeeper_GetNodesStakedTokens(t *testing.T) {
	ctx, vals, _, _, k := createTestInput(t, false)
	assert.NotZero(t, len(vals))
	tokens := vals[0].StakedTokens
	assert.Equal(t, k.GetNodesStakedTokens(ctx), tokens.Mul(types.NewInt(int64(len(vals)))))
}

func TestKeeper_GetAppsStakedTokens(t *testing.T) {
	ctx, _, apps, _, k := createTestInput(t, false)
	assert.NotZero(t, len(apps))
	tokens := apps[0].StakedTokens
	assert.Equal(t, k.GetNodesStakedTokens(ctx), tokens.Mul(types.NewInt(int64(len(apps)))))
}

func TestKeeper_GetTotalStakedTokens(t *testing.T) {
	ctx, vals, apps, _, k := createTestInput(t, false)
	assert.NotZero(t, len(apps))
	appToken := apps[0].StakedTokens
	appTokens := appToken.Mul(types.NewInt(int64(len(apps))))
	valToken := vals[0].StakedTokens
	valTokens := valToken.Mul(types.NewInt(int64(len(vals))))
	assert.Equal(t, k.GetTotalStakedTokens(ctx), appTokens.Add(valTokens))
}

//func TestKeeper_GetTotalTokens(t *testing.T) { todo
//	ctx, vals, apps, accs, k := createTestInput(t, false)
//	assert.NotZero(t, len(apps))
//	appToken := apps[0].StakedTokens
//	appTokens := appToken.Mul(types.NewInt(int64(len(apps))))
//	valToken := vals[0].StakedTokens
//	valTokens := valToken.Mul(types.NewInt(int64(len(vals))))
//	accToken := accs[0].GetCoins().AmountOf(k.StakeDenom(ctx))
//	accTokens:= accToken.Mul(types.NewInt(int64(len(accs))))
//	fmt.Println("ACC Tokens " + accTokens.String())
//	fmt.Println("Val Tokens " + valTokens.String())
//	fmt.Println("App Tokens " + appTokens.String())
//	fmt.Println("total toakens " + appTokens.Add(valTokens).Add(accTokens).String())
//	fmt.Println("actual tokens ", k.GetTotalTokens(ctx))
//	assert.Equal(t, k.GetTotalTokens(ctx), appTokens.Add(valTokens))
//}
