package types

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const (
	tagKeySeparator = "/"
)

// TxIndex is the simplest possible indexer, backed by key-value storage (levelDB).
type TxIndex struct {
	store                dbm.DB
	compositeKeysToIndex []string
	indexAllEvents       bool
	cdc                  *amino.Codec
	index                uint32
	height               int64
	cache                *Cache
	indexMap             *sync.Map
}

type keyAndHash struct {
	key  string
	hash []byte
}

// NewTxIndex creates new KV indexer.
func NewTxIndex(store dbm.DB, cdc *amino.Codec, cacheSize int, options ...func(*TxIndex)) *TxIndex {
	txi := &TxIndex{store: store, compositeKeysToIndex: make([]string, 0), indexAllEvents: false, index: 0, cdc: cdc, cache: NewCache(cacheSize), indexMap: &sync.Map{}}
	for _, option := range options {
		option(txi)
	}
	return txi
}
func (txi *TxIndex) UpdateHeight(height int64) {
	if txi.height != height {
		// Reset cache & indexes
		txi.cache.Purge()
		txi.indexMap = nil // releases the map to OS NOTE: consider using an unsafe pointer and manually releasing
		txi.indexMap = &sync.Map{}
	}
	txi.height = height
}

// Get gets transaction from the TxIndex storage and returns it or nil if the
// transaction is not found.
func (txi *TxIndex) Get(hash []byte) (*tmTypes.TxResult, error) {
	if len(hash) == 0 {
		return nil, txindex.ErrorEmptyHash
	}
	if txRes, ok := txi.cache.Get(string(hash)); ok {
		txResult := new(tmTypes.TxResult)
		err := txi.cdc.UnmarshalBinaryBare(txRes.([]byte), &txResult)
		if err != nil {
			return nil, fmt.Errorf("error reading TxResult: %v", err)
		}
		return txResult, nil
	}

	rawBytes, _ := txi.store.Get(hash)
	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(tmTypes.TxResult)
	err := txi.cdc.UnmarshalBinaryBare(rawBytes, &txResult)
	if err != nil {
		return nil, fmt.Errorf("error reading TxResult: %v", err)
	}
	txi.cache.Add(string(hash), rawBytes)

	return txResult, nil
}

// IndexEvents is an option for setting which composite keys to index.
func IndexEvents(compositeKeys []string) func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.compositeKeysToIndex = compositeKeys
	}
}

// IndexAllEvents is an option for indexing all events.
func IndexAllEvents() func(*TxIndex) {
	return func(txi *TxIndex) {
		txi.indexAllEvents = true
	}
}
func (txi *TxIndex) IndexWithSender(result *tmTypes.TxResult, sender string) error {
	b := txi.store.NewBatch()
	defer b.Close()

	// index tx by events
	txi.indexEvents(result, result.Tx.Hash(), b, sender)

	// index tx by height
	if txi.indexAllEvents || stringInSlice(tmTypes.TxHeightKey, txi.compositeKeysToIndex) {
		b.Set(keyForHeight(result), result.Tx.Hash())
	}
	// index tx by hash
	rawBytes, err := txi.cdc.MarshalBinaryBare(result)
	if err != nil {
		return err
	}

	b.Set(result.Tx.Hash(), rawBytes)
	b.WriteSync()

	txi.index++
	return nil
}

func (txi *TxIndex) Index(result *tmTypes.TxResult) error {
	return nil
}

func (txi *TxIndex) indexEvents(result *tmTypes.TxResult, hash []byte, store dbm.SetDeleter, sender string) {
	for _, event := range result.Result.Events {
		// only index events with a non-empty type
		if len(event.Type) == 0 {
			continue
		}

		for _, attr := range event.Attributes {
			if len(attr.Key) == 0 {
				continue
			}

			compositeTag := fmt.Sprintf("%s.%s", event.Type, string(attr.Key))
			if txi.indexAllEvents || stringInSlice(compositeTag, txi.compositeKeysToIndex) {
				if compositeTag == "message.sender" && string(attr.Value) != sender {
					continue // Index value does not match sender cannot index !!
				}
				store.Set(keyForEventBytes(compositeTag, attr.Value, result, txi.indexMap), hash)
			}
		}
	}
}

