package types

import (
	"context"
	"fmt"
	"testing"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

var memCdc *codec.Codec

func TestTxIndex(t *testing.T) {
	allowedKeys := []string{"account.number", "account.owner", "account.date", "message.sender"}
	memCdc = amino.NewCodec()
	indexer := NewTxIndex(db.NewMemDB(), memCdc, 10, IndexEvents(allowedKeys))

	tx := tmTypes.Tx("HELLO WORLD")
	txResult := &tmTypes.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx,
		Result: abci.ResponseDeliverTx{
			Data: []byte{0},
			Code: abci.CodeTypeOK, Log: "", Events: nil,
		},
	}
	hash := tx.Hash()

	err := indexer.IndexWithSender(txResult, "a signer")
	require.NoError(t, err)

	loadedTxResult, err := indexer.Get(hash)
	require.NoError(t, err)
	assert.Equal(t, txResult, loadedTxResult)

	loadedTxResult, err = indexer.Get(hash) // Make sure cache works
	require.NoError(t, err)
	assert.Equal(t, txResult, loadedTxResult)
}

func TestTxSearch(t *testing.T) {
	memCdc = amino.NewCodec()
	allowedKeys := []string{"account.number", "account.owner", "account.date", "message.sender"}
	indexer := NewTxIndex(db.NewMemDB(), memCdc, 10, IndexEvents(allowedKeys))

	txResult := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 1", 1, 0)

	txResult2 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 2", 1, 1)

	txResult3 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("2")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "message", Attributes: []kv.Pair{{Key: []byte("sender"), Value: []byte("address")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 3", 2, 0)

	txResult4 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 4", 3, 0)

	txResult5 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 5", 3, 1)

	txResult6 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 6", 5, 0)

	txResult7 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 7", 5, 1)

	txResult8 := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("number"), Value: []byte("1")}}},
		{Type: "account", Attributes: []kv.Pair{{Key: []byte("owner"), Value: []byte("Ivan")}}},
		{Type: "", Attributes: []kv.Pair{{Key: []byte("not_allowed"), Value: []byte("Vlad")}}},
	}, "Hello world 8", 5, 2)

	err := indexer.IndexWithSender(txResult, "1")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult2, "1")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult3, "a third signer")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult4, "still signer")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult5, "a signer")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult6, "1")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult7, "1")
	require.NoError(t, err)
	err = indexer.IndexWithSender(txResult8, "1")
	require.NoError(t, err)

	tt := []struct {
		name          string
		q             string
		resultsLength int
		pagination    bool
		size          int
		skip          int
		sort          string
		orderedResult []*types.TxResult
	}{
		{name: "search by hash", q: fmt.Sprintf("tx.hash = '%X'", hash), resultsLength: 1},
		{
			name:          "search by exact match",
			q:             "account.number = 1",
			resultsLength: 7,
			pagination:    true,
			size:          7,
			skip:          0,
			sort:          "asc",
			orderedResult: []*types.TxResult{txResult, txResult2, txResult4, txResult5, txResult6, txResult7, txResult8},
		},
		{
			name:          "search by exact match with page limit",
			q:             "account.number = 1",
			resultsLength: 2,
			pagination:    true,
			size:          2,
			skip:          0,
			sort:          "asc",
			orderedResult: []*types.TxResult{txResult, txResult2},
		},
		{
			name:          "search by exact match with descending order",
			q:             "account.number = 1",
			resultsLength: 4,
			pagination:    true,
			size:          4,
			skip:          0,
			sort:          "desc",
			orderedResult: []*types.TxResult{txResult8, txResult7, txResult6, txResult5},
		},
		{
			name:          "cannot find txResult for miss indexed tx (signer != sender)",
			q:             "message.sender = 'address'",
			resultsLength: 0,
			pagination:    true,
			size:          5,
			skip:          0,
			sort:          "desc",
			orderedResult: []*types.TxResult{},
		},
	}
	for _, tc := range tt {
		ctx := context.Background()
		t.Run(tc.name, func(t *testing.T) {
			q := query.MustParse(tc.q)
			if tc.pagination {
				q.Pagination = &query.Page{tc.size, tc.skip, tc.sort}
			}
			results, err := indexer.Search(ctx, q)
			assert.NoError(t, err)

			assert.Len(t, results, tc.resultsLength)
			if len(results) > 1 {
				assert.Equal(t, tc.orderedResult, results)
			}
		})
	}
}

func txResultWithEvents(events []abci.Event, txMsg string, height int64, index uint32) *types.TxResult {
	tx := types.Tx(txMsg)
	return &types.TxResult{
		Height: height,
		Index:  index,
		Tx:     tx,
		Result: abci.ResponseDeliverTx{
			Data:   []byte{0},
			Code:   abci.CodeTypeOK,
			Log:    "",
			Events: events,
		},
	}
}

// func IndexEvents(compositeKeys []string) func(*TxIndex) {
// 	return func(txi *TxIndex) {
// 		txi.compositeKeysToIndex = compositeKeys
// 	}
// }
