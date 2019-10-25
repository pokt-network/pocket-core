package types

import (
	"github.com/pokt-network/pocket-core/crypto"
	tndmt "github.com/tendermint/tendermint/abci/types"
)

// wrapper around tendermints block id
type BlockID tndmt.BlockID

func (blkid BlockID) HashHex() string {
	return crypto.HexEncodeToString(blkid.Hash)
}
