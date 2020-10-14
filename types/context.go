package types

import (
	"context"
	"errors"
	"github.com/tendermint/tendermint/store"
	"time"

	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/pokt-network/pocket-core/store/gaskv"
	stypes "github.com/pokt-network/pocket-core/store/types"
)

/*
Context is an immutable object contains all information needed to
process a request.

It contains a context.Context object inside if you want to use that,
but please do not over-use it. We try to keep all data structured
and standard additions here would be better just to add to the Context struct
*/

//var _ Ctx = Context

type Context struct {
	ctx           context.Context
	ms            MultiStore
	blockstore    *store.BlockStore
	header        abci.Header
	chainID       string
	txBytes       []byte
	logger        log.Logger
	voteInfo      []abci.VoteInfo
	gasMeter      GasMeter
	blockGasMeter GasMeter
	checkTx       bool
	minGasPrice   DecCoins
	consParams    *abci.ConsensusParams
	eventManager  *EventManager
	appVersion    string
}

type Ctx interface {
	Context() context.Context
	MultiStore() MultiStore
	BlockHeight() int64
	BlockTime() time.Time
	ChainID() string
	TxBytes() []byte
	Logger() log.Logger
	VoteInfos() []abci.VoteInfo
	GasMeter() GasMeter
	BlockGasMeter() GasMeter
	IsCheckTx() bool
	MinGasPrices() DecCoins
	EventManager() *EventManager
	BlockHeader() abci.Header
	ConsensusParams() *abci.ConsensusParams
	MustGetPrevCtx(height int64) Context
	PrevCtx(height int64) (Context, error)
	WithBlockStore(bs *store.BlockStore) Context
	BlockStore() *store.BlockStore
	WithContext(ctx context.Context) Context
	WithMultiStore(ms MultiStore) Context
	WithBlockHeader(header abci.Header) Context
	WithBlockTime(newTime time.Time) Context
	WithProposer(addr Address) Context
	WithBlockHeight(height int64) Context
	WithChainID(chainID string) Context
	WithTxBytes(txBytes []byte) Context
	WithLogger(logger log.Logger) Context
	WithVoteInfos(voteInfo []abci.VoteInfo) Context
	WithGasMeter(meter GasMeter) Context
	WithBlockGasMeter(meter GasMeter) Context
	WithIsCheckTx(isCheckTx bool) Context
	WithMinGasPrices(gasPrices DecCoins) Context
	WithConsensusParams(params *abci.ConsensusParams) Context
	WithEventManager(em *EventManager) Context
	WithValue(key, value interface{}) Context
	Value(key interface{}) interface{}
	KVStore(key StoreKey) KVStore
	TransientStore(key StoreKey) KVStore
	CacheContext() (cc Context, writeCache func())
	IsZero() bool
	AppVersion() string
}

// Proposed rename, not done to avoid API breakage
type Request = Context

// Read-only accessors
func (c Context) Context() context.Context    { return c.ctx }
func (c Context) MultiStore() MultiStore      { return c.ms }
func (c Context) BlockHeight() int64          { return c.header.Height }
func (c Context) BlockTime() time.Time        { return c.header.Time }
func (c Context) ChainID() string             { return c.chainID }
func (c Context) TxBytes() []byte             { return c.txBytes }
func (c Context) Logger() log.Logger          { return c.logger }
func (c Context) VoteInfos() []abci.VoteInfo  { return c.voteInfo }
func (c Context) GasMeter() GasMeter          { return c.gasMeter }
func (c Context) BlockGasMeter() GasMeter     { return c.blockGasMeter }
func (c Context) IsCheckTx() bool             { return c.checkTx }
func (c Context) MinGasPrices() DecCoins      { return c.minGasPrice }
func (c Context) EventManager() *EventManager { return c.eventManager }
func (c Context) AppVersion() string          { return c.appVersion }

// clone the header before returning
func (c Context) BlockHeader() abci.Header {
	var msg = proto.Clone(&c.header).(*abci.Header)
	return *msg
}

func (c Context) ConsensusParams() *abci.ConsensusParams {
	return proto.Clone(c.consParams).(*abci.ConsensusParams)
}

// create a new context
func NewContext(ms MultiStore, header abci.Header, isCheckTx bool, logger log.Logger) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	return Context{
		ctx:          context.Background(),
		ms:           ms,
		header:       header,
		chainID:      header.ChainID,
		checkTx:      isCheckTx,
		logger:       logger,
		gasMeter:     stypes.NewInfiniteGasMeter(),
		minGasPrice:  DecCoins{},
		eventManager: NewEventManager(),
	}
}

func (c Context) MustGetPrevCtx(height int64) Context {
	con, err := c.PrevCtx(height)
	if err != nil {
		panic(err)
	}
	return con
}

