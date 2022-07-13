package dedup

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

//func HeightKey(height int64, prefix string, k []byte) []byte {
//	return []byte(fmt.Sprintf("%d/%s/%s", height, prefix, string(k)))
//}
//
//func FromHeightKey(heightKey string) (height int64, prefix string, k []byte) {
//	var delim = "/"
//	arr := strings.Split(heightKey, delim)
//	// get height
//	height, err := strconv.ParseInt(arr[0], 10, 64)
//	if err != nil {
//		panic("unable to parse height from height key: " + heightKey)
//	}
//	prefix = arr[1]
//	k = []byte(strings.Join(arr[2:], delim))
//	return
//}

func HeightKey(height int64, prefix string, k []byte) []byte {
	h := make([]byte, 4)
	binary.LittleEndian.PutUint32(h, uint32(height))
	return append(append(h, PrefixToByte(prefix)), k...)
}

func FromHeightKey(heightKey string) (height int64, prefix string, k []byte) {
	return int64(binary.LittleEndian.Uint32([]byte(heightKey[:4]))), ByteToPrefix(heightKey[4]), []byte(heightKey[5:])
}

func KeyFromHeightKey(heightKey []byte) (k []byte) {
	_, _, k = FromHeightKey(string(heightKey))
	return
}

func HashKey(key []byte) []byte {
	return sdk.Hash(key)[:16]
}

const orphanPrefix = "o/"

func OrphanPrefixKey(key []byte) []byte {
	return append([]byte(orphanPrefix), key...)
}

func OrphanKey(height int64, prefix string, dataKey []byte) []byte {
	heightKey := HeightKey(height+1, prefix, dataKey)
	return append([]byte(orphanPrefix), heightKey...)
}

func KeyFromOrphanKey(orphanKey []byte) []byte {
	return orphanKey[len(orphanPrefix):]
}

// util

const (
	DefaultCacheKeepHeights = 15
)

func getPreloadStartHeight(latestHeight int64) int64 {
	startHeight := latestHeight - DefaultCacheKeepHeights
	if startHeight < 0 {
		startHeight = 0
	}
	return startHeight
}

type OperationType int

const (
	Set OperationType = iota + 1
	Del
)

const (
	POSPrefix            = "pos"
	POSPrefixByte        = byte(0xa)
	AppPrefix            = "application"
	AppPrefixByte        = byte(0xb)
	PocketCorePrefix     = "pocketcore"
	PocketCorePrefixByte = byte(0xc)
	GovPrefix            = "gov"
	GovPrefixByte        = byte(0xd)
	AuthPrefix           = "auth"
	AuthPrefixByte       = byte(0xe)
	ParamsPrefix         = "params"
	ParamsPrefixByte     = byte(0xf)
	MainPrefix           = "main"
	MainPrefixByte       = byte(0x0)
)

func PrefixToByte(prefix string) (p byte) {
	switch prefix {
	case POSPrefix:
		return POSPrefixByte
	case AppPrefix:
		return AppPrefixByte
	case PocketCorePrefix:
		return PocketCorePrefixByte
	case GovPrefix:
		return GovPrefixByte
	case AuthPrefix:
		return AuthPrefixByte
	case ParamsPrefix:
		return ParamsPrefixByte
	case MainPrefix:
		return MainPrefixByte
	default:
		panic("unknown prefix: " + prefix)
	}
}

func ByteToPrefix(p byte) (prefix string) {
	switch p {
	case POSPrefixByte:
		return POSPrefix
	case AppPrefixByte:
		return AppPrefix
	case PocketCorePrefixByte:
		return PocketCorePrefix
	case GovPrefixByte:
		return GovPrefix
	case AuthPrefixByte:
		return AuthPrefix
	case ParamsPrefixByte:
		return ParamsPrefix
	case MainPrefixByte:
		return MainPrefix
	default:
		panic(fmt.Sprintf("unknown prefix byte %v", p))
	}
}
