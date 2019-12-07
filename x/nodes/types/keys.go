package types

import (
	"encoding/binary"
	sdk "github.com/pokt-network/posmint/types"
	"time"
)

const (
	ModuleName = "pos"
	StoreKey   = ModuleName // StoreKey is the string store representation
	// TStoreKey is the string transient store representation
	TStoreKey    = "transient_" + ModuleName
	QuerierRoute = ModuleName // QuerierRoute is the querier route for the staking module
	RouterKey    = ModuleName // RouterKey is the msg router key for the staking module
)

//nolint
var ( // Keys for store prefixes
	ProposerKey                     = []byte{0x01} // key for the proposer address used for rewards
	ValidatorSigningInfoKey         = []byte{0x11} // Prefix for signing info used in slashing
	ValidatorMissedBlockBitArrayKey = []byte{0x12} // Prefix for missed block bit array used in slashing
	AddrPubkeyRelationKey           = []byte{0x13} // Prefix for address-pubkey relation used in slashing
	AllValidatorsKey                = []byte{0x21} // prefix for each key to a validator
	AllValidatorsByConsensusAddrKey = []byte{0x22} // prefix for each key to a validator index, by pubkey
	StakedValidatorsKey             = []byte{0x23} // prefix for each key to a staked validator index, sorted by power
	PrevStateValidatorsPowerKey     = []byte{0x31} // prefix for the key to the validators of the prevState state
	PrevStateTotalPowerKey          = []byte{0x32} // prefix for the total power of the prevState state
	UnstakingValidatorsKey          = []byte{0x41} // prefix for unstaking validator
	UnstakedValidatorsKey           = []byte{0x42} // prefix for unstaked validators
	AwardValidatorKey               = []byte{0x51} // prefix for awarding validators
	BurnValidatorKey                = []byte{0x52} // prefix for awarding validators
	WaitingToBeginUnstakingKey      = []byte{0x43}
)

func KeyForValWaitingToBeginUnstaking(addr sdk.ValAddress) []byte {
	return append(WaitingToBeginUnstakingKey, addr.Bytes()...)
}

// generates the key for the validator with address
func KeyForValByAllVals(addr sdk.ValAddress) []byte {
	return append(AllValidatorsKey, addr.Bytes()...)
}

// generates the key for the validator with consensus address
func KeyForValidatorByConsAddr(addr sdk.ConsAddress) []byte {
	return append(AllValidatorsByConsensusAddrKey, addr.Bytes()...)
}

// generates the key for unstaking validators by the unstakingtime
func KeyForUnstakingValidators(unstakingTime time.Time) []byte {
	bz := sdk.FormatTimeBytes(unstakingTime)
	return append(UnstakingValidatorsKey, bz...) // use the unstaking time as part of the key
}

// generates the key for a validator in the staking set
func KeyForValidatorInStakingSet(validator Validator) []byte {
	// NOTE the address doesn't need to be stored because counter bytes must always be different
	return getStakedValPowerRankKey(validator)
}

// generates the key for a validator in the prevState state
func KeyForValidatorPrevStateStateByPower(address sdk.ValAddress) []byte {
	return append(PrevStateValidatorsPowerKey, address...)
}

// generates the award key for a validator in the current state
func KeyForValidatorAward(address sdk.ValAddress) []byte {
	return append(AwardValidatorKey, address...)
}

func KeyForValidatorBurn(address sdk.ValAddress) []byte {
	return append(BurnValidatorKey, address...)
}

// Removes the prefix bytes from a key to expose true address
func AddressFromPrevStateValidatorPowerKey(key []byte) []byte {
	return key[1:] // remove prefix bytes
}

// get the power ranking key of a validator
// NOTE the larger values are of higher value
func getStakedValPowerRankKey(validator Validator) []byte {
	// get the consensus power
	consensusPower := sdk.TokensToConsensusPower(validator.StakedTokens)
	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	// generate the key for this validator by deriving it from the main key
	key[0] = StakedValidatorsKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(validator.Address)
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}

// parse the validators address from power rank key
func ParseValidatorPowerRankKey(key []byte) (operAddr []byte) {
	powerBytesLen := 8
	if len(key) != 1+powerBytesLen+sdk.AddrLen {
		panic("Invalid validator power rank key length")
	}
	operAddr = sdk.CopyBytes(key[powerBytesLen+1:])
	for i, b := range operAddr {
		operAddr[i] = ^b
	}
	return operAddr
}

// generates the key for validator signing information by consensus addr
func GetValidatorSigningInfoKey(v sdk.ConsAddress) []byte {
	return append(ValidatorSigningInfoKey, v.Bytes()...)
}

// extract the address from a validator signing info key
func GetValidatorSigningInfoAddress(key []byte) (v sdk.ConsAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return addr
}

// generates the prefix key for missing val who missed block through consensus addr
func GetValMissedBlockPrefixKey(v sdk.ConsAddress) []byte {
	return append(ValidatorMissedBlockBitArrayKey, v.Bytes()...)
}

// generates the key for missing val who missed block through consensus addr
func GetValMissedBlockKey(v sdk.ConsAddress, i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return append(GetValMissedBlockPrefixKey(v), b...)
}

// generates pubkey relation key used to get the pubkey from the address
func GetAddrPubkeyRelationKey(address []byte) []byte {
	return append(AddrPubkeyRelationKey, address...)
}
