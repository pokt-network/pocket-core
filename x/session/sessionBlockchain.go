package session

import (
	"github.com/pokt-network/pocket-core/types"
)

// extension/wrapper of legacy.Blockchain for session
// TODO should be refactored to remove legacy
type SessionBlockchain types.AminoBuffer
