package keeper

import (
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/x/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParamKeyTable(t *testing.T) {
	p := params.NewKeyTable().RegisterParamSet(&types.Params{})
	assert.Equal(t, ParamKeyTable(), p)
}

func TestKeeper_SessionNodeCount(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	sessNodeCount := keeper.SessionNodeCount(ctx)
	assert.NotNil(t, sessNodeCount)
	assert.NotEmpty(t, sessNodeCount)
	assert.Equal(t, types.DefaultSessionNodeCount, sessNodeCount)
}

func TestKeeper_ClaimExpiration(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	claimExpiration := keeper.ClaimExpiration(ctx)
	assert.NotNil(t, claimExpiration)
	assert.NotEmpty(t, claimExpiration)
	assert.Equal(t, types.DefaultClaimExpiration, claimExpiration)
}

func TestKeeper_SessionFrequency(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	sessFrequency := keeper.SessionFrequency(ctx)
	assert.NotNil(t, sessFrequency)
	assert.NotEmpty(t, sessFrequency)
	assert.Equal(t, int64(nodeTypes.DefaultSessionBlocktime), sessFrequency)
}

func TestKeeper_ProofWaitingPeriod(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	proofWaiting := keeper.ProofWaitingPeriod(ctx)
	assert.NotNil(t, proofWaiting)
	assert.NotEmpty(t, proofWaiting)
	assert.Equal(t, types.DefaultProofWaitingPeriod, proofWaiting)
}

func TestKeeper_SupportedBlockchains(t *testing.T) {
	ctx, _, _, _, keeper := createTestInput(t, false)
	supportedBlockchains := keeper.SupportedBlockchains(ctx)
	assert.Equal(t, types.DefaultSupportedBlockchains, supportedBlockchains)
}

func TestKeeper_GetParams(t *testing.T) {
	ctx, _, _, _, k := createTestInput(t, false)
	p := types.Params{
		SessionNodeCount:     k.SessionNodeCount(ctx),
		ProofWaitingPeriod:   k.ProofWaitingPeriod(ctx),
		SupportedBlockchains: k.SupportedBlockchains(ctx),
		ClaimExpiration:      k.ClaimExpiration(ctx),
	}
	paramz := k.GetParams(ctx)
	assert.NotNil(t, paramz)
	assert.Equal(t, p, paramz)
}

func TestKeeper_SetParams(t *testing.T) {
	ctx, _, _, _, k := createTestInput(t, false)
	sessionNodeCount := int64(17)
	pwp := int64(22)
	sb := []string{"ethereum"}
	p := types.Params{
		SessionNodeCount:     sessionNodeCount,
		ProofWaitingPeriod:   pwp,
		SupportedBlockchains: sb,
	}
	k.SetParams(ctx, p)
	paramz := k.GetParams(ctx)
	assert.Equal(t, paramz, p)
}
