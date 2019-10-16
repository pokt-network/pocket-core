package session

import (
	"github.com/pokt-network/pocket-core/types"
)

// extension/wrapper of legacy.Blockchain for session
// TODO non native chains need to be defined by config, for now will be hash
type SessionBlockchain types.AminoBuffer

func (sbc SessionBlockchain) Validate() error {
	// todo more validation
	if len(sbc) == 0 {
		return EmptyNonNativeChainError
	}
	return nil
}
