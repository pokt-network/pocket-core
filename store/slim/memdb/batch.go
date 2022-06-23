package memdb

import (
	"github.com/pkg/errors"
	db "github.com/tendermint/tm-db"
)

var _ db.Batch = (*PocketMemDBBatch)(nil)

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
		opType: opTypeSet,
		key:    key,
		value:  value,
	})
}

func (p *PocketMemDBBatch) Delete(key []byte) {
	p.ops = append(p.ops, operation{
		opType: opTypeDelete,
		key:    key,
	})
}

func (p *PocketMemDBBatch) Write() error {
	for _, op := range p.ops {
		switch op.opType {
		case opTypeSet:
			_ = p.db.Set(op.key, op.value)
		case opTypeDelete:
			_ = p.db.Delete(op.key)
		default:
			return errors.Errorf("unknown operation type %v (%v)", op.opType, op)
		}
	}
	return nil
}

func (p *PocketMemDBBatch) WriteSync() error { return p.Write() }
func (p *PocketMemDBBatch) Close()           { p.ops = nil }

// util struct below

type opType int

const (
	opTypeSet opType = iota + 1
	opTypeDelete
)

type operation struct {
	opType
	key   []byte
	value []byte
}
