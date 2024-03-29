// /*
// Package baseapp contains data structures that provide basic data storage
// functionality and act as a bridge between the ABCI interface and the SDK
// abstractions.
//
// BaseApp has no state except the CommitMultiStore you provide upon init.
// */
package baseapp

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/tendermint/tendermint/evidence"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/state/txindex"
	tmStore "github.com/tendermint/tendermint/store"
	"io"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"syscall"

	"errors"

	"github.com/gogo/protobuf/proto"

	rootMulti "github.com/pokt-network/pocket-core/store/rootmulti"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/store"
	sdk "github.com/pokt-network/pocket-core/types"
)

var ABCILogging bool
var cdc = codec.NewCodec(types.NewInterfaceRegistry())

// Key to store the consensus params in the main store.
var mainConsensusParamsKey = []byte("consensus_params")

// Enum mode for app.runTx
type runTxMode uint8

const (
	// Check a transaction
	runTxModeCheck runTxMode = iota
	// Simulate a transaction
	runTxModeSimulate runTxMode = iota
	// Deliver a transaction
	runTxModeDeliver runTxMode = iota

	// MainStoreKey is the string representation of the main store
	MainStoreKey = "main"

	codeDuplicateTransaction = 6
	authCodespace            = "auth"
)

// BaseApp reflects the ABCI application implementation.
type BaseApp struct {
	// initialized on creation
	logger           log.Logger
	name             string               // application name from abci.Info
	db               dbm.DB               // common DB backend
	tmNode           *node.Node           // <---- todo updated here
	txIndexer        txindex.TxIndexer    // <---- todo updated here
	blockstore       *tmStore.BlockStore  // <---- todo updated here
	evidencePool     *evidence.Pool       // <---- todo updated here
	cms              sdk.CommitMultiStore // Main (uncached) state
	cdc              *codec.Codec
	router           sdk.Router      // handle any kind of message
	queryRouter      sdk.QueryRouter // router for redirecting query calls
	txDecoder        sdk.TxDecoder   // unmarshal []byte into sdk.Tx
	transactionCache map[string]struct{}

	// set upon RollbackVersion or LoadLatestVersion.
	baseKey *sdk.KVStoreKey // Main KVStore in cms

	anteHandler    sdk.AnteHandler  // ante handler for fee and auth
	initChainer    sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker   sdk.BeginBlocker // logic to run before any txs
	endBlocker     sdk.EndBlocker   // logic to run after all txs, and to determine valset changes
	addrPeerFilter sdk.PeerFilter   // filter peers by address and port
	idPeerFilter   sdk.PeerFilter   // filter peers by node ID
	fauxMerkleMode bool             // if true, IAVL MountStores uses MountStoresDB for simulation speed.

	// --------------------
	// Volatile state
	// checkState is set on initialization and reset on Commit.
	// deliverState is set in InitChain and BeginBlock and cleared on Commit.
	// See methods setCheckState and setDeliverState.
	checkState   *state          // for CheckTx
	deliverState *state          // for DeliverTx
	voteInfos    []abci.VoteInfo // absent validators from begin block

	// consensus params
	// TODO: Move this in the future to baseapp param store on main store.
	consensusParams *abci.ConsensusParams

	// flag for sealing options and parameters to a BaseApp
	sealed bool

	// block height at which to halt the chain and gracefully shutdown
	haltHeight uint64

	// minimum block time (in Unix seconds) at which to halt the chain and gracefully shutdown
	haltTime uint64

	// application's version string
	appVersion string
}

var _ abci.Application = (*BaseApp)(nil)

// NewBaseApp returns a reference to an initialized BaseApp. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
//
// NOTE: The db is used to store the version number for now.
func NewBaseApp(name string, logger log.Logger, db dbm.DB, cache bool, iavlCacheSize int64, txDecoder sdk.TxDecoder, cdc *codec.Codec, options ...func(*BaseApp)) *BaseApp {
	app := &BaseApp{
		logger:           logger,
		name:             name,
		db:               db,
		cdc:              cdc,
		cms:              store.NewCommitMultiStore(db, cache, iavlCacheSize),
		router:           NewRouter(),
		queryRouter:      NewQueryRouter(),
		transactionCache: make(map[string]struct{}),
		txDecoder:        txDecoder,
		fauxMerkleMode:   false,
	}
	for _, option := range options {
		option(app)
	}

	return app
}

func (app *BaseApp) SetTendermintNode(node *node.Node) {
	app.tmNode = node
}

