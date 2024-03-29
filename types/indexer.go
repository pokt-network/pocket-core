package types

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/jordanorelli/lexnum"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/state/txindex"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	_ txindex.TxIndexer = &TransactionIndexer{}
	// Since LevelDB comparators order lexicongraphically, the implementation uses ELEN to encode numbers to ensure alphanumerical
	// ordering at insertion time. https://www.zanopha.com/docs/elen.pdf
	// Since the keys are sorted alphanumerically from the start, we don't have to:
	// (a) load all results to memory (b) paginate and sort transactions after
	// This indexer inserts in sorted order so it can paginate and return based on the db iterator
	elenEncoder = lexnum.NewEncoder('=', '-')
)

const (
	// Transaction(tx) Block height index key.
	//
	// Note: Block height is cast as int when querying using the tx height index due to the lexicographic encoder library,
	// this is safe because v0 block height should never exceed 2^63-1 on 64-bit systems or 2^31-1 on 32-bit systems.
	TxHeightKey         = "tx.height"
	TxSignerKey         = "tx.signer"
	TxRecipientKey      = "tx.recipient"
	TxHashKey           = "tx.hash"
	SortAscending       = "asc"
	SortDescending      = "desc"
	AuthCodespace       = "auth"
	sep                 = "/"
	maxPerPage          = 10000
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
func (t *TransactionIndexer) Search(ctx context.Context, q *query.Query) (res []*types.TxResult, total int, err error) {
	conditions, err := q.Conditions()
	if err != nil {
		return nil, 0, errors.Wrap(err, "error during parsing conditions from query")
	}

	if q.Pagination.Size > maxPerPage {
		q.Pagination.Size = maxPerPage
	}

	for _, condition := range conditions {
		if condition.Op != query.OpEqual {
			return nil, 0, fmt.Errorf("transaction indexer only supports op.Equal not %v", condition.Op)
		}
	}

	primaryCondition := conditions[0]
	secondaryCondition := query.Condition{}
	if len(conditions) > 1 {
		secondaryCondition = conditions[1]
		if secondaryCondition.CompositeKey != TxHeightKey {
			return nil, 0, fmt.Errorf("transaction indexer only supports secondary condition on tx.height not %v", secondaryCondition.CompositeKey)
		}
	}

	switch primaryCondition.CompositeKey {
	case TxHeightKey:
		return t.heightQuery(primaryCondition, q.Pagination)
	case TxSignerKey:
		return t.signerQuery(primaryCondition, secondaryCondition, q.Pagination)
	case TxRecipientKey:
		return t.recipientQuery(primaryCondition, secondaryCondition, q.Pagination)
	case TxHashKey:
		return t.hashQuery(primaryCondition)
	default:
		return nil, 0, fmt.Errorf("Condition.CompositeKey: %v not supported on this indexer", primaryCondition.CompositeKey)
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

func (t *TransactionIndexer) hashQuery(condition query.Condition) (res []*types.TxResult, total int, err error) {
	hash, err := hex.DecodeString(condition.Operand.(string))
	if err != nil {
		return nil, 0, errors.Wrap(err, "error during searching for a hash in the query")
	}
	result, err := t.Get(hash)
	if err == nil {
		total = 1
	}
	return []*types.TxResult{result}, total, err
}

func (t *TransactionIndexer) heightQuery(condition query.Condition, pagination *query.Page) (res []*types.TxResult, total int, err error) {
	height, ok := condition.Operand.(int64)
	if !ok {
		return nil, 0, errors.New("error during searching for a height in the query, c.Operand not type int64")
	}
	return t.getByPrefix(prefixKeyForHeight(height), pagination)
}

func (t *TransactionIndexer) signerQuery(primaryCondition query.Condition, secondaryCondition query.Condition, pagination *query.Page) (res []*types.TxResult, total int, err error) {
	signer, err := hex.DecodeString(primaryCondition.Operand.(string))
	if err != nil {
		return nil, 0, errors.Wrap(err, "error during searching for a address in the query")
	}
	if secondaryCondition.CompositeKey == TxHeightKey {
		height, ok := secondaryCondition.Operand.(int64)
		if !ok {
			return nil, 0, errors.New("error during searching for a height in the query, c.Operand not type int64")
		}
		return t.getByPrefix(prefixKeyForSignerAndHeight(signer, height), pagination)
	}
	return t.getByPrefix(prefixKeyForSigner(signer), pagination)
}

func (t *TransactionIndexer) recipientQuery(primaryCondition query.Condition, secondaryCondition query.Condition, pagination *query.Page) (res []*types.TxResult, total int, err error) {
	recipient, err := hex.DecodeString(primaryCondition.Operand.(string))
	if err != nil {
		return nil, 0, errors.Wrap(err, "error during searching for a address in the query")
	}
	if secondaryCondition.CompositeKey == TxHeightKey {
		height, ok := secondaryCondition.Operand.(int64)
		if !ok {
			return nil, 0, errors.New("error during searching for a height in the query, c.Operand not type int64")
		}
		return t.getByPrefix(prefixKeyForRecipientAndHeight(recipient, height), pagination)
	}
	return t.getByPrefix(prefixKeyForRecipient(recipient), pagination)
}

func (t *TransactionIndexer) getByPrefix(prefix []byte, pagination *query.Page) (res []*types.TxResult, total int, err error) {
	it, err := PrefixIterator(t.store, prefix, pagination.Sort)
	if err != nil {
		return nil, 0, errors.Wrap(err, "error creating prefix iterator")
	}
	defer it.Close()
	for i, skipCount := 0, 0; it.Valid(); it.Next() {
		if skipCount < pagination.Skip {
			skipCount++
			total++
			continue
		}
		val, err := t.Get(it.Value())
		if err != nil {
			return nil, 0, errors.Wrap(err, "error during query iteration get()")
		}
		if i < pagination.Size {
			res = append(res, val)
		}
		total++
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

func prefixKeyForSignerAndHeight(signer Address, height int64) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s",
		TxSignerKey,
		signer,
		elenEncoder.EncodeInt(int(height)),
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

func prefixKeyForRecipientAndHeight(recipient Address, height int64) []byte {
	return []byte(fmt.Sprintf("%s/%s/%s",
		TxRecipientKey,
		recipient,
		elenEncoder.EncodeInt(int(height)),
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
