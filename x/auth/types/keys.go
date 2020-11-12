package types

import (
	sdk "github.com/pokt-network/pocket-core/types"
)

const (
	// module name
	ModuleName = "auth"
	// storeKey is string representation of the store key for auth
	StoreKey = ModuleName
	// FeeCollectorName the root string for the fee collector account address
	FeeCollectorName = "fee_collector"
	// QuerierRoute is the querier route for auth
	QuerierRoute = StoreKey
	// default codespace
	DefaultCodespace = ModuleName
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	SupplyKeyPrefix       = []byte{0x00}
	AddressStoreKeyPrefix = []byte{0x01}
)

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.Address) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}
