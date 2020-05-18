package types

import (
	"fmt"
	"time"

	sdk "github.com/pokt-network/posmint/types"
)

// Signing information of the validator is needed for tracking bad acting within the block signing process
type ValidatorSigningInfo struct {
	Address             sdk.Address `json:"address" yaml:"address"`                             // validator consensus address
	StartHeight         int64       `json:"start_height" yaml:"start_height"`                   // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64       `json:"index_offset" yaml:"index_offset"`                   // index offset into signed block bit array
	JailedUntil         time.Time   `json:"jailed_until" yaml:"jailed_until"`                   // timestamp validator cannot be unjailed until
	Tombstoned          bool        `json:"tombstoned" yaml:"tombstoned"`                       // whether or not a validator has been tombstoned (killed out of validator set)
	MissedBlocksCounter int64       `json:"missed_blocks_counter" yaml:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
	JailedBlocksCounter int64       `json:"jailed_blocks_counter" yaml:"jailed_blocks_counter"` // jailed blocks counter (to avoid scanning the array every time)
}

// Return human readable signing info
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Entropy Offset:        %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d
  Jailed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter, i.JailedBlocksCounter)
}