func (app *BaseApp) SetTxIndexer(txindexer txindex.TxIndexer) {
	app.txIndexer = txindexer
}

func (app *BaseApp) Txindexer() (txindexer txindex.TxIndexer) {
	return app.txIndexer
}

func (app *BaseApp) SetBlockstore(blockstore *tmStore.BlockStore) {
	app.blockstore = blockstore
}

func (app *BaseApp) Blockstore() (blockstore *tmStore.BlockStore) {
	return app.blockstore
}

func (app *BaseApp) SetEvidencePool(evidencePool *evidence.Pool) {
	app.evidencePool = evidencePool
}

func (app *BaseApp) EvidencePool() (evidencePool *evidence.Pool) {
	return app.evidencePool
}

// Name returns the name of the BaseApp.
func (app *BaseApp) Name() string {
	return app.name
}

// AppVersion returns the application's version string.
func (app *BaseApp) AppVersion() string {
	return app.appVersion
}

// Logger returns the logger of the BaseApp.
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

func (app *BaseApp) Store() sdk.CommitMultiStore {
	return app.cms
}

func (app *BaseApp) BlockStore() *tmStore.BlockStore {
	return app.blockstore
}

func (app *BaseApp) TMNode() *node.Node {
	return app.tmNode
}

// SetCommitMultiStoreTracer sets the store tracer on the BaseApp's underlying
// CommitMultiStore.
func (app *BaseApp) SetCommitMultiStoreTracer(w io.Writer) {
	app.cms.SetTracer(w)
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...sdk.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			if !app.fauxMerkleMode {
				app.MountStore(key, sdk.StoreTypeIAVL)
			} else {
				// StoreTypeDB doesn't do anything upon commit, and it doesn't
				// retain history, but it's useful for faster simulation.
				app.MountStore(key, sdk.StoreTypeDB)
			}
		case *sdk.TransientStoreKey:
			app.MountStore(key, sdk.StoreTypeTransient)
		default:
			fmt.Println("Unrecognized store key type " + reflect.TypeOf(key).Name())
			os.Exit(1)
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountKVStores(keys map[string]*sdk.KVStoreKey) {
	keys[sdk.ParamsKey.Name()] = sdk.ParamsKey
	for _, key := range keys {
		if !app.fauxMerkleMode {
			app.MountStore(key, sdk.StoreTypeIAVL)
		} else {
			// StoreTypeDB doesn't do anything upon commit, and it doesn't
			// retain history, but it's useful for faster simulation.
			app.MountStore(key, sdk.StoreTypeDB)
		}
	}
}

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountTransientStores(keys map[string]*sdk.TransientStoreKey) {
	keys[sdk.ParamsTKey.Name()] = sdk.ParamsTKey
	for _, key := range keys {
		app.MountStore(key, sdk.StoreTypeTransient)
	}
}

// MountStoreWithDB mounts a store to the provided key in the BaseApp
// multistore, using a specified DB.
func (app *BaseApp) MountStoreWithDB(key sdk.StoreKey, typ sdk.StoreType, db dbm.DB) {
	app.cms.MountStoreWithDB(key, typ, db)
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion(baseKey *sdk.KVStoreKey) error {
	err := app.cms.LoadLatestVersion()
	if err != nil {
		return err
	}
	return app.initFromMainStore(baseKey)
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

// initializes the remaining logic from app.cms
func (app *BaseApp) initFromMainStore(baseKey *sdk.KVStoreKey) error {
	mainStore := app.cms.GetKVStore(baseKey)
	if mainStore == nil {
		return errors.New("baseapp expects MultiStore with 'main' KVStore")
	}

	// memoize baseKey
	if app.baseKey != nil {
		return fmt.Errorf("app.baseKey expected to be nil; possible duplicate init")

	}
	app.baseKey = baseKey

	// Load the consensus params from the main store. If the consensus params are
	// nil, it will be saved later during InitChain.
	//
	// TODO: assert that InitChain hasn't yet been called.
	consensusParamsBz, _ := mainStore.Get(mainConsensusParamsKey)
	if consensusParamsBz != nil {
		var consensusParams = &abci.ConsensusParams{}

		err := proto.Unmarshal(consensusParamsBz, consensusParams)
		if err != nil {
			return err
		}

		app.setConsensusParams(consensusParams)
	}

	// needed for the export command which inits from store but never calls initchain
	app.setCheckState(abci.Header{})
	app.Seal()

	return nil
}

func (app *BaseApp) setHaltHeight(haltHeight uint64) {
	app.haltHeight = haltHeight
}

func (app *BaseApp) setHaltTime(haltTime uint64) {
	app.haltTime = haltTime
}

// Router returns the router of the BaseApp.
func (app *BaseApp) Router() sdk.Router {
	return app.router
}

// QueryRouter returns the QueryRouter of a BaseApp.
func (app *BaseApp) QueryRouter() sdk.QueryRouter { return app.queryRouter }

// Seal seals a BaseApp. It prohibits any further modifications to a BaseApp.
func (app *BaseApp) Seal() { app.sealed = true }

// IsSealed returns true if the BaseApp is sealed and false otherwise.
func (app *BaseApp) IsSealed() bool { return app.sealed }

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and Commit()
func (app *BaseApp) setCheckState(header abci.Header) { // todo <- modified here
	ms := app.cms
	context := sdk.NewContext(ms, header, true, app.logger).WithAppVersion(app.appVersion).WithBlockStore(app.blockstore)
	app.checkState = &state{
		ms:  ms.CacheMultiStore(),
		ctx: context,
	}
}

// setCheckState sets checkState with the cached multistore and
// the context wrapping it.
// It is called by InitChain() and BeginBlock(),
// and deliverState is set nil on Commit().
func (app *BaseApp) setDeliverState(header abci.Header) { // todo <- modified here
	ms := app.cms
	context := sdk.NewContext(ms, header, true, app.logger).WithAppVersion(app.appVersion).WithBlockStore(app.blockstore)
	app.deliverState = &state{
		ms:  ms.CacheMultiStore(),
		ctx: context,
	}
}

// setConsensusParams memoizes the consensus params.
func (app *BaseApp) setConsensusParams(consensusParams *abci.ConsensusParams) {
	app.consensusParams = consensusParams
}

// setConsensusParams stores the consensus params to the main store.
func (app *BaseApp) storeConsensusParams(consensusParams *abci.ConsensusParams) {
	consensusParamsBz, err := proto.Marshal(consensusParams)
	if err != nil {
		fmt.Println("error marshalling consensus params in storeConsensusParams err:" + err.Error())
		return
	}
	mainStore := app.cms.GetKVStore(app.baseKey)
	_ = mainStore.Set(mainConsensusParamsKey, consensusParamsBz)
}

// getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// if maximum block gas is less than negative one and returns zero if negative
// one.
func (app *BaseApp) getMaximumBlockGas() uint64 {
	if app.consensusParams == nil || app.consensusParams.Block == nil {
		return 0
	}

	maxGas := app.consensusParams.Block.MaxGas
	switch {
	case maxGas < -1:
		fmt.Println(fmt.Errorf("invalid maximum block gas: %d", maxGas))
		return 0
	case maxGas == -1:
		return 0

	default:
		return uint64(maxGas)
	}
}

// ----------------------------------------------------------------------------
// ABCI

// Info implements the ABCI interface.
func (app *BaseApp) Info(_ abci.RequestInfo) abci.ResponseInfo {
	lastCommitID := app.cms.LastCommitID()

	return abci.ResponseInfo{
		Data:             app.name,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}
}

// SetOption implements the ABCI interface.
func (app *BaseApp) SetOption(_ abci.RequestSetOption) (res abci.ResponseSetOption) {
	// TODO: Implement!
	return
}

// InitChain implements the ABCI interface. It runs the initialization logic
// directly on the CommitMultiStore.
func (app *BaseApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {
	// stash the consensus params in the cms main store and memoize
	if req.ConsensusParams != nil {
		app.setConsensusParams(req.ConsensusParams)
		app.storeConsensusParams(req.ConsensusParams)
	}

	initHeader := abci.Header{ChainID: req.ChainId, Time: req.Time}

	// initialize the deliver state and check state with a correct header
	app.setDeliverState(initHeader)
	app.setCheckState(initHeader)

	if app.initChainer == nil {
		return
	}

	// add block gas meter for any genesis transactions (allow infinite gas)
	app.deliverState.ctx = app.deliverState.ctx.
		WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	res = app.initChainer(app.deliverState.ctx, req)

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(res.Validators) {
			fmt.Println(fmt.Errorf(
				"len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(res.Validators)))
			os.Exit(1)
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(res.Validators))
		for i, val := range res.Validators {
			if !val.Equal(req.Validators[i]) {
				fmt.Println(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
				os.Exit(1)
			}
		}
	}

	// NOTE: We don't commit, but BeginBlock for block 1 starts from this
	// deliverState.
	return
}

// FilterPeerByAddrPort filters peers by address/port.
func (app *BaseApp) FilterPeerByAddrPort(info string) abci.ResponseQuery {
	if app.addrPeerFilter != nil {
		return app.addrPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

// FilterPeerByIDfilters peers by node ID.
func (app *BaseApp) FilterPeerByID(info string) abci.ResponseQuery {
	if app.idPeerFilter != nil {
		return app.idPeerFilter(info)
	}
	return abci.ResponseQuery{}
}

// Splits a string path using the delimiter '/'.
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func splitPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")
	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}
	return path
}

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *BaseApp) Query(req abci.RequestQuery) (res abci.ResponseQuery) {
	path := splitPath(req.Path)
	if len(path) == 0 {
		msg := "no query path provided"
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}

	switch path[0] {
	// "/app" prefix for special application queries
	case "app":
		return handleQueryApp(app, path, req)
	case "store":
		return handleQueryStore(app, path, req)

	case "p2p":
		return handleQueryP2P(app, path, req)

	case "custom":
		return handleQueryCustom(app, path, req)
	}

	msg := "unknown query path"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryApp(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	if len(path) >= 2 {
		var result sdk.Result

		switch path[1] {
		case "simulate":
			txBytes := req.Data
			tx, err := app.txDecoder(txBytes, req.Height)
			if err != nil {
				result = err.Result()
			} else {
				result = app.Simulate(txBytes, tx)
			}

		case "version":
			return abci.ResponseQuery{
				Code:      uint32(sdk.CodeOK),
				Codespace: string(sdk.CodespaceRoot),
				Height:    req.Height,
				Value:     []byte(app.appVersion),
			}

		default:
			result = sdk.ErrUnknownRequest(fmt.Sprintf("Unknown query: %s", path)).Result()
		}

		value, _ := cdc.MarshalBinaryLengthPrefixed(&result, req.Height)
		return abci.ResponseQuery{
			Code:      uint32(sdk.CodeOK),
			Codespace: string(sdk.CodespaceRoot),
			Height:    req.Height,
			Value:     value,
		}
	}

	msg := "Expected second parameter to be either simulate or version, neither was present"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryStore(app *BaseApp, path []string, req abci.RequestQuery) abci.ResponseQuery {
	// "/store" prefix for store queries
	queryable, ok := app.cms.(sdk.Queryable)
	if !ok {
		msg := "multistore doesn't support queries"
		return sdk.ErrUnknownRequest(msg).QueryResult()
	}

	req.Path = "/" + strings.Join(path[1:], "/")

	// when a client did not provide a query height, manually inject the latest
	if req.Height == 0 {
		req.Height = app.LastBlockHeight()
	}

	if req.Height <= 1 && req.Prove {
		return sdk.ErrInternal("cannot query with proof when height <= 1; please provide a valid height").QueryResult()
	}

	resp := queryable.Query(req)
	resp.Height = req.Height

	return resp
}

func handleQueryP2P(app *BaseApp, path []string, _ abci.RequestQuery) (res abci.ResponseQuery) {
	// "/p2p" prefix for p2p queries
	if len(path) >= 4 {
		cmd, typ, arg := path[1], path[2], path[3]
		switch cmd {
		case "filter":
			switch typ {
			case "addr":
				return app.FilterPeerByAddrPort(arg)

			case "id":
				return app.FilterPeerByID(arg)
			}

		default:
			msg := "Expected second parameter to be filter"
			return sdk.ErrUnknownRequest(msg).QueryResult()
		}
	}

	msg := "Expected path is p2p filter <addr|id> <parameter>"
	return sdk.ErrUnknownRequest(msg).QueryResult()
}

func handleQueryCustom(app *BaseApp, path []string, req abci.RequestQuery) (res abci.ResponseQuery) {
	// path[0] should be "custom" because "/custom" prefix is required for keeper
	// queries.
	//
	// The queryRouter routes using path[1]. For example, in the path
	// "custom/gov/proposal", queryRouter routes using "gov".
	if len(path) < 2 || path[1] == "" {
		return sdk.ErrUnknownRequest("No route for custom query specified").QueryResult()
	}

	querier := app.queryRouter.Route(path[1])
	if querier == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("no custom querier found for route %s", path[1])).QueryResult()
	}

	// when a client did not provide a query height, manually inject the latest
	if req.Height == 0 {
		req.Height = app.LastBlockHeight()
	}

	if req.Height <= 1 && req.Prove {
		return sdk.ErrInternal("cannot query with proof when height <= 1; please provide a valid height").QueryResult()
	}
	// new multistore for copy
	store, err := app.cms.(*rootMulti.Store).LoadLazyVersion(req.Height)
	if err != nil {
		return sdk.ErrInternal(
			fmt.Sprintf(
				"failed to load state at height %d; %s (latest height: %d)",
				req.Height, err, app.LastBlockHeight(),
			),
		).QueryResult()
	}
	newMS, ok := (*store).(*rootMulti.Store)
	if !ok {
		return sdk.ErrInternal(
			fmt.Sprintf(
				"failed to convert store to rootMulti at height %d; (latest height: %d)",
				req.Height, app.LastBlockHeight(),
			),
		).QueryResult()
	}
	// cache wrap the commit-multistore for safety
	ctx := sdk.NewContext(
		newMS, app.checkState.ctx.BlockHeader(), true, app.logger,
	).WithBlockStore(app.checkState.ctx.BlockStore()).WithAppVersion(app.appVersion)

	// Passes the rest of the path as an argument to the querier.
	//
	// For example, in the path "custom/gov/proposal/test", the gov querier gets
	// []string{"proposal", "test"} as the path.
	resBytes, queryErr := querier(ctx, path[2:], req)
	if queryErr != nil {
		return abci.ResponseQuery{
			Code:      uint32(queryErr.Code()),
			Codespace: string(queryErr.Codespace()),
			Height:    req.Height,
			Log:       queryErr.ABCILog(),
		}
	}

	return abci.ResponseQuery{
		Code:   uint32(sdk.CodeOK),
		Height: req.Height,
		Value:  resBytes,
	}
}

func (app *BaseApp) validateHeight(req abci.RequestBeginBlock) error {
	if req.Header.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Header.Height)
	}

	prevHeight := app.LastBlockHeight()
	if req.Header.Height != prevHeight+1 {
		return fmt.Errorf("invalid height: %d; expected: %d", req.Header.Height, prevHeight+1)
	}

	return nil
}

// BeginBlock implements the ABCI application interface.
func (app *BaseApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
	if req.Header.Height == codec.GetCodecUpgradeHeight() {
		app.cdc.SetUpgradeOverride(true)
		app.txDecoder = auth.DefaultTxDecoder(app.cdc)
	}
	if app.cms.TracingEnabled() {
		app.cms.SetTracingContext(sdk.TraceContext(
			map[string]interface{}{"blockHeight": req.Header.Height},
		))
	}

	if err := app.validateHeight(req); err != nil {
		fmt.Println(fmt.Errorf("unable to validate height for req: %v err: %s", req, err))
		os.Exit(1)
	}
	// Initialize the DeliverTx state. If this is the first block, it should
	// already be initialized in InitChain. Otherwise app.deliverState will be
	// nil, since it is reset on Commit.
	if app.deliverState == nil {
		app.setDeliverState(req.Header)
	} else {
		// In the first block, app.deliverState.ctx will already be initialized
		// by InitChain. Context is now updated with Header information.
		app.deliverState.ctx = app.deliverState.ctx.
			WithBlockHeader(req.Header).
			WithBlockHeight(req.Header.Height)
	}

	// add block gas meter
	var gasMeter sdk.GasMeter
	if maxGas := app.getMaximumBlockGas(); maxGas > 0 {
		gasMeter = sdk.NewGasMeter(maxGas)
	} else {
		gasMeter = sdk.NewInfiniteGasMeter()
	}

	app.deliverState.ctx = app.deliverState.ctx.WithBlockGasMeter(gasMeter)

	if app.beginBlocker != nil {
		res = app.beginBlocker(app.deliverState.ctx, req)
	}

	// set the signed validators for addition to context in deliverTx
	app.voteInfos = req.LastCommitInfo.GetVotes()
	return
}

// CheckTx implements the ABCI interface. It runs the "basic checks" to see
// whether or not a transaction can possibly be executed, first decoding and then
// the ante handler (which checks signatures/fees/ValidateBasic).
//
// NOTE:CheckTx does not run the actual ProtoMsg handler function(s).
func (app *BaseApp) CheckTx(req abci.RequestCheckTx) (res abci.ResponseCheckTx) {
	var result sdk.Result
	//if _, ok := app.transactionCache[TxCacheKey(req.Tx, runTxModeCheck)]; ok && && cdc.IsAfterNamedFeatureActivationHeight(app.LastBlockHeight(), "some_key_here") {
	//	return abci.ResponseCheckTx{
	//		Code:      uint32(codeDuplicateTransaction),
	//		Data:      result.Data,
	//		Log:       result.Log,
	//		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
	//		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
	//		Events:    result.Events.ToABCIEvents(),
	//	}
	//}

	tx, err := app.txDecoder(req.Tx, app.LastBlockHeight())
	if err != nil {
		result = err.Result()
	} else {
		result, _ = app.runTx(runTxModeCheck, req.Tx, tx)
	}

	return abci.ResponseCheckTx{
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
		Events:    result.Events.ToABCIEvents(),
	}
}

// DeliverTx implements the ABCI interface.
func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	var result sdk.Result
	var signerPK crypto.PublicKey
	var signer sdk.Address
	var recipient sdk.Address
	var messageType string
	var duplicateTransaction bool

	if _, ok := app.transactionCache[TxCacheKey(req.Tx, runTxModeDeliver)]; ok {
		duplicateTransaction = true
	} else {
		app.transactionCache[TxCacheKey(req.Tx, runTxModeDeliver)] = struct{}{}
	}
	tx, err := app.txDecoder(req.Tx, app.LastBlockHeight())
	if err != nil {
		result = err.Result()
	} else {
		if duplicateTransaction && cdc.IsAfterNamedFeatureActivationHeight(app.LastBlockHeight(), codec.TxCacheEnhancementKey) {
			app.logger.Debug("Duplicate Tx Found")
			result = sdk.Result{
				Code:      codeDuplicateTransaction,
				Codespace: authCodespace,
			}
		} else {
			result, signerPK = app.runTx(runTxModeDeliver, req.Tx, tx)
		}
		msg := tx.GetMsg()
		messageType = msg.Type()
		recipient = msg.GetRecipient()
		if signerPK == nil {
			if signers := msg.GetSigners(); len(signers) >= 1 {
				signer = signers[0]
			}
		} else {
			signer = sdk.Address(signerPK.Address())
		}
	}

	return abci.ResponseDeliverTx{
		Code:        uint32(result.Code),
		Data:        result.Data,
		Log:         result.Log,
		GasWanted:   int64(result.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:     int64(result.GasUsed),   // TODO: Should type accept unsigned ints?
		Events:      result.Events.ToABCIEvents(),
		Codespace:   string(result.Codespace),
		Signer:      signer,
		Recipient:   recipient,
		MessageType: messageType,
	}
}

// validateBasicTxMsgs executes basic validator calls for messages.
func validateBasicTxMsgs(msg sdk.Msg) sdk.Error {
	if msg == nil {
		return sdk.ErrUnknownRequest("Tx.GetMsg() must return at least one message")
	}
	// Validate the ProtoMsg.
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getContextForTx(mode runTxMode, txBytes []byte) (ctx sdk.Ctx) {
	ctx = app.getState(mode).ctx.
		WithTxBytes(txBytes).
		WithVoteInfos(app.voteInfos).
		WithConsensusParams(app.consensusParams)

	if mode == runTxModeSimulate {
		ctx, _ = ctx.CacheContext()
	}

	return
}

// runMsg iterates through all the messages and executes them.
// nolint: gocyclo
func (app *BaseApp) runMsg(ctx sdk.Ctx, msg sdk.Msg, mode runTxMode, signer crypto.PublicKey) (result sdk.Result) {
	var msgLogs sdk.ABCIMessageLogs

	if GetABCILogging() {
		msgLogs = make(sdk.ABCIMessageLogs, 0, 1)
	}

	var (
		data      []byte
		code      sdk.CodeType
		codespace sdk.CodespaceType
	)
	events := sdk.EmptyEvents()
	// NOTE: GasWanted is determined by ante handler and GasUsed by the GasMeter.
	// match message route
	msgRoute := msg.Route()
	handler := app.router.Route(msgRoute)
	if handler == nil {
		return sdk.ErrUnknownRequest("unrecognized ProtoMsg type: " + msgRoute).Result()
	}
	var msgResult sdk.Result
	// skip actual execution for CheckTx mode
	if mode != runTxModeCheck {
		msgResult = handler(ctx, msg, signer)
	}
	// Each message result's Data must be length prefixed in order to separate
	// each result.
	data = append(data, msgResult.Data...)
	// append events from the message's execution and a message action event
	events = events.AppendEvent(sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type())))
	events = events.AppendEvents(msgResult.Events)
	// stop execution and return on first failed message
	if !msgResult.IsOK() {
		if GetABCILogging() {
			msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint32(0), false, msgResult.Log, events))
		}
		code = msgResult.Code
		codespace = msgResult.Codespace
	}
	if GetABCILogging() {
		msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint32(0), true, msgResult.Log, events))
	}
	result = sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Log:       strings.TrimSpace(msgLogs.String()),
		GasUsed:   0,
		Events:    events,
	}

	return result
}

