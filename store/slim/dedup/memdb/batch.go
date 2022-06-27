package memdb

import (
	"github.com/pkg/errors"
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
		operationType: set,
		key:           key,
		value:         value,
	})
}

func (p *PocketMemDBBatch) Delete(key []byte) {
	p.ops = append(p.ops, operation{
		operationType: del,
		key:           key,
	})
}

func (p *PocketMemDBBatch) Write() error {
	for _, op := range p.ops {
		switch op.operationType {
		case set:
			_ = p.db.Set(op.key, op.value)
		case del:
			_ = p.db.Delete(op.key)
		default:
			return errors.Errorf("unknown operation type %v (%v)", op.operationType, op)
		}
	}
	return nil
}

func (p *PocketMemDBBatch) WriteSync() error { return p.Write() }
func (p *PocketMemDBBatch) Close()           { p.ops = nil }

type operationType int

const (
	set operationType = iota + 1
	del
)

type operation struct {
	operationType
	key   []byte
	value []byte
}
