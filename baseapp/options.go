// nolint: golint
package baseapp

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
)

// File for storing in-package BaseApp optional functions,
// for options that need access to non-exported fields of the BaseApp

// SetPruning sets a pruning option on the multistore associated with the app
func SetPruning(opts sdk.PruningOptions) func(*BaseApp) {
	return func(bap *BaseApp) { bap.cms.SetPruning(opts) }
}

// SetHaltHeight returns a BaseApp option function that sets the halt block height.
func SetHaltHeight(blockHeight uint64) func(*BaseApp) {
	return func(bap *BaseApp) { bap.setHaltHeight(blockHeight) }
}

// SetHaltTime returns a BaseApp option function that sets the halt block time.
func SetHaltTime(haltTime uint64) func(*BaseApp) {
	return func(bap *BaseApp) { bap.setHaltTime(haltTime) }
}

func (app *BaseApp) SetName(name string) {
	app.name = name
}

// SetAppVersion sets the application's version string.
func (app *BaseApp) SetAppVersion(v string) {
	app.appVersion = v
}

func (app *BaseApp) SetDB(db dbm.DB) {
	app.db = db
}

func (app *BaseApp) SetCMS(cms store.CommitMultiStore) {
	app.cms = cms
}

func (app *BaseApp) SetInitChainer(initChainer sdk.InitChainer) {
	app.initChainer = initChainer
}

func (app *BaseApp) SetBeginBlocker(beginBlocker sdk.BeginBlocker) {
	app.beginBlocker = beginBlocker
}

func (app *BaseApp) SetEndBlocker(endBlocker sdk.EndBlocker) {
	app.endBlocker = endBlocker
}

func (app *BaseApp) SetAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}

func (app *BaseApp) SetAddrPeerFilter(pf sdk.PeerFilter) {
	app.addrPeerFilter = pf
}

func (app *BaseApp) SetIDPeerFilter(pf sdk.PeerFilter) {
	app.idPeerFilter = pf
}

func (app *BaseApp) SetFauxMerkleMode() {
	app.fauxMerkleMode = true
}