func (txi *TxIndex) ReducedSearch(ctx context.Context, q *query.Query) ([]*types.TxResult, error) {
	return txi.Search(ctx, q)
}

// Search performs a search using a reduced query operations
// CONTRACT: will only look for single condition (Eq, Contains) will not look for ranges.
func (txi *TxIndex) Search(ctx context.Context, q *query.Query) ([]*types.TxResult, error) {
	// get a list of conditions (like "tx.height > 5")
	condition, err := q.Condition()
	if err != nil {
		return nil, errors.Wrap(err, "error during parsing condition from query")
	}

	// if there is a hash condition, return the result immediately
	hash, ok, err := lookForHash(condition)
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a hash in the query")
	} else if ok {
		res, err := txi.Get(hash)
		switch {
		case err != nil:
			return []*types.TxResult{}, errors.Wrap(err, "error while retrieving the result")
		case res == nil:
			return []*types.TxResult{}, nil
		default:
			return []*types.TxResult{res}, nil
		}
	}

	height := lookForHeight(condition)

	matchedKeys := txi.keys(ctx, condition, startKeyForCondition(condition, height), q.Pagination)

	switch q.Pagination.Sort {
	case "asc":
		sort.Slice(matchedKeys, func(i, j int) bool {
			a := strings.Split(matchedKeys[i].key, "/")
			b := strings.Split(matchedKeys[j].key, "/")
			aHeight, _ := strconv.Atoi(a[2])
			bHeight, _ := strconv.Atoi(b[2])
			if aHeight == bHeight {
				aIndex, _ := strconv.Atoi(a[3])
				bIndex, _ := strconv.Atoi(b[3])
				return aIndex < bIndex
			}
			return aHeight < bHeight
		})
	default:
		sort.Slice(matchedKeys, func(i, j int) bool {
			a := strings.Split(matchedKeys[i].key, "/")
			b := strings.Split(matchedKeys[j].key, "/")
			aHeight, _ := strconv.Atoi(a[2])
			bHeight, _ := strconv.Atoi(b[2])
			if aHeight == bHeight {
				aIndex, _ := strconv.Atoi(a[3])
				bIndex, _ := strconv.Atoi(b[3])
				return aIndex > bIndex
			}
			return aHeight > bHeight
		})
	}
	skipCount := 0
	results := make([]*types.TxResult, 0, q.Pagination.Size)
	for _, key := range matchedKeys {
		// skip keys
		if skipCount < q.Pagination.Skip {
			skipCount++
			continue
		}
		res, err := txi.Get(key.hash)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get Tx{%X}", key)
		}
		results = append(results, res)
		if len(results) == cap(results) {
			return results, nil
		}
	}
	return results, nil
}

// Retrieves the keys from the iterator based on condition
// NOTE: filteredHashes may be empty if no previous condition has matched.
func (txi *TxIndex) keys(
	ctx context.Context,
	c query.Condition,
	startKeyBz []byte,
	pagination *query.Page,
) []*keyAndHash {
	hashes := make([]*keyAndHash, 0)
	switch {
	case c.Op == query.OpEqual:
		var it dbm.Iterator
		switch pagination.Sort {
		case "asc":
			it, _ = dbm.IteratePrefix(txi.store, startKeyBz)
		case "desc", "":
			it, _ = reverseIteratePrefix(txi.store, startKeyBz)
		}
		defer it.Close()
		for ; it.Valid(); it.Next() {
			// Potentially exit early.
			select {
			case <-ctx.Done():
				break
			default:
				// hashes[string(it.Key())] = it.Value()
				hashes = append(hashes, &keyAndHash{key: string(it.Key()), hash: it.Value()})
			}

		}
	case c.Op == query.OpContains:
		// XXX: startKey does not apply here.
		// For example, if startKey = "account.owner/an/" and search query = "account.owner CONTAINS an"
		// we can't iterate with prefix "account.owner/an/" because we might miss keys like "account.owner/Ulan/"

		var it dbm.Iterator
		switch pagination.Sort {
		case "asc":
			it, _ = dbm.IteratePrefix(txi.store, startKey(c.CompositeKey))
		case "desc", "":
			it, _ = reverseIteratePrefix(txi.store, startKey(c.CompositeKey))
		}
		defer it.Close()

		skipCount := 0
		for ; it.Valid(); it.Next() {
			// Potentially exit early.
			select {
			case <-ctx.Done():
				break
			default:
				if skipCount < pagination.Skip {
					skipCount++
					continue
				}
				if !isTagKey(it.Key()) {
					continue
				}

				if strings.Contains(extractValueFromKey(it.Key()), c.Operand.(string)) {
					hashes = append(hashes, &keyAndHash{key: string(it.Key()), hash: it.Value()})
				}
			}
		}
	default:
		panic("other operators should be handled already")
	}
	return hashes
}

