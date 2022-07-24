package slim

import (
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/store/slim/cache"
	"github.com/pokt-network/pocket-core/store/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	dbm "github.com/tendermint/tm-db"
)

var cdc = codec.NewCodec(types2.NewInterfaceRegistry())

const (
	latestVersionKey = "s/latest"
	commitInfoKeyFmt = "s/%d"
)

// NOTE: almost all of this is legacy cosmos sdk

var _ types.KVStore = &Store{}
var _ types.CommitStore = &Store{}

func multiStoreToStore(db dbm.DB, cachedb *cache.CacheDB, lastcommit types.CommitID, newStores map[types.StoreKey]types.CommitStore) *types.Store {
	newMultiStore := types.Store(&MultiStore{
		DB:         db,
		CacheDB:    cachedb,
		Stores:     newStores,
		LastCommit: lastcommit,
	})
	return &newMultiStore
}

func getLatestVersion(db dbm.DB) int64 {
	var latest sdk.Int64
	latestBytes, _ := db.Get([]byte(latestVersionKey))
	if latestBytes == nil {
		return 0
	}
	err := cdc.ProtoUnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}
	return int64(latest)
}

func setLatestVersion(batch dbm.Batch, version int64) {
	v := sdk.Int64(version)
	latestBytes, _ := cdc.ProtoMarshalBinaryLengthPrefixed(&v)
	batch.Set([]byte(latestVersionKey), latestBytes)
}

func getCommitInfo(db dbm.DB, ver int64) (cInfo CommitInfo, err error) {
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, ver)
	cInfoBytes, _ := db.Get([]byte(cInfoKey))
	if cInfoBytes == nil {
		return CommitInfo{}, fmt.Errorf("failed to get CacheMultiStore: no data")
	}
	err = cdc.ProtoUnmarshalBinaryLengthPrefixed(cInfoBytes, &cInfo)
	if err != nil {
		return CommitInfo{}, fmt.Errorf("failed to get CacheMultiStore: %v", err)
	}
	return cInfo, nil
}

func setCommitInfo(batch dbm.Batch, version int64, cInfo CommitInfo) {
	cInfoBytes, err := cdc.ProtoMarshalBinaryLengthPrefixed(&cInfo)
	if err != nil {
		panic(err)
	}
	cInfoKey := fmt.Sprintf(commitInfoKeyFmt, version)
	batch.Set([]byte(cInfoKey), cInfoBytes)
}

func getCommitID(db dbm.DB, ver int64) (cID types.CommitID, err error) {
	cInfo, err := getCommitInfo(db, ver)
	if err != nil {
		return types.CommitID{}, err
	}
	return cInfo.CommitID(), nil
}

func (ci *CommitInfo) CommitID() types.CommitID {
	return types.CommitID{
		Version: ci.Version,
		Hash:    ci.Hash(),
	}
}

func (ci *CommitInfo) Hash() []byte {
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.Hash()
	}
	return merkle.SimpleHashFromMap(m)
}

func (si StoreInfo) Hash() []byte {
	bz := si.Core.CommitID.Hash
	hasher := tmhash.New()
	_, err := hasher.Write(bz)
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}
