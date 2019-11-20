package types

import (
	"fmt"

	sdk "github.com/pokt-network/posmint/types"
)

// names used as root for pool module accounts:
// - StakingPool -> "application_staked_tokens_pool"
const (
	StakedPoolName = "application_staked_tokens_pool"
)

type Pool struct {
	Tokens sdk.Int
}

// Tokens - tracking staked token supply
type StakingPool Pool

// NewPool creates a new Tokens instance used for queries
func NewPool(tokens sdk.Int) Pool {
	return Pool{
		Tokens: tokens,
	}
}

// String returns a human readable string representation of a pool.
func (bp StakingPool) String() string {
	return fmt.Sprintf(`Staked Tokens:      %s`, bp.Tokens)
}
