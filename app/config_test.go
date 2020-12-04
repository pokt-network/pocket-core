package app

import (
	"github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	c := types.DefaultConfig("~/.pocket")

	// Check default Tx indexing params
	assert.EqualValues(t, types.DefaultTxIndexer, c.TendermintConfig.TxIndex.Indexer)
	assert.EqualValues(t, types.DefaultTxIndexTags, c.TendermintConfig.TxIndex.IndexKeys)
}
