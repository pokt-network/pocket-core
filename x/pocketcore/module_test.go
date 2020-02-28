package pocketcore

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
	"testing"
)

func TestAppModule_Name(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, am.Name(), types.ModuleName)
	assert.Equal(t, am.Name(), types.ModuleName)
}

func TestAppModule_InitExportGenesis(t *testing.T) {
	p := types.Params{
		SessionNodeCount:     10,
		ProofWaitingPeriod:   22,
		SupportedBlockchains: []string{"eth"},
		ClaimExpiration:      55,
	}
	genesisState := types.GenesisState{
		Params: p,
		Proofs: []types.StoredEvidence(nil),
		Claims: []types.MsgClaim(nil),
	}
	ctx, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	data, err := types.ModuleCdc.MarshalJSON(genesisState)
	assert.Nil(t, err)
	am.InitGenesis(ctx, data)
	genesisbz := am.ExportGenesis(ctx)
	var genesis types.GenesisState
	err = types.ModuleCdc.UnmarshalJSON(genesisbz, &genesis)
	assert.Nil(t, err)
	assert.Equal(t, genesis, genesisState)
	am.InitGenesis(ctx, nil)
	var genesis2 types.GenesisState
	genesis2bz := am.ExportGenesis(ctx)
	err = types.ModuleCdc.UnmarshalJSON(genesis2bz, &genesis2)
	assert.Equal(t, genesis2, types.DefaultGenesisState())
}

func TestAppModule_NewQuerierHandler(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, reflect.ValueOf(keeper.NewQuerier(k)).String(), reflect.ValueOf(am.NewQuerierHandler()).String())
}

func TestAppModule_Route(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, am.Route(), types.RouterKey)
}

func TestAppModule_QuerierRoute(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, am.QuerierRoute(), types.ModuleName)
}

func TestAppModule_EndBlock(t *testing.T) {
	ctx, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, am.EndBlock(ctx, abci.RequestEndBlock{}), []abci.ValidatorUpdate{})
}

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	assert.Equal(t, []byte(am.DefaultGenesis()), []byte(types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())))
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	_, nk, ak, k := createTestInput(t, false)
	am := NewAppModule(k, nk, ak)
	p := types.Params{
		SessionNodeCount:     10,
		ProofWaitingPeriod:   22,
		SupportedBlockchains: []string{hex.EncodeToString(types.Hash([]byte("eth")))},
		ClaimExpiration:      55,
	}
	genesisState := types.GenesisState{
		Params: p,
		Proofs: []types.StoredEvidence(nil),
		Claims: []types.MsgClaim(nil),
	}
	p2 := types.Params{
		SessionNodeCount:     -1,
		ProofWaitingPeriod:   22,
		SupportedBlockchains: []string{"eth"},
		ClaimExpiration:      55,
	}
	genesisState2 := types.GenesisState{
		Params: p2,
		Proofs: []types.StoredEvidence(nil),
		Claims: []types.MsgClaim(nil),
	}
	validBz, err := types.ModuleCdc.MarshalJSON(genesisState)
	assert.Nil(t, err)
	invalidBz, err := types.ModuleCdc.MarshalJSON(genesisState2)
	assert.Nil(t, err)
	assert.True(t, nil == am.ValidateGenesis(validBz))
	assert.False(t, nil == am.ValidateGenesis(invalidBz))
}
