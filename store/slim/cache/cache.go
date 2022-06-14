package cache

import (
	"github.com/pokt-network/pocket-core/codec"
	types2 "github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/store/types"
	dbm "github.com/tendermint/tm-db"
	"log"
	"sort"
	"strconv"
	"sync"
)

type CacheStore struct {
	cache  *CacheDB
	Height int64
	Prefix string
}

func NewStore(height int64, prefix string, cacheDB *CacheDB) *CacheStore {
	cacheDB.L.Lock()
	if _, ok := cacheDB.M[height]; !ok {
		cacheDB.M[height] = NewKVSlice()
	}
	cacheDB.L.Unlock()
	return &CacheStore{
		cache:  cacheDB,
		Height: height,
		Prefix: prefix,
	}
}

func (c *CacheStore) Get(key []byte) ([]byte, error) {
	c.cache.L.RLock()
	defer c.cache.L.RUnlock()
	kv := c.cache.M[c.Height]
	bz, _ := kv.Get(c.PrefixKey(key))
	return bz, nil
}

func (c *CacheStore) Has(key []byte) (bool, error) {
	c.cache.L.RLock()
	defer c.cache.L.RUnlock()
	kv := c.cache.M[c.Height]
	_, found := kv.Get(c.PrefixKey(key))
	return found, nil
}

func (c *CacheStore) Set(key, value []byte) {
	c.cache.L.Lock()
	defer c.cache.L.Unlock()
	kv := c.cache.M[c.Height]
	kv.Set(c.PrefixKey(key), value)
	c.cache.M[c.Height] = kv
}

func (c *CacheStore) Delete(key []byte) {
	c.cache.L.Lock()
	defer c.cache.L.Unlock()
	kv := c.cache.M[c.Height]
	kv.Delete(c.PrefixKey(key))
	c.cache.M[c.Height] = kv
}

func (c *CacheStore) Iterator(start, end []byte) (types.Iterator, error) {
	return NewIterator(c.cache, c.Height, c.Prefix, start, end, false), nil
}

func (c *CacheStore) ReverseIterator(start, end []byte) (types.Iterator, error) {
	return NewIterator(c.cache, c.Height, c.Prefix, start, end, true), nil
}

func (c *CacheStore) Commit()                     { c.Height++ }
func (c *CacheStore) PrefixKey(key []byte) string { return PrefixKey(c.Prefix, key) }

// cache db structure: the global 'cache' database that is used in multistore (outside the store abstraction)

type CacheDB struct {
	Persisted    dbm.DB             // a db to persist the cache
	LatestHeight int64              // latest height of the store
	KeepHeights  int64              // how many heights are kept in cache
	M            map[int64]*KVSlice // height -> key/value ordered slice
	L            sync.RWMutex       // thread safety
}

func (cDB *CacheDB) DeleteHeight(h int64) {
	cDB.L.Lock()
	defer cDB.L.Unlock()
	delete(cDB.M, h)
}

func (cDB *CacheDB) PrepareNextHeight(nextHeight int64) {
	cDB.L.Lock()
	defer cDB.L.Unlock()
	cDB.M[nextHeight] = cDB.M[nextHeight-1].Copy()
}

func (cDB *CacheDB) Commit() {
	bz, err := cdc.ProtoMarshalBinaryBare(&cDB.M[cDB.LatestHeight].S)
	if err != nil {
		panic(err)
	}
	if err = cDB.Persisted.Set(cDB.PreloadKey(), bz); err != nil {
		panic(err)
	}
	cDB.LatestHeight++
	cDB.PrepareNextHeight(cDB.LatestHeight)
}

func (cDB *CacheDB) Preload() {
	log.Println("preloading into cache...")
	latestHeight := cDB.LatestHeight
	if latestHeight == 0 {
		return
	}
	for i := latestHeight - cDB.KeepHeights + 1; i < latestHeight; i++ {
		cDB.LatestHeight = i
		bz, err := cDB.Persisted.Get(cDB.PreloadKey())
		if err != nil {
			panic(err)
		}
		cDB.M[i] = NewKVSlice()
		if err = cdc.ProtoUnmarshalBinaryBare(bz, &cDB.M[i].S); err != nil {
			panic(err)
		}
		cDB.M[i].size = len(cDB.M[i].S.KVSlice)
	}
	cDB.PrepareNextHeight(latestHeight)
	cDB.LatestHeight++
}

func (cDB *CacheDB) PreloadKey() []byte {
	// uses modulo to rotate slots like a circular queue
	offset := (cDB.LatestHeight % cDB.KeepHeights) + 1
	return []byte(PreloadPrefix + strconv.Itoa(int(offset)))
}

func NewCacheDB(persisted dbm.DB, latestHeight int64, keepHeights int) *CacheDB {
	return &CacheDB{
		Persisted:    persisted,
		LatestHeight: latestHeight,
		KeepHeights:  int64(keepHeights),
		M:            make(map[int64]*KVSlice),
		L:            sync.RWMutex{},
	}
}

// key value slice

type KVSlice struct {
	S    KVs
	size int
}

func NewKVSlice() *KVSlice {
	return &KVSlice{
		S:    KVs{make([]KV, 0)},
		size: 0,
	}
}

func (kv *KVSlice) Set(key string, value []byte) {
	i, found := kv.Search(key)
	if found {
		kv.S.KVSlice[i].Value = value
		return
	}
	kv.S.KVSlice = append(kv.S.KVSlice, KV{})
	copy(kv.S.KVSlice[i+1:], kv.S.KVSlice[i:])
	kv.S.KVSlice[i] = KV{
		Key:   key,
		Value: value,
	}
	kv.size++
}

func (kv *KVSlice) Get(key string) ([]byte, bool) {
	i, found := kv.Search(key)
	if !found {
		return nil, false
	}
	return kv.S.KVSlice[i].Value, true
}

func (kv *KVSlice) Delete(key string) {
	i, found := kv.Search(key)
	if !found {
		return
	}
	kv.S.KVSlice = append(kv.S.KVSlice[:i], kv.S.KVSlice[i+1:]...)
	kv.size--
}

func (kv *KVSlice) Search(key string) (int, bool) {
	i := sort.Search(kv.size, func(i int) bool {
		return kv.S.KVSlice[i].Key >= key
	})
	if i != kv.size && kv.S.KVSlice[i].Key == key {
		return i, true
	}
	return i, false
}

func (kv *KVSlice) Copy() *KVSlice {
	kvSlice := make([]KV, kv.size)
	copy(kvSlice, kv.S.KVSlice)
	return &KVSlice{
		S:    KVs{kvSlice},
		size: kv.size,
	}
}

// util

const (
	PreloadPrefix = "pl"
)

var cdc = codec.NewCodec(types2.NewInterfaceRegistry())
