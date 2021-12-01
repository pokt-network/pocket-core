package pocketcore

import (
	"testing"

	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
)

func TestInitExportGenesis(t *testing.T) {
	ctx, _, _, k, _ := createTestInput(t, false)
	p := types.Params{
		SessionNodeCount:      10,
		ClaimSubmissionWindow: 22,
		SupportedBlockchains:  []string{"eth"},
		ClaimExpiration:       55,
		MinimumNumberOfProofs: int64(5),
	}
	genesisState := types.GenesisState{
		Params: p,
		Claims: []types.MsgClaim(nil),
	}
	InitGenesis(ctx, k, genesisState)
	assert.Equal(t, k.GetParams(ctx), p)
	gen := ExportGenesis(ctx, k)
	assert.Equal(t, genesisState, gen)
}