// Returns the applications's deliverState if app is in runTxModeDeliver,
// otherwise it returns the application's checkstate.
func (app *BaseApp) getState(mode runTxMode) *state {
	if mode == runTxModeCheck || mode == runTxModeSimulate {
		return app.checkState
	}

	return app.deliverState
}

// txContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) txContext(ctx sdk.Ctx, txBytes []byte) (
	sdk.Context, sdk.MultiStore) { // todo edit here!!!
	newMS := store.MultiStore((*app.cms.(store.CommitMultiStore).(*rootMulti.Store).CopyStore()).(*rootMulti.Store))
	if newMS.TracingEnabled() {
		newMS = newMS.SetTracingContext(
			map[string]interface{}{
				"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
			},
		)
	}
	return ctx.WithMultiStore(newMS), newMS
}

// txContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Ctx, txBytes []byte) (
	sdk.Context, sdk.CacheMultiStore) {

	ms := ctx.MultiStore()
	// TODO: https://github.com/cosmos/cosmos-sdk/issues/2824
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.SetTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

// runTx processes a transaction. The transactions is processed via an
// anteHandler. The provided txBytes may be nil in some cases, eg. in tests. For
// further details on transaction execution, reference the BaseApp SDK
// documentation.
func (app *BaseApp) runTx(mode runTxMode, txBytes []byte, tx sdk.Tx) (result sdk.Result, signer crypto.PublicKey) {
	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	var gasWanted uint64

	ctx := app.getContextForTx(mode, txBytes)
	ms := ctx.MultiStore()

	// only run the tx if there is block gas remaining
	if mode == runTxModeDeliver && ctx.BlockGasMeter().IsOutOfGas() {
		return sdk.ErrOutOfGas("no block gas left to run tx").Result(), nil
	}

	var startingGas uint64
	if mode == runTxModeDeliver {
		startingGas = ctx.BlockGasMeter().GasConsumed()
	}

	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
					rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
				)
				result = sdk.ErrOutOfGas(log).Result()
			default:
				log := fmt.Sprintf("recovered: %v\nstack:\n%v", r, string(debug.Stack()))
				result = sdk.ErrInternal(log).Result()
			}
		}

		result.GasWanted = gasWanted
		//result.GasUsed = ctx.GasMeter().GasConsumed()
		result.GasUsed = 0
	}()

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer func() {
		if mode == runTxModeDeliver {
			ctx.BlockGasMeter().ConsumeGas(
				ctx.GasMeter().GasConsumedToLimit(),
				"block gas meter",
			)

			if ctx.BlockGasMeter().GasConsumed() < startingGas {
				fmt.Println(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"}) // todo remove w/ gas
				os.Exit(1)
			}
		}
	}()
	var msgs = tx.GetMsg()
	if err := validateBasicTxMsgs(msgs); err != nil {
		return err.Result(), nil
	}

	if app.anteHandler != nil {
		anteCtx, newCtx := sdk.Ctx(sdk.Context{}), sdk.Ctx(sdk.Context{})
		abort := false
		var msCache sdk.CacheMultiStore // todo edit here

		// Cache wrap context before anteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// Ref: https://github.com/pokt-network/pocket-core/issues/2772
		//
		// NOTE: Alternatively, we could require that anteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
		// performance benefits, but it'll be more difficult to get right.
		anteCtx, msCache = app.cacheTxContext(ctx, txBytes)
		newCtx, result, signer, abort = app.anteHandler(anteCtx, tx, txBytes, app.txIndexer, mode == runTxModeSimulate)
		if newCtx != nil && !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache-wrapped, or something else
			// replaced by the ante handler. We want the original multistore, not one
			// which was cache-wrapped for the ante handler.
			//
			// Also, in the case of the tx aborting, we need to track gas consumed via
			// the instantiated gas meter in the ante handler, so we update the context
			// prior to returning.
			ctx = newCtx.WithMultiStore(ms)
		}

		gasWanted = result.GasWanted

		if abort {
			return result, signer
		}

		if mode == runTxModeDeliver {
			msCache.Write()
		}
	}

	// Create a new context based off of the existing context with a cache wrapped
	// multi-store in case message processing fails.
	runMsgCtx, newMS := app.txContext(ctx, txBytes) // todo edit here!!!
	result = app.runMsg(runMsgCtx, msgs, mode, signer)
	result.GasWanted = gasWanted

	// Safety check: don't write the cache state unless we're in DeliverTx.
	if mode != runTxModeDeliver {
		return result, signer
	}

	// only update state if all messages pass
	if result.IsOK() {
		newMS.CacheMultiStore().Write() // todo edit here!!!
	}

	return result, signer
}

