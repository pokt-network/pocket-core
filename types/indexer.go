package types

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/jordanorelli/lexnum"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"math"
)

var (
	_           txindex.TxIndexer = &TransactionIndexer{}
	// Since LevelDB comparators order lexicongraphically, the implementation uses ELEN to encode numbers to ensure alphanumerical
	// ordering at insertion time. https://www.zanopha.com/docs/elen.pdf
	// Since the keys are sorted alphanumerically from the start, we don't have to:
	// (a) load all results to memory (b) paginate and sort transactions after
	// This indexer inserts in sorted order so it can paginate and return based on the db iterator
	elenEncoder                   = lexnum.NewEncoder('=', '-')
)

const (
	TxHeightKey         = "tx.height"
	TxSignerKey         = "tx.signer"
	TxRecipientKey      = "tx.recipient"
	TxHashKey           = "tx.hash"
	SortAscending       = "asc"
	SortDescending      = "desc"
	AuthCodespace       = "auth"
	sep                 = "/"
	maxPerPage          = 1000
	AnteHandlerMaxError = 10
)

type TransactionIndexer struct {
	store dbm.DB
}

func NewTransactionIndexer(store dbm.DB) *TransactionIndexer {
	return &TransactionIndexer{store: store}
}

func (t *TransactionIndexer) AddBatch(b *txindex.Batch) error {
	storeBatch := t.store.NewBatch()
	defer storeBatch.Close()

	for _, result := range b.Ops { // iterate through all the transaction results
		if result.Result.Codespace == AuthCodespace && result.Result.Code < AnteHandlerMaxError {
			continue // don't index any ante handler level errors
		}
		hash := result.Tx.Hash()

		// index tx by sender
		if result.Result.Signer != nil {
			storeBatch.Set(keyForSigner(result), hash)
		}

		// index tx by height
		storeBatch.Set(keyForHeight(result), hash)

		// index tx by hash
		rawBytes, err := cdc.MarshalBinaryBare(result, 0) // TODO make protobuf compatible
		if err != nil {
			return err
		}
		storeBatch.Set(hash, rawBytes)
	}

	return storeBatch.WriteSync()
}

func (t *TransactionIndexer) Index(result *types.TxResult) error {
	storeBatch := t.store.NewBatch()
	defer storeBatch.Close()
	if result.Result.Codespace == AuthCodespace && result.Result.Code < AnteHandlerMaxError {
		return nil // no indexing for ante handler level errors
	}
	hash := result.Tx.Hash()
	// index tx by sender
	if result.Result.Signer != nil {
		storeBatch.Set(keyForSigner(result), hash)
	}

	// index tx by recipient
	if result.Result.Recipient != nil {
		storeBatch.Set(keyForRecipient(result), hash)
	}

	// index tx by height
	storeBatch.Set(keyForHeight(result), hash)

	// index tx by hash
	rawBytes, err := cdc.MarshalBinaryBare(result, 0) // TODO make protobuf compatible
	if err != nil {
		return err
	}
	storeBatch.Set(hash, rawBytes)

	return storeBatch.WriteSync()
}

func (t *TransactionIndexer) Get(hash []byte) (*types.TxResult, error) {
	if len(hash) == 0 {
		return nil, txindex.ErrorEmptyHash
	}

	rawBytes, _ := t.store.Get(hash)
	if rawBytes == nil {
		return nil, nil
	}

	txResult := new(types.TxResult)
	err := cdc.UnmarshalBinaryBare(rawBytes, &txResult, 0)
	if err != nil {
		return nil, fmt.Errorf("error reading TxResult: %v", err)
	}

	return txResult, nil
}

// NOTE: Only supports op.Equal for hash, height, signer, or recipient, we only support op.Equal for simplicity and
// optimization of our use case
func (t *TransactionIndexer) Search(ctx context.Context, q *query.Query) ([]*types.TxResult, error) {
	condition, err := q.Condition()
	if err != nil {
		return nil, errors.Wrap(err, "error during parsing condition from query")
	}

	if q.Pagination.Size > maxPerPage {
		q.Pagination.Size = maxPerPage
	}

	if condition.Op != query.OpEqual {
		return nil, fmt.Errorf("transaction indexer only supports op.Equal not %v", condition.Op)
	}

	switch condition.CompositeKey {
	case TxHeightKey:
		return t.heightQuery(condition, q.Pagination)
	case TxSignerKey:
		return t.signerQuery(condition, q.Pagination)
	case TxRecipientKey:
		return t.recipientQuery(condition, q.Pagination)
	case TxHashKey:
		return t.hashQuery(condition)
	default:
		return nil, fmt.Errorf("Condition.CompositeKey: %v not supported on this indexer", condition.CompositeKey)
	}
}

