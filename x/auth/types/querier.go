package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

// query endpoints supported by the auth Querier
const (
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address sdk.Address
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr sdk.Address) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
