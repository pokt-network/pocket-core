package types

import (
	"encoding/binary"
	sdk "github.com/pokt-network/posmint/types"
	"time"
)

const (
	ModuleName = "application"
	StoreKey   = ModuleName // StoreKey is the string store representation
	// TStoreKey is the string transient store representation
	TStoreKey    = "transient_" + ModuleName
	QuerierRoute = ModuleName // QuerierRoute is the querier route for the staking module
	RouterKey    = ModuleName // RouterKey is the msg router key for the staking module
)

//nolint
var ( // Keys for store prefixes
	AllApplicationsKey                = []byte{0x01} // prefix for each key to a application
	AllApplicationsByConsensusAddrKey = []byte{0x02} // prefix for each key to a application index, by pubkey
	StakedAppsKey                     = []byte{0x03} // prefix for each key to a staked application index, sorted by power
	UnstakingAppsKey                  = []byte{0x04} // prefix for unstaking application
	UnstakedAppsKey                   = []byte{0x05} // prefix for unstaked applications
	BurnApplicationKey                = []byte{0x06} // prefix for awarding applications
)

// Removes the prefix bytes from a key to expose true address
func AddressFromKey(key []byte) []byte {
	return key[1:] // remove prefix bytes
}

// generates the key for the application with address
func KeyForAppByAllApps(addr sdk.ValAddress) []byte {
	return append(AllApplicationsKey, addr.Bytes()...)
}

// generates the key for the application with consensus address
func KeyForAppByConsAddr(addr sdk.ConsAddress) []byte {
	return append(AllApplicationsByConsensusAddrKey, addr.Bytes()...)
}

// generates the key for unstaking applications by the unstakingtime
func KeyForUnstakingApps(unstakingTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(unstakingTime)
	return append(UnstakingAppsKey, bz...) // use the unstaking time as part of the key
}

// generates the key for a application in the staking set
func KeyForAppInStakingSet(app Application) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	return getStakedValPowerRankKey(app)
}

func KeyForAppBurn(address sdk.ValAddress) []byte {
	return append(BurnApplicationKey, address...)
}

// get the power ranking key of a application
// NOTE the larger values are of higher value
func getStakedValPowerRankKey(application Application) []byte {
	// get the consensus power
	consensusPower := sdk.TokensToConsensusPower(application.StakedTokens)
	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	// generate the key for this application by deriving it from the main key
	key[0] = StakedAppsKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(application.Address)
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}
