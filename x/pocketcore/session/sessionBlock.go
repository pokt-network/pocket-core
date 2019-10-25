package session

import (
	"github.com/pokt-network/pocket-core/types"
)

// wrapper around blockID type for session module
type SessionBlockID types.BlockID // todo keep hash as string possibly?

func (sbid SessionBlockID) Validate() error {
	// todo more header validation
	if len(sbid.Hash) == 0 {
		return EmptyBlockIDError
	}
	return nil
}
