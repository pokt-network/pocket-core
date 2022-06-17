package rootmulti

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/store/cachemulti"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	dbm "github.com/tendermint/tm-db"
	"strings"

	"github.com/pokt-network/pocket-core/store/iavl"
	"github.com/pokt-network/pocket-core/store/types"
)

const (
	latestVersionKey = "s/latest"
	commitInfoKeyFmt = "s/%d" // s/<version>
)

var cdc = codec.NewCodec(types2.NewInterfaceRegistry())

// Store is composed of many CommitStores. Name contrasts with
// cacheMultiStore which is for cache-wrapping other MultiStores. It implements
// the CommitMultiStore interface.
type Store struct {
	DB           dbm.DB
	lastCommitID types.CommitID
	storesParams map[types.StoreKey]storeParams
	stores       map[types.StoreKey]types.CommitStore
	keysByName   map[string]types.StoreKey
}

func (rs *Store) CopyStore() *types.Store {
	newParams := make(map[types.StoreKey]storeParams)
	for k, v := range rs.storesParams {
		newParams[k] = v
	}
	newStores := make(map[types.StoreKey]types.CommitStore)
	for k, v := range rs.stores {
		newStores[k] = v
	}
	newKeysByName := make(map[string]types.StoreKey)
	for k, v := range rs.keysByName {
		newKeysByName[k] = v
	}
	s := types.Store(&Store{
		DB:           rs.DB,
		lastCommitID: rs.lastCommitID,
		storesParams: newParams,
		stores:       newStores,
		keysByName:   newKeysByName,
	})
	return &s
}

var _ types.CommitMultiStore = (*Store)(nil)

func NewStore(db dbm.DB) *Store {
	return &Store{
		DB:           db,
		storesParams: make(map[types.StoreKey]storeParams),
		stores:       make(map[types.StoreKey]types.CommitStore),
		keysByName:   make(map[string]types.StoreKey),
	}
}

// Implements Store.
func (rs *Store) GetStoreType() types.StoreType {
	return types.StoreTypeMulti
}

// Implements CommitMultiStore.
func (rs *Store) MountStoreWithDB(key types.StoreKey, typ types.StoreType, db dbm.DB) {
	if key == nil {
		panic("MountIAVLStore() Key cannot be nil")
	}
	if _, ok := rs.storesParams[key]; ok {
		panic(fmt.Sprintf("CacheMultiStore duplicate store Key %v", key))
	}
	if _, ok := rs.keysByName[key.Name()]; ok {
		panic(fmt.Sprintf("CacheMultiStore duplicate store Key name %v", key))
	}
	rs.storesParams[key] = storeParams{
		key: key,
		typ: typ,
		db:  db,
	}
	rs.keysByName[key.Name()] = key
}

// Implements CommitMultiStore.
func (rs *Store) LoadLatestVersion() error {
	ver := getLatestVersion(rs.DB)
	if ver == 0 {
		// Special logic for version 0 where there is no need to get commit
		// information.
		for key, storeParams := range rs.storesParams {
			store, err := rs.loadCommitStoreFromParams(key, types.CommitID{}, storeParams)
			if err != nil {
				return fmt.Errorf("failed to load CacheMultiStore: %v", err)
			}

			rs.stores[key] = store
		}

		rs.lastCommitID = types.CommitID{}
		return nil
	}

	cInfo, err := getCommitInfo(rs.DB, ver)
	if err != nil {
		return err
	}

	// convert StoreInfos slice to map
	infos := make(map[types.StoreKey]StoreInfo)
	for _, storeInfo := range cInfo.StoreInfos {
		infos[rs.nameToKey(storeInfo.Name)] = storeInfo
	}

	// load each CacheMultiStore
	var newStores = make(map[types.StoreKey]types.CommitStore)
	for key, storeParams := range rs.storesParams {
		var id types.CommitID

		info, ok := infos[key]
		if ok {
			id = info.Core.CommitID
		}

		store, err := rs.loadCommitStoreFromParams(key, id, storeParams)
		if err != nil {
			return fmt.Errorf("failed to load CacheMultiStore: %v", err)
		}

		newStores[key] = store
	}

	rs.lastCommitID = cInfo.CommitID()
	rs.stores = newStores

	return nil
}

