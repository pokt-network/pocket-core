package store

import (
	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/store/rootmulti"
	"github.com/pokt-network/pocket-core/store/types"
)

func NewCommitMultiStore(db dbm.DB, cache bool, iavlCacheSize int64) types.CommitMultiStore {
	return rootmulti.NewStore(db, cache, iavlCacheSize)
}