// EndBlock implements the ABCI interface.
func (app *BaseApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	//if app.deliverState.ms.TracingEnabled() {
	//	app.deliverState.ms = app.deliverState.ms.SetTracingContext(nil).(sdk.CacheMultiStore)
	//} // todo edit here!!!!

	if app.endBlocker != nil {
		res = app.endBlocker(app.deliverState.ctx, req)
	}
	app.transactionCache = make(map[string]struct{})
	return
}

// Commit implements the ABCI interface. It will commit all state that exists in
// the deliver state's multi-store and includes the resulting commit ID in the
// returned abci.ResponseCommit. Commit will set the check state based on the
// latest header and reset the deliver state. Also, if a non-zero halt height is
// defined in config, Commit will execute a deferred function call to check
// against that height and gracefully halt if it matches the latest committed
// height.
func (app *BaseApp) Commit() (res abci.ResponseCommit) {
	header := app.deliverState.ctx.BlockHeader()

	var halt bool

	switch {
	case app.haltHeight > 0 && uint64(header.Height) >= app.haltHeight:
		halt = true

	case app.haltTime > 0 && header.Time.Unix() >= int64(app.haltTime):
		halt = true
	}

	if halt {
		app.halt()

		// Note: State is not actually committed when halted. Logs from Tendermint
		// can be ignored.
		return abci.ResponseCommit{}
	}

	// Write the DeliverTx state which is cache-wrapped and commit the MultiStore.
	// The write to the DeliverTx state writes all state transitions to the root
	// MultiStore (app.cms) so when Commit() is called is persists those values.
	app.deliverState.ms.Write()
	commitID := app.cms.Commit()
	app.logger.Debug("Commit synced", "commit", fmt.Sprintf("%X", commitID))

	// Reset the Check state to the latest committed.
	//
	// NOTE: This is safe because Tendermint holds a lock on the mempool for
	// Commit. Use the header from this latest block.
	app.setCheckState(header)

	// empty/reset the deliver state
	app.deliverState = nil

	return abci.ResponseCommit{
		Data: commitID.Hash,
	}
}

