package types

import "github.com/tendermint/tendermint/node"

// Handler defines the core of the state transition function of an application.
type Handler func(ctx Ctx, msg Msg) Result

// AnteHandler authenticates transactions, before their internal messages are handled.
// If newCtx.IsZero(), ctx is used instead.
type AnteHandler func(ctx Ctx, tx Tx, txBz []byte, tmNode *node.Node, simulate bool) (newCtx Ctx, result Result, abort bool)