func (c Context) PrevCtx(height int64) (Context, error) {
	if height == c.BlockHeight() {
		header := c.BlockHeader()
		if header.LastBlockId.Hash == nil {
			header.LastBlockId.Hash = c.BlockHeader().ConsensusHash
		}
		return c.WithBlockHeader(header), nil
	}
	ms := c.ms.(CommitMultiStore).CopyStore()
	err := (*ms).(CommitMultiStore).LoadVersion(height)
	if err != nil {
		return Context{}, err
	}
	meta := c.blockstore.LoadBlockMeta(height)
	if meta == nil {
		return Context{}, errors.New("block at height not found")
	}
	hash := meta.Header.LastBlockID.Hash
	if hash == nil {
		hash = meta.Header.ConsensusHash
	}
	var header = abci.Header{
		Version: abci.Version{
			Block: meta.Header.Version.Block.Uint64(),
			App:   meta.Header.Version.App.Uint64(),
		},
		ChainID:  meta.Header.ChainID,
		Height:   meta.Header.Height,
		Time:     meta.Header.Time,
		NumTxs:   meta.Header.NumTxs,
		TotalTxs: meta.Header.TotalTxs,
		LastBlockId: abci.BlockID{
			Hash: hash,
			PartsHeader: abci.PartSetHeader{
				Total: int32(meta.Header.LastBlockID.PartsHeader.Total),
				Hash:  meta.Header.Hash(),
			},
		},
		LastCommitHash:     meta.Header.LastCommitHash,
		DataHash:           meta.Header.DataHash,
		ValidatorsHash:     meta.Header.ValidatorsHash,
		NextValidatorsHash: meta.Header.NextValidatorsHash,
		ConsensusHash:      meta.Header.ConsensusHash,
		AppHash:            meta.Header.AppHash,
		LastResultsHash:    meta.Header.LastResultsHash,
		EvidenceHash:       meta.Header.EvidenceHash,
		ProposerAddress:    meta.Header.ProposerAddress,
	}
	return NewContext((*ms).(MultiStore), header, false, c.logger).WithAppVersion(c.appVersion).WithBlockStore(c.blockstore).WithConsensusParams(c.consParams), nil
}

func (c Context) WithBlockStore(bs *store.BlockStore) Context {
	c.blockstore = bs
	return c
}

func (c Context) WithAppVersion(version string) Context {
	c.appVersion = version
	return c
}

func (c Context) BlockStore() *store.BlockStore {
	return c.blockstore
}

func (c Context) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

func (c Context) WithMultiStore(ms MultiStore) Context {
	c.ms = ms
	return c
}

func (c Context) WithBlockHeader(header abci.Header) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	c.header = header
	return c
}

func (c Context) WithBlockTime(newTime time.Time) Context {
	newHeader := c.BlockHeader()
	// https://github.com/gogo/protobuf/issues/519
	newHeader.Time = newTime.UTC()
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithProposer(addr Address) Context {
	newHeader := c.BlockHeader()
	newHeader.ProposerAddress = addr.Bytes()
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithBlockHeight(height int64) Context {
	newHeader := c.BlockHeader()
	newHeader.Height = height
	return c.WithBlockHeader(newHeader)
}

func (c Context) WithChainID(chainID string) Context {
	c.chainID = chainID
	return c
}

func (c Context) WithTxBytes(txBytes []byte) Context {
	c.txBytes = txBytes
	return c
}

func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

func (c Context) WithVoteInfos(voteInfo []abci.VoteInfo) Context {
	c.voteInfo = voteInfo
	return c
}

func (c Context) WithGasMeter(meter GasMeter) Context {
	c.gasMeter = meter
	return c
}

func (c Context) WithBlockGasMeter(meter GasMeter) Context {
	c.blockGasMeter = meter
	return c
}

func (c Context) WithIsCheckTx(isCheckTx bool) Context {
	c.checkTx = isCheckTx
	return c
}

func (c Context) WithMinGasPrices(gasPrices DecCoins) Context {
	c.minGasPrice = gasPrices
	return c
}

func (c Context) WithConsensusParams(params *abci.ConsensusParams) Context {
	c.consParams = params
	return c
}

func (c Context) WithEventManager(em *EventManager) Context {
	c.eventManager = em
	return c
}

// TODO: remove???
func (c Context) IsZero() bool {
	return c.ms == nil
}

// WithValue is deprecated, provided for backwards compatibility
// Please use
//     ctx = ctx.WithContext(context.WithValue(ctx.Context(), key, false))
// instead of
//     ctx = ctx.WithValue(key, false)
func (c Context) WithValue(key, value interface{}) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}

// Value is deprecated, provided for backwards compatibility
// Please use
//     ctx.Context().Value(key)
// instead of
//     ctx.Value(key)
func (c Context) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// ----------------------------------------------------------------------------
// Store / Caching
// ----------------------------------------------------------------------------

// KVStore fetches a KVStore from the MultiStore.
func (c Context) KVStore(key StoreKey) KVStore {
	return gaskv.NewStore(c.MultiStore().GetKVStore(key), c.GasMeter(), stypes.KVGasConfig())
}

// TransientStore fetches a TransientStore from the MultiStore.
func (c Context) TransientStore(key StoreKey) KVStore {
	return gaskv.NewStore(c.MultiStore().GetKVStore(key), c.GasMeter(), stypes.TransientGasConfig())
}

// CacheContext returns a new Context with the multi-store cached and a new
// EventManager. The cached context is written to the context when writeCache
// is called.
func (c Context) CacheContext() (cc Context, writeCache func()) {
	cms := c.MultiStore().CacheMultiStore()
	cc = c.WithMultiStore(cms).WithEventManager(NewEventManager())
	return cc, cms.Write
}
