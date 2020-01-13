package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitExportGenesis(t *testing.T) {
	ctx, _, _, k := createTestInput(t, false)
	p := types.Params{
		SessionNodeCount:     10,
		ProofWaitingPeriod:   22,
		SupportedBlockchains: []string{"eth"},
		ClaimExpiration:      55,
	}
	genesisState := types.GenesisState{
		Params: p,
		Proofs: []types.StoredInvoice(nil),
		Claims: []types.MsgClaim(nil),
	}
	InitGenesis(ctx, k, genesisState)
	assert.Equal(t, k.GetParams(ctx), p)
	gen := ExportGenesis(ctx, k)
	assert.Equal(t, genesisState, gen)
}
