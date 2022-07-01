package memdb

import (
	"github.com/pkg/errors"
	"github.com/pokt-network/pocket-core/store/slim/dedup"
	db "github.com/tendermint/tm-db"
)

var _ db.Batch = (*PocketMemDBBatch)(nil)

// The batch is a simple design, the idea is to save all of the write operations into a slice and then execute them
// at once to the parent db upon the Write() call.

type PocketMemDBBatch struct {
	db  *PocketMemDB
	ops []operation
}

func NewPocketMemDBBatch(db *PocketMemDB) *PocketMemDBBatch {
	return &PocketMemDBBatch{
		db:  db,
		ops: make([]operation, 0),
	}
}

func (p *PocketMemDBBatch) Set(key, value []byte) {
	p.ops = append(p.ops, operation{
		OperationType: dedup.Set,
		key:           key,
		value:         value,
	})
}

func (p *PocketMemDBBatch) Delete(key []byte) {
	p.ops = append(p.ops, operation{
		OperationType: dedup.Del,
		key:           key,
	})
}

func (p *PocketMemDBBatch) Write() error {
	for _, op := range p.ops {
		switch op.OperationType {
		case dedup.Set:
			_ = p.db.Set(op.key, op.value)
		case dedup.Del:
			_ = p.db.Delete(op.key)
		default:
			return errors.Errorf("unknown operation type %v (%v)", op.OperationType, op)
		}
	}
	return nil
}

func (p *PocketMemDBBatch) WriteSync() error { return p.Write() }
func (p *PocketMemDBBatch) Close()           { p.ops = nil }

type operation struct {
	dedup.OperationType
	key   []byte
	value []byte
}