func (t *TransactionIndexer) DeleteFromHeight(ctx context.Context, height int64) error {
	startKey := []byte(fmt.Sprintf("%s/%s",
		TxHeightKey,
		elenEncoder.EncodeInt(int(height)),
	))
	endKey := []byte(fmt.Sprintf("%s/%s",
		TxHeightKey,
		elenEncoder.EncodeInt(math.MaxInt64),
	))
	it, err := t.store.ReverseIterator(startKey, endKey)
	if err != nil {
		return errors.Wrap(err, "error creating the reverse iterator for deleteFromHeight")
	}
	defer it.Close()
	b := t.store.NewBatch()
	defer b.Close()
	for ; it.Valid(); it.Next() {
		b.Delete(it.Value())
	}
	return b.WriteSync()
}

func (t *TransactionIndexer) hashQuery(condition query.Condition) (res []*types.TxResult, err error) {
	hash, err := hex.DecodeString(condition.Operand.(string))
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a hash in the query")
	}
	result, err := t.Get(hash)
	return []*types.TxResult{result}, err
}

func (t *TransactionIndexer) heightQuery(condition query.Condition, pagination *query.Page) (res []*types.TxResult, err error) {
	height, ok := condition.Operand.(int64)
	if !ok {
		return nil, errors.New("error during searching for a height in the query, c.Operand not type int64")
	}
	return t.getByPrefix(prefixKeyForHeight(height), pagination)
}

func (t *TransactionIndexer) signerQuery(condition query.Condition, pagination *query.Page) (res []*types.TxResult, err error) {
	signer, err := hex.DecodeString(condition.Operand.(string))
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a address in the query")
	}
	return t.getByPrefix(prefixKeyForSigner(signer), pagination)
}

func (t *TransactionIndexer) recipientQuery(condition query.Condition, pagination *query.Page) (res []*types.TxResult, err error) {
	recipient, err := hex.DecodeString(condition.Operand.(string))
	if err != nil {
		return nil, errors.Wrap(err, "error during searching for a address in the query")
	}
	return t.getByPrefix(prefixKeyForRecipient(recipient), pagination)
}

func (t *TransactionIndexer) getByPrefix(prefix []byte, pagination *query.Page) (res []*types.TxResult, err error) {
	it, err := PrefixIterator(t.store, prefix, pagination.Sort)
	if err != nil {
		return nil, errors.Wrap(err, "error creating prefix iterator")
	}
	defer it.Close()
	for i, skipCount := 0, 0; it.Valid() && i < pagination.Size; it.Next() {
		if skipCount < pagination.Skip {
			skipCount++
			continue
		}
		val, err := t.Get(it.Value())
		if err != nil {
			return nil, errors.Wrap(err, "error during query iteration get()")
		}
		res = append(res, val)
		i++
	}
	return
}

func keyForHeight(result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s",
		TxHeightKey,
		elenEncoder.EncodeInt(int(result.Height)),
		elenEncoder.EncodeInt(int(result.Index)),
	))
}

func prefixKeyForHeight(height int64) []byte {
	return []byte(fmt.Sprintf("%s/%s/",
		TxHeightKey,
		elenEncoder.EncodeInt(int(height)),
	))
}

func keyForSigner(result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s",
		TxSignerKey,
		Address(result.Result.Signer),
		elenEncoder.EncodeInt(int(result.Height)),
		elenEncoder.EncodeInt(int(result.Index)),
	))
}

func prefixKeyForSigner(signer Address) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s",
		TxSignerKey,
		signer,
		elenEncoder.EncodeInt(0),
	))
}

func keyForRecipient(result *types.TxResult) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s/%s",
		TxRecipientKey,
		Address(result.Result.Recipient),
		elenEncoder.EncodeInt(int(result.Height)),
		elenEncoder.EncodeInt(int(result.Index)),
	))
}

func prefixKeyForRecipient(recipient Address) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s",
		TxRecipientKey,
		recipient,
		elenEncoder.EncodeInt(0),
	))
}

// contract: caller must close iterator
func PrefixIterator(db dbm.DB, prefix []byte, order string) (dbm.Iterator, error) {
	switch order {
	case SortAscending:
		return db.ReverseIterator(prefix, endKey(prefix))
	case SortDescending:
		return db.Iterator(prefix, endKey(prefix))
	default:
		return nil, fmt.Errorf("sorting order: %v not supported", order)
	}
}

func endKey(prefix []byte) []byte {
	bz := bytes.Split(prefix, []byte(sep))
	bz = append(bz[:len(bz)-1], []byte(elenEncoder.EncodeInt(math.MaxInt64)))
	return bytes.Join(bz[:], []byte(sep))
}
