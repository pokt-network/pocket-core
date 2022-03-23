/*
Package module contains application module patterns and associated "manager" functionality.
The module pattern has been broken down by:
 - independent module functionality (AppModuleBasic)
 - inter-dependent module genesis functionality (AppModuleGenesis)
 - inter-dependent module full functionality (AppModule)

inter-dependent module functionality is module functionality which somehow
depends on other modules, typically through the module keeper.  Many of the
module keepers are dependent on each other, thus in order to access the full
set of module functionality we need to define all the keepers/params-store/keys
etc. This full set of advanced functionality is defined by the AppModule interface.

Independent module functions are separated to allow for the construction of the
basic application structures required early on in the application definition
and used to enable the definition of full module functionality later in the
process. This separation is necessary, however we still want to allow for a
high level pattern for modules to follow - for instance, such that we don't
have to manually register all of the codecs for all the modules. This basic
procedure as well as other basic patterns are handled through the use of
BasicManager.

Lastly the interface for genesis functionality (AppModuleGenesis) has been
separated out from full module functionality (AppModule) so that modules which
are only used for genesis can take advantage of the Module patterns without
needlessly defining many placeholder functions
*/
package module

import (
	"encoding/json"
	"time"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//__________________________________________________________________________________________
// AppModuleBasic is the standard form for basic non-dependant elements of an application module.
type AppModuleBasic interface {
	Name() string
	RegisterCodec(*codec.Codec)

	// genesis
	DefaultGenesis() json.RawMessage
	ValidateGenesis(json.RawMessage) error
}

// collections of AppModuleBasic
type BasicManager map[string]AppModuleBasic

func NewBasicManager(modules ...AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return moduleMap
}

// RegisterCodecs registers all module codecs
func (bm BasicManager) RegisterCodec(cdc *codec.Codec) {
	for _, b := range bm {
		b.RegisterCodec(cdc)
	}
}

// Provided default genesis information for all modules
func (bm BasicManager) DefaultGenesis() map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range bm {
		genesis[b.Name()] = b.DefaultGenesis()
	}
	return genesis
}

// Provided default genesis information for all modules
func (bm BasicManager) ValidateGenesis(genesis map[string]json.RawMessage) error {
	for _, b := range bm {
		if err := b.ValidateGenesis(genesis[b.Name()]); err != nil {
			return err
		}
	}
	return nil
}

//_________________________________________________________
// AppModuleGenesis is the standard form for an application module genesis functions
type AppModuleGenesis interface {
	AppModuleBasic
	InitGenesis(sdk.Ctx, json.RawMessage) []abci.ValidatorUpdate
	ExportGenesis(sdk.Ctx) json.RawMessage
}

// AppModule is the standard form for an application module
type AppModule interface {
	AppModuleGenesis

	// registers
	RegisterInvariants(sdk.InvariantRegistry)

	// routes
	Route() string
	NewHandler() sdk.Handler
	QuerierRoute() string
	NewQuerierHandler() sdk.Querier

	BeginBlock(sdk.Ctx, abci.RequestBeginBlock)
	EndBlock(sdk.Ctx, abci.RequestEndBlock) []abci.ValidatorUpdate
	UpgradeCodec(sdk.Ctx)
}

//___________________________
// app module
type GenesisOnlyAppModule struct {
	AppModuleGenesis
}

func (gam GenesisOnlyAppModule) UpgradeCodec(sdk.Ctx) {}

// NewGenesisOnlyAppModule creates a new GenesisOnlyAppModule object
func NewGenesisOnlyAppModule(amg AppModuleGenesis) AppModule {
	return GenesisOnlyAppModule{
		AppModuleGenesis: amg,
	}
}

