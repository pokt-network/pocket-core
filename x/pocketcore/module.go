package pocketcore

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic "AppModuleBasic" - The fundamental building block of a sdk module
type AppModuleBasic struct{}

// Name "Name" - Returns the name of the module
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec "RegisterCodec" - Registers the codec for the module
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis "DefaultGenesis" - Returns the default genesis for the module
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis "ValidateGenesis" - Validation check for genesis state bytes
func (AppModuleBasic) ValidateGenesis(bytes json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bytes, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return types.ValidateGenesis(data)
}

// AppModule "AppModule" - The higher level building block for a module
type AppModule struct {
	AppModuleBasic               // a fundamental structure for all mods
	keeper         keeper.Keeper // responsible for store operations
}

func (am AppModule) ConsensusParamsUpdate(ctx sdk.Ctx) *abci.ConsensusParams {

	return am.keeper.ConsensusParamUpdate(ctx)
}

func (am AppModule) UpgradeCodec(ctx sdk.Ctx) {
	am.keeper.UpgradeCodec(ctx)
}

// NewAppModule "NewAppModule" - Creates a new AppModule Object
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// RegisterInvariants "RegisterInvariants" - Unused crisis checking
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route "Route" - returns the route of the module
func (am AppModule) Route() string {
	return types.RouterKey
}

// NewHandler "NewHandler" - returns the handler for the module
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute "QuerierRoute" - returns the route of the module for queries
func (am AppModule) QuerierRoute() string {
	return types.ModuleName
}

// NewQuerierHandler "NewQuerierHandler" - returns the query handler for the module
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// BeginBlock "BeginBlock" - Functionality that is called at the beginning of (every) block
func (am AppModule) BeginBlock(ctx sdk.Ctx, req abci.RequestBeginBlock) {
	ActivateAdditionalParameters(ctx, am)
	// delete the expired claims
	am.keeper.DeleteExpiredClaims(ctx)
}

// ActivateAdditionalParameters activate additional parameters on their respective upgrade heights
func ActivateAdditionalParameters(ctx sdk.Ctx, am AppModule) {
	if am.keeper.Cdc.IsOnNamedFeatureActivationHeight(ctx.BlockHeight(), codec.BlockSizeModifyKey) {
		//on the height we set the default value
		params := am.keeper.GetParams(ctx)
		params.BlockByteSize = types.DefaultBlockByteSize
		am.keeper.SetParams(ctx, params)
	}
}

// EndBlock "EndBlock" - Functionality that is called at the end of (every) block
func (am AppModule) EndBlock(ctx sdk.Ctx, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// get blocks per session
	blocksPerSession := am.keeper.BlocksPerSession(ctx)

	// run go routine because cannot access TmNode during end-block period
	go func() {
		// use this sleep timer to bypass the beginBlock lock over transactions
		minSleep := 2000
		maxSleep := 5000
		time.Sleep(time.Duration(rand.Intn(maxSleep-minSleep)+minSleep) * time.Millisecond)

		// check the consensus reactor sync status
		status, err := am.keeper.TmNode.ConsensusReactorStatus()
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("could not get status for tendermint node (cannot submit claims/proofs in this state): %s", err.Error()))
			return
		}

		if status.IsCatchingUp {
			ctx.Logger().Error("tendermint is currently syncing still (cannot submit claims/proofs in this state)")
			return
		}

		for _, node := range types.GlobalPocketNodes {
			address := node.GetAddress()
			if (ctx.BlockHeight()+int64(address[0]))%blocksPerSession == 1 && ctx.BlockHeight() != 1 {
				// auto send the proofs
				am.keeper.SendClaimTx(ctx, am.keeper, am.keeper.TmNode, node, ClaimTx)
				// auto claim the proofs
				am.keeper.SendProofTx(ctx, am.keeper.TmNode, node, ProofTx)
				// clear session cache and db
				types.ClearSessionCache(node.SessionStore)
			}
		}
	}()
	return []abci.ValidatorUpdate{}
}

// InitGenesis "InitGenesis" - Inits the module genesis from raw json
func (am AppModule) InitGenesis(ctx sdk.Ctx, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	if data == nil {
		genesisState = types.DefaultGenesisState()
	} else {
		types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	}
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis "ExportGenesis" - Exports the genesis from raw json
func (am AppModule) ExportGenesis(ctx sdk.Ctx) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}
