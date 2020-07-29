package keeper

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

func TestInitGenesis(t *testing.T) {
	gs := types.GenesisState{
		Params: types.Params{
			ACL:      createTestACL(),
			DAOOwner: getRandomValidatorAddress(),
			Upgrade:  types.Upgrade{},
		},
		DAOTokens: sdk.NewInt(1000),
	}
	ctx, k := createTestKeeperAndContext(t, false)
	assert.Equal(t, k.InitGenesis(ctx, gs), []abci.ValidatorUpdate{})
}

func TestExportGenesis(t *testing.T) {
	ctx, k := createTestKeeperAndContext(t, false)
	d := types.DefaultGenesisState()
	d.Params.ACL = createTestACL()
	d.Params.Upgrade = types.Upgrade{}
	assert.Equal(t, k.ExportGenesis(ctx).Params.ACL.String(), d.Params.ACL.String())
	assert.Equal(t, k.ExportGenesis(ctx).DAOTokens.Int64(), d.DAOTokens.Int64())
}
