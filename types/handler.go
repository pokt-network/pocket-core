package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/tendermint/tendermint/state/txindex"
)

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Ctx, msg Msg, signer crypto.PublicKey) Result

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Ctx, tx Tx, txBz []byte, txIndexer txindex.TxIndexer, simulate bool) (newCtx Ctx, result Result, signer crypto.PublicKey, abort bool)