// register invariants
func (GenesisOnlyAppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// module message route ngame
func (GenesisOnlyAppModule) Route() string { return "" }

// module handler
func (GenesisOnlyAppModule) NewHandler() sdk.Handler { return nil }

// module querier route ngame
func (GenesisOnlyAppModule) QuerierRoute() string { return "" }

// module querier
func (gam GenesisOnlyAppModule) NewQuerierHandler() sdk.Querier { return nil }

// module begin-block
func (gam GenesisOnlyAppModule) BeginBlock(ctx sdk.Ctx, req abci.RequestBeginBlock) {}

// module end-block
func (GenesisOnlyAppModule) EndBlock(_ sdk.Ctx, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________
// module manager provides the high level utility for managing and executing
// operations for a group of modules
type Manager struct {
	Modules            map[string]AppModule
	OrderInitGenesis   []string
	OrderExportGenesis []string
	OrderBeginBlockers []string
	OrderEndBlockers   []string
}

// NewModuleManager creates a new Manager object
func NewManager(modules ...AppModule) *Manager {

	moduleMap := make(map[string]AppModule)
	var modulesStr []string
	for _, module := range modules {
		moduleMap[module.Name()] = module
		modulesStr = append(modulesStr, module.Name())
	}

	return &Manager{
		Modules:            moduleMap,
		OrderInitGenesis:   modulesStr,
		OrderExportGenesis: modulesStr,
		OrderBeginBlockers: modulesStr,
		OrderEndBlockers:   modulesStr,
	}
}

func (m Manager) GetModule(name string) *AppModule {
	appModule := m.Modules[name]
	return &appModule
}

func (m *Manager) SetModule(name string, module AppModule) {
	m.Modules[name] = module
}

// set the order of init genesis calls
func (m *Manager) SetOrderInitGenesis(moduleNames ...string) {
	m.OrderInitGenesis = moduleNames
}

// set the order of export genesis calls
func (m *Manager) SetOrderExportGenesis(moduleNames ...string) {
	m.OrderExportGenesis = moduleNames
}

// set the order of set begin-blocker calls
func (m *Manager) SetOrderBeginBlockers(moduleNames ...string) {
	m.OrderBeginBlockers = moduleNames
}

// set the order of set end-blocker calls
func (m *Manager) SetOrderEndBlockers(moduleNames ...string) {
	m.OrderEndBlockers = moduleNames
}

// register all module routes and module querier routes
func (m *Manager) RegisterInvariants(ir sdk.InvariantRegistry) {
	for _, module := range m.Modules {
		module.RegisterInvariants(ir)
	}
}

// register all module routes and module querier routes
func (m *Manager) RegisterRoutes(router sdk.Router, queryRouter sdk.QueryRouter) {
	for _, module := range m.Modules {
		if module.Route() != "" {
			router.AddRoute(module.Route(), module.NewHandler())
		}
		if module.QuerierRoute() != "" {
			queryRouter.AddRoute(module.QuerierRoute(), module.NewQuerierHandler())
		}
	}
}

// perform init genesis functionality for modules
func (m *Manager) InitGenesis(ctx sdk.Ctx, genesisData map[string]json.RawMessage) abci.ResponseInitChain {
	var validatorUpdates []abci.ValidatorUpdate
	for _, moduleName := range m.OrderInitGenesis {
		var moduleValUpdates []abci.ValidatorUpdate
		if genesisData[moduleName] == nil {
			m.Modules[moduleName].InitGenesis(ctx, nil)
			continue
		}
		err := m.Modules[moduleName].ValidateGenesis(genesisData[moduleName])
		if err != nil {
			panic(err)
		}
		moduleValUpdates = m.Modules[moduleName].InitGenesis(ctx, genesisData[moduleName])

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator InitGenesis updates already set by a previous module")
			}
			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

// perform export genesis functionality for modules
func (m *Manager) ExportGenesis(ctx sdk.Ctx) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	for _, moduleName := range m.OrderExportGenesis {
		genesisData[moduleName] = m.Modules[moduleName].ExportGenesis(ctx)
	}
	return genesisData
}

// BeginBlock performs begin block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) BeginBlock(ctx sdk.Ctx, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	defer sdk.TimeTrack(time.Now())
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	if ctx.IsOnUpgradeHeight() {
		//for _, name := range m.OrderBeginBlockers {
		//	m.Modules[name].UpgradeCodec(ctx)
		//}
		for _, mod := range m.Modules {
			mod.UpgradeCodec(ctx)
		}
	}

	for _, moduleName := range m.OrderBeginBlockers {
		m.Modules[moduleName].BeginBlock(ctx, req)
	}
	return abci.ResponseBeginBlock{
		Events: ctx.EventManager().ABCIEvents(),
	}
}

// EndBlock performs end block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) EndBlock(ctx sdk.Ctx, req abci.RequestEndBlock) abci.ResponseEndBlock {
	defer sdk.TimeTrack(time.Now())
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	validatorUpdates := []abci.ValidatorUpdate{}

	for _, moduleName := range m.OrderEndBlockers {
		moduleValUpdates := m.Modules[moduleName].EndBlock(ctx, req)

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator EndBlock updates already set by a previous module")
			}

			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}
}

// DONTCOVER