// NOTE: satisfy tx indexer interface not intended for cusotm usage
func (txi *TxIndex) AddBatch(b *txindex.Batch) error {
	return nil
}

// NOTE: satisfy tx indexer interface not intended for cusotm usage
func (txi *TxIndex) DeleteFromHeight(ctx context.Context, height int64) error {
	return nil
}

func keyForEventBytes(key string, value []byte, result *tmTypes.TxResult, m *sync.Map) []byte {
	base := fmt.Sprintf("%s/%s/%d", key, value, result.Height)
	index, loaded := m.LoadOrStore(base, uint32(0))
	idx := index.(uint32)
	switch loaded {
	case false:
		result.Index = idx
		m.Store(base, idx) // make sure to increase
	default:
		result.Index = idx + 1
		m.Store(base, result.Index) // make sure to increase
	}
	keyForEvent := fmt.Sprintf("%s/%d",
		base,
		result.Index,
	)
	return []byte(keyForEvent)
}

func keyForHeight(result *tmTypes.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%d/%d/%d",
		tmTypes.TxHeightKey,
		result.Height,
		result.Height,
		result.Index,
	))
}

func Hash(tx []byte) []byte {
	return tmhash.Sum(tx)
}

func stringInSlice(s string, sl []string) bool {
	for _, a := range sl {
		if a == s {
			return true
		}
	}
	return false
}

func reverseIteratePrefix(db dbm.DB, prefix []byte) (dbm.Iterator, error) {
	var start, end []byte
	if len(prefix) == 0 {
		start = nil
		end = nil
	} else {
		start = cp(prefix)
		end = cpIncr(prefix)
	}
	itr, err := db.ReverseIterator(start, end)
	if err != nil {
		return nil, err
	}
	return itr, nil
}

func cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func cpIncr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpIncr expects non-zero bz length")
	}
	ret = cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {
			// Overflow
			return nil
		}
	}
	return nil
}

func startKey(fields ...interface{}) []byte {
	var b bytes.Buffer
	for _, f := range fields {
		b.Write([]byte(fmt.Sprintf("%v", f) + tagKeySeparator))
	}
	return b.Bytes()
}

func isTagKey(key []byte) bool {
	return strings.Count(string(key), tagKeySeparator) == 3
}

func extractValueFromKey(key []byte) string {
	parts := strings.SplitN(string(key), tagKeySeparator, 3)
	return parts[1]
}

func lookForHash(conditions ...query.Condition) (hash []byte, ok bool, err error) {
	for _, c := range conditions {
		if c.CompositeKey == types.TxHashKey {
			decoded, err := hex.DecodeString(c.Operand.(string))
			return decoded, true, err
		}
	}
	return
}

// lookForHeight returns a height if there is an "height=X" condition.
func lookForHeight(conditions ...query.Condition) (height int64) {
	for _, c := range conditions {
		if c.CompositeKey == types.TxHeightKey && c.Op == query.OpEqual {
			return c.Operand.(int64)
		}
	}
	return 0
}

func startKeyForCondition(c query.Condition, height int64) []byte {
	if height > 0 {
		return startKey(c.CompositeKey, c.Operand, height)
	}
	return startKey(c.CompositeKey, c.Operand)
}