func (rs *Store) LoadVersion(ver int64) (*types.Store, error) {
	newStores := make(map[types.StoreKey]types.CommitStore)
	for k, v := range rs.stores {
		a, ok := (v).(*iavl.Store)
		if !ok {
			return nil, fmt.Errorf("cannot convert store into iavl store in get immutable")
		}
		s, err := a.LazyLoadStore(ver)
		if err != nil {
			return nil, fmt.Errorf("error loading store: %s, in LoadVersion: %s", k, err.Error())
		}
		newStores[k] = s
	}
	newParams := make(map[types.StoreKey]storeParams)
	for k, v := range rs.storesParams {
		newParams[k] = v
	}
	newKeysByName := make(map[string]types.StoreKey)
	for k, v := range rs.keysByName {
		newKeysByName[k] = v
	}
	s := types.Store(&Store{
		DB:           rs.DB,
		lastCommitID: rs.lastCommitID,
		storesParams: newParams,
		stores:       newStores,
		keysByName:   newKeysByName,
	})
	return &s, nil
}

//----------------------------------------
// +CommitStore

// Implements Committer/CommitStore.
func (rs *Store) LastCommitID() types.CommitID {
	return rs.lastCommitID
}

// Implements Committer/CommitStore.
func (rs *Store) Commit() types.CommitID {

	// Commit stores.
	version := rs.lastCommitID.Version + 1
	commitInfo := commitStores(version, rs.stores)

	// Need to update atomically.
	batch := rs.DB.NewBatch()
	defer batch.Close()
	setCommitInfo(batch, version, commitInfo)
	setLatestVersion(batch, version)
	_ = batch.Write()

	// Prepare for next version.
	commitID := types.CommitID{
		Version: version,
		Hash:    commitInfo.Hash(),
	}
	rs.lastCommitID = commitID
	return commitID
}

// Implements CacheWrapper/CacheMultiStore/CommitStore.
func (rs *Store) CacheWrap() types.CacheWrap {
	return rs.CacheMultiStore().(types.CacheWrap)
}

//----------------------------------------
// +MultiStore

// Implements MultiStore.
func (rs *Store) CacheMultiStore() types.CacheMultiStore {
	stores := make(map[types.StoreKey]types.CacheWrapper)
	for k, v := range rs.stores {
		stores[k] = v
	}

	return cachemulti.NewCacheMulti(rs.stores)
}

// GetKVStore implements the MultiStore interface.
// the original KVStore will be returned.
// If the store does not exist, panics.
func (rs *Store) GetKVStore(key types.StoreKey) types.KVStore {
	return rs.stores[key].(types.KVStore)
}

// Implements MultiStore

// getStoreByName will first convert the original name to
// a special Key, before looking up the CommitStore.
// This is not exposed to the extensions (which will need the
// storeKey), but is useful in main, and particularly app.Query,
// in order to convert human strings into CommitStores.
func (rs *Store) getStoreByName(name string) types.Store {
	key := rs.keysByName[name]
	if key == nil {
		return nil
	}
	return rs.stores[key]
}

//---------------------- Query ------------------
// parsePath expects a format like /<storeName>[/<subpath>]
// Must start with /, subpath may be empty
// Returns error if it doesn't start with /
func parsePath(path string) (storeName string, subpath string, err sdk.Error) {
	if !strings.HasPrefix(path, "/") {
		err = sdk.ErrUnknownRequest(fmt.Sprintf("invalid path: %s", path))
		return
	}

	paths := strings.SplitN(path[1:], "/", 2)
	storeName = paths[0]

	if len(paths) == 2 {
		subpath = "/" + paths[1]
	}

	return
}

//----------------------------------------

func (rs *Store) loadCommitStoreFromParams(key types.StoreKey, id types.CommitID, params storeParams) (store types.CommitStore, err error) {
	var db dbm.DB

	if params.db != nil {
		db = dbm.NewPrefixDB(params.db, []byte("s/_/"))
	} else {
		db = dbm.NewPrefixDB(rs.DB, []byte("s/k:"+params.key.Name()+"/"))
	}

	switch params.typ {
	case types.StoreTypeMulti:
		panic("recursive MultiStores not yet supported")

	case types.StoreTypeIAVL:
		return iavl.LoadStore(db, id)

	default:
		panic(fmt.Sprintf("unrecognized store type %v", params.typ))
	}
}

