package types

import (
	"fmt"

	sdk "github.com/pokt-network/posmint/types"
)

// names used as root for pool module accounts:
// - StakingPool -> "staked_tokens_pool"
const (
	DAOPoolName    = "dao_pool"
	StakedPoolName = "staked_tokens_pool"
)

type Pool struct {
	Tokens sdk.Int
}

// Tokens - tracking staked token supply
type StakingPool Pool

// Tokens - tracking dao token supply
type DAOPool Pool

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

// String returns a human readable string representation of a pool.
func (dao DAOPool) String() string {
	return fmt.Sprintf(`Tokens:	
  Tokens In DAO Tokens:      %s`,
		dao.Tokens)
}
