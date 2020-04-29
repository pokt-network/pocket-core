package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig("~/.pocket")
	// Check default dbbackend
	assert.EqualValues(t, DefaultDBBackend, c.TendermintConfig.DBBackend)

	// Check default Tx indexing params
	assert.EqualValues(t, DefaultTxIndexer, c.TendermintConfig.TxIndex.Indexer)
	assert.EqualValues(t, DefaultTxIndexTags, c.TendermintConfig.TxIndex.IndexTags)
}