func (rs *Store) nameToKey(name string) types.StoreKey {
	for key := range rs.storesParams {
		if key.Name() == name {
			return key
		}
	}
	panic("Unknown name " + name)
}

//----------------------------------------
// storeParams

type storeParams struct {
	key types.StoreKey
	db  dbm.DB
	typ types.StoreType
}

//----------------------------------------
// commitInfo

// NOTE: Keep commitInfo a simple immutable struct.
// type commitInfo struct {
// 	// Version
// 	Version int64

// 	// CacheMultiStore info for
// 	StoreInfos []StoreInfo
// }

// Hash returns the simple merkle root hash of the stores sorted by name.
func (ci *CommitInfo) Hash() []byte {
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.Hash()
	}

	return merkle.SimpleHashFromMap(m)
}

func (ci *CommitInfo) CommitID() types.CommitID {
	return types.CommitID{
		Version: ci.Version,
		Hash:    ci.Hash(),
	}
}

//----------------------------------------
// StoreInfo

// StoreInfo contains the name and core reference for an
// underlying store.  It is the leaf of the Stores top
// level simple merkle tree.
//type StoreInfo struct {
//	Name string
//	Core StoreCore
//}
//
//type StoreCore struct {
//	// StoreType StoreType
//	CommitID types.CommitID
//	// ... maybe add more state
//}

// Implements merkle.Hasher.
func (si StoreInfo) Hash() []byte {
	// Doesn't write Name, since merkle.SimpleHashFromMap() will
	// include them via the keys.
	bz := si.Core.CommitID.Hash
	hasher := tmhash.New()

	_, err := hasher.Write(bz)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}

	return hasher.Sum(nil)
}

//----------------------------------------
// Misc.

func getLatestVersion(db dbm.DB) int64 {
	var latest sdk.Int64
	latestBytes, _ := db.Get([]byte(latestVersionKey))
	if latestBytes == nil {
		return 0
	}
	err := cdc.LegacyUnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}

	return int64(latest)
}

// Set the latest version.
func setLatestVersion(batch dbm.Batch, version int64) {
	v := sdk.Int64(version)
	latestBytes, _ := cdc.LegacyMarshalBinaryLengthPrefixed(&v)
	batch.Set([]byte(latestVersionKey), latestBytes)
}

// Commits each store and returns a new commitInfo.
func commitStores(version int64, storeMap map[types.StoreKey]types.CommitStore) CommitInfo {
	storeInfos := make([]StoreInfo, 0, len(storeMap))

	for key, store := range storeMap {
		// Commit
		commitID := store.Commit()

		// Record CommitID
		si := StoreInfo{}
		si.Name = key.Name()
		si.Core.CommitID = commitID
		// si.Core.StoreType = store.GetStoreType()
		storeInfos = append(storeInfos, si)
	}

	ci := CommitInfo{
		Version:    version,
		StoreInfos: storeInfos,
	}
	return ci
}

// Gets commitInfo from disk.
func getCommitInfo(db dbm.DB, ver int64) (CommitInfo, error) {

	// Get from DB.
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, ver)
	cInfoBytes, _ := db.Get([]byte(cInfoKey))
	if cInfoBytes == nil {
		return CommitInfo{}, fmt.Errorf("failed to get CacheMultiStore: no data")
	}

	var cInfo CommitInfo

	err := cdc.LegacyUnmarshalBinaryLengthPrefixed(cInfoBytes, &cInfo)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("failed to get CacheMultiStore: %v", err)
	}

	return cInfo, nil
}

// Set a commitInfo for given version.
func setCommitInfo(batch dbm.Batch, version int64, cInfo CommitInfo) {
	cInfoBytes, err := cdc.LegacyMarshalBinaryLengthPrefixed(&cInfo)
	if err != nil {
		panic(err)
	}
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, version)
	batch.Set([]byte(cInfoKey), cInfoBytes)
}
