package gov

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/types/module"
	"github.com/pokt-network/pocket-core/x/gov/keeper"
	"github.com/pokt-network/pocket-core/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

const moduleName = "gov"

// AppModuleBasic app module basics object
type AppModuleBasic struct{}

// Name module name
func (AppModuleBasic) Name() string {
	return moduleName
}

// RegisterCodec register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis module validate genesis
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error { return nil }

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func (am AppModule) ConsensusParamsUpdate(ctx sdk.Ctx) *abci.ConsensusParams {
	return &abci.ConsensusParams{}
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

// Name returns the staking module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {
}

func (am AppModule) UpgradeCodec(ctx sdk.Ctx) {
	am.keeper.UpgradeCodec(ctx)
}

// Route returns the message routing key for the staking module.
func (AppModule) Route() string {
	return types.RouterKey
}

// NewHandler returns an sdk.Handler for the staking module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns the staking module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// NewQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// InitGenesis performs genesis initialization for the pos module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Ctx, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	if ctx.AppVersion() == "" {
		fmt.Println(fmt.Errorf("must set app version in context, set with ctx.WithAppVersion(<version>)").Error())
		os.Exit(1)
	}
	if data == nil {
		genesisState = types.DefaultGenesisState()
	} else {
		types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	}
	return am.keeper.InitGenesis(ctx, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Ctx) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock module begin-block
func (am AppModule) BeginBlock(ctx sdk.Ctx, req abci.RequestBeginBlock) {

	am.activateAdditionalParametersACL(ctx)

	u := am.keeper.GetUpgrade(ctx)

	// We have an issue if the context version is behind the upgrade version but the context block height is ahead of the upgrade height
	if ctx.AppVersion() < u.Version && ctx.BlockHeight() >= u.UpgradeHeight() && ctx.BlockHeight() != 0 {
		ctx.Logger().Error("MUST UPGRADE TO NEXT VERSION: ", u.Version)
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventMustUpgrade, sdk.NewAttribute("VERSION:", u.UpgradeVersion())))

		ctx.Logger().Error(fmt.Sprintf("GRACEFULLY EXITING FOR UPGRADE, AT HEIGHT: %d", ctx.BlockHeight()))
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			ctx.Logger().Error(err.Error())
			os.Exit(1)
		}
		err = p.Signal(os.Interrupt)
		if err != nil {
			ctx.Logger().Error(err.Error())
			os.Exit(1)
		}
		os.Exit(2) // Send a SIGINT signal if neither of the above worked
		select {}  // why do we need this?
	}
}

// activateAdditionalParametersACL activates additional parameters on their respective upgrade heights
func (am AppModule) activateAdditionalParametersACL(ctx sdk.Ctx) {

	// Activate BlockSizeModify params
	if am.keeper.GetCodec().IsOnNamedFeatureActivationHeight(ctx.BlockHeight(), codec.BlockSizeModifyKey) {
		gParams := am.keeper.GetParams(ctx)
		// Adding a neww ACL owner for the block parameter key (before its activation)
		gParams.ACL.SetOwner(types.NewACLKey(types.PocketcoreSubspace, "BlockByteSize"), am.keeper.GetDAOOwner(ctx))
		// Update all the params
		am.keeper.SetParams(ctx, gParams)
	}

	// Activate RSCALKey params
	if am.keeper.GetCodec().IsOnNamedFeatureActivationHeight(ctx.BlockHeight(), codec.RSCALKey) {
		params := am.keeper.GetParams(ctx)
		params.ACL.SetOwner(types.NewACLKey(types.NodesSubspace, "ServicerStakeFloorMultiplier"), am.keeper.GetDAOOwner(ctx))
		params.ACL.SetOwner(types.NewACLKey(types.NodesSubspace, "ServicerStakeWeightMultiplier"), am.keeper.GetDAOOwner(ctx))
		params.ACL.SetOwner(types.NewACLKey(types.NodesSubspace, "ServicerStakeWeightCeiling"), am.keeper.GetDAOOwner(ctx))
		params.ACL.SetOwner(types.NewACLKey(types.NodesSubspace, "ServicerStakeFloorMultiplierExponent"), am.keeper.GetDAOOwner(ctx))
		am.keeper.SetParams(ctx, params)
	}
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Ctx, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
