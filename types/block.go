package types

import (
	tndmt "github.com/tendermint/tendermint/types"
)

// wrapper around tendermint block structure
type Block tndmt.Block

// wrapper around tendermints block id
type BlockID tndmt.BlockID
