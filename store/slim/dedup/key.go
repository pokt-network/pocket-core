package dedup

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"strconv"
	"strings"
)

func HeightKey(height int64, prefix string, k []byte) []byte {
	return []byte(fmt.Sprintf("%d/%s/%s", height, prefix, string(k)))
}

func FromHeightKey(heightKey string) (height int64, prefix string, k []byte) {
	var delim = "/"
	arr := strings.Split(heightKey, delim)
	// get height
	height, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		panic("unable to parse height from height key: " + heightKey)
	}
	prefix = arr[1]
	k = []byte(strings.Join(arr[2:], delim))
	return
}

func KeyFromHeightKey(heightKey []byte) (k []byte) {
	_, _, k = FromHeightKey(string(heightKey))
	return
}

func HashKey(key []byte) []byte {
	return sdk.Hash(key)
}

const orphanPrefix = "orphan/"

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
