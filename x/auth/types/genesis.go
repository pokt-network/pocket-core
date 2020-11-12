package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params   Params    `json:"params" yaml:"params"`
	Accounts Accounts  `json:"accounts" yaml:"accounts"`
	Supply   sdk.Coins `json:"supply" yaml:"supply"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, accounts Accounts, supply sdk.Coins) GenesisState {
	return GenesisState{
		Params:   params,
		Accounts: accounts,
		Supply:   supply,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), nil, nil)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, account := range data.Accounts {
		if account.GetPubKey().PubKey() == nil {
			return fmt.Errorf("PubKey should never be nil")
		}
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if err := NewSupply(data.Supply).ValidateBasic(); err != nil {
		return err
	}
	return nil
}
