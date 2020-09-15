package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
	"gopkg.in/yaml.v2"
)

// Implements Delegation interface
var _ exported.SupplyI = Supply{}

// SetTotal sets the total supply.
func (supply Supply) SetTotal(total sdk.Coins) exported.SupplyI {
	supply.Total = total
	return supply
}

// GetTotal returns the supply total.
func (supply Supply) GetTotal() sdk.Coins {
	return supply.Total
}

// NewSupply creates a new Supply instance
func NewSupply(total sdk.Coins) exported.SupplyI {
	return Supply{total}
}

// DefaultSupply creates an empty Supply
func DefaultSupply() exported.SupplyI {
	return NewSupply(sdk.NewCoins())
}

// Inflate adds coins to the total supply
func (supply Supply) Inflate(amount sdk.Coins) exported.SupplyI {
	supply.Total = supply.Total.Add(amount)
	return supply
}

// Deflate subtracts coins from the total supply
func (supply Supply) Deflate(amount sdk.Coins) exported.SupplyI {
	supply.Total = supply.Total.Sub(amount)
	return supply
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	b, err := yaml.Marshal(supply)
	if err != nil {
		fmt.Println("error converting supply to string: " + err.Error())
		return ""
	}
	return string(b)
}

// ValidateBasic validates the Supply coins and returns error if invalid
func (supply Supply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}
	return nil
}