// halt attempts to gracefully shutdown the node via SIGINT and SIGTERM falling
// back on os.Exit if both fail.
func (app *BaseApp) halt() {
	app.logger.Info("halting node per configuration", "height", app.haltHeight, "time", app.haltTime)

	p, err := os.FindProcess(os.Getpid())
	if err == nil {
		// attempt cascading signals in case SIGINT fails (os dependent)
		sigIntErr := p.Signal(syscall.SIGINT)
		sigTermErr := p.Signal(syscall.SIGTERM)

		if sigIntErr == nil || sigTermErr == nil {
			return
		}
	}

	// Resort to exiting immediately if the process could not be found or killed
	// via SIGINT/SIGTERM signals.
	app.logger.Info("failed to send SIGINT/SIGTERM; exiting...")
	os.Exit(0)
}

// ----------------------------------------------------------------------------
// State

type state struct {
	ms  sdk.CacheMultiStore
	ctx sdk.Context
}

func (st *state) CacheMultiStore() sdk.CacheMultiStore {
	return st.ms.CacheMultiStore()
}

func (st *state) Context() sdk.Context {
	return st.ctx
}

//ABCI logging

func SetABCILogging(value bool) {
	ABCILogging = value
}

func GetABCILogging() bool {
	return ABCILogging
}

func TxCacheKey(txBytes []byte, mode runTxMode) string {
	return hex.EncodeToString(txBytes) + "/" + string(mode)
}
