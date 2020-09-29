package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params    Params     `json:"params" yaml:"params"`
	DAOTokens sdk.BigInt `json:"DAO_Tokens"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, daoTokens sdk.BigInt) GenesisState {
	return GenesisState{
		Params:    params,
		DAOTokens: daoTokens,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), sdk.ZeroInt())
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.ACL == nil {
		return ErrInvalidACL(ModuleName, fmt.Errorf("nil acl"))
	}
	return nil
}
