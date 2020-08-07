package keeper

import (
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

func TestParamKeyTable(t *testing.T) {
	p := sdk.NewKeyTable().RegisterParamSet(&types.Params{})
	assert.Equal(t, ParamKeyTable(), p)
}

func TestKeeper_SessionNodeCount(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	sessNodeCount := keeper.SessionNodeCount(ctx)
	assert.NotNil(t, sessNodeCount)
	assert.NotEmpty(t, sessNodeCount)
	assert.Equal(t, types.DefaultSessionNodeCount, sessNodeCount)
}

func TestKeeper_ClaimExpiration(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	claimExpiration := keeper.ClaimExpiration(ctx)
	assert.NotNil(t, claimExpiration)
	assert.NotEmpty(t, claimExpiration)
	assert.Equal(t, types.DefaultClaimExpiration, claimExpiration)
}

func TestKeeper_ReplayAttackBurnMultiplier(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	rabm := keeper.ReplayAttackBurnMultiplier(ctx)
	assert.NotNil(t, rabm)
	assert.NotEmpty(t, rabm)
	assert.Equal(t, types.DefaultReplayAttackBurnMultiplier, rabm)
}

func TestKeeper_SessionFrequency(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	sessFrequency := keeper.BlocksPerSession(ctx)
	assert.NotNil(t, sessFrequency)
	assert.NotEmpty(t, sessFrequency)
	assert.Equal(t, int64(nodeTypes.DefaultSessionBlocktime), sessFrequency)
}

func TestKeeper_ClaimSubmissionWindow(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	proofWaiting := keeper.ClaimSubmissionWindow(ctx)
	assert.NotNil(t, proofWaiting)
	assert.NotEmpty(t, proofWaiting)
	assert.Equal(t, types.DefaultClaimSubmissionWindow, proofWaiting)
}

func TestKeeper_SupportedBlockchains(t *testing.T) {
	ctx, _, _, _, keeper, _, _ := createTestInput(t, false)
	supportedBlockchains := keeper.SupportedBlockchains(ctx)
	assert.Equal(t, []string{getTestSupportedBlockchain()}, supportedBlockchains)
}

func TestKeeper_GetParams(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	p := types.Params{
		SessionNodeCount:           k.SessionNodeCount(ctx),
		ClaimSubmissionWindow:      k.ClaimSubmissionWindow(ctx),
		SupportedBlockchains:       k.SupportedBlockchains(ctx),
		ClaimExpiration:            k.ClaimExpiration(ctx),
		ReplayAttackBurnMultiplier: k.ReplayAttackBurnMultiplier(ctx),
		MinimumNumberOfProofs:      k.MinimumNumberOfProofs(ctx),
	}
	paramz := k.GetParams(ctx)
	assert.NotNil(t, paramz)
	assert.Equal(t, p, paramz)
}

func TestKeeper_SetParams(t *testing.T) {
	ctx, _, _, _, k, _, _ := createTestInput(t, false)
	sessionNodeCount := int64(17)
	pwp := int64(22)
	sb := []string{"ethereum"}
	p := types.Params{
		SessionNodeCount:      sessionNodeCount,
		ClaimSubmissionWindow: pwp,
		SupportedBlockchains:  sb,
	}
	k.SetParams(ctx, p)
	paramz := k.GetParams(ctx)
	assert.Equal(t, paramz, p)
}
