package types

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb/opt"
	dbm "github.com/tendermint/tm-db"
)

var (
	// This is set at compile time. Could be cleveldb, defaults is goleveldb.
	DBBackend         = ""
	VbCCache          *Cache
	ShowTimeTrackData = false
)

func init() {
	VbCCache = NewCache(1200)
	ShowTimeTrackData = false
}

func GetCacheKey(height int, value string) (key string) {
	sh := strconv.Itoa(height)
	key = fmt.Sprintf("%s-%s", sh, value)

	return key
}

func IsBetween(target, minInclusive, maxInclusive int64) bool {
	return minInclusive <= target && target <= maxInclusive
}

// SortedJSON takes any JSON and returns it sorted by keys. Also, all white-spaces
// are removed.
// This method can be used to canonicalize JSON to be returned by GetSignBytes,
// e.g. for the ledger integration.
// If the passed JSON isn't valid it will return an error.
func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}

// Uint64ToBigEndian - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint64ToBigEndian(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// Slight modification of the RFC3339Nano but it right pads all zeros and drops the time zone info
const SortableTimeFormat = "2006-01-02T15:04:05.000000000"

// Formats a time.Time into a []byte that can be sorted
func FormatTimeBytes(t time.Time) []byte {
	return []byte(t.UTC().Round(0).Format(SortableTimeFormat))
}

// Parses a []byte encoded using FormatTimeKey back into a time.Time
func ParseTimeBytes(bz []byte) (time.Time, error) {
	str := string(bz)
	t, err := time.Parse(SortableTimeFormat, str)
	if err != nil {
		return t, err
	}
	return t.UTC().Round(0), nil
}

// NewLevelDB instantiate a new LevelDB instance according to DBBackend.
func NewLevelDB(name, dir string, o *opt.Options) (db dbm.DB, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()
	db, err = dbm.NewGoLevelDBWithOpts(name, dir, o)
	if err != nil {
		return nil, err
	}
	return db, err
}

// Raw is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type Raw []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m Raw) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *Raw) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

var _ json.Marshaler = (*Raw)(nil)

func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	if ShowTimeTrackData {
		log.Println(fmt.Sprintf("%s took %s", name, elapsed))
	}
}

// Compares two version strings, which are expected to be dot-delimited
// integers like "1.2.3.4".  The result is similar to strcmp in C, negative
// if the first version string is considered to be earlier, positive if the
// second version string is considered to be earlier, and zero if both version
// strings are the same.  If any of the given version strings is not
// dot-delimited, the function returns an error.
// For more details, see Test_CompareVersionStrings in utils_test.go.
func CompareVersionStrings(verStr1, verStr2 string) (int, error) {
	ver1 := strings.Split(verStr1, ".")
	ver2 := strings.Split(verStr2, ".")
	lenVer1 := len(ver1)
	lenVer2 := len(ver2)

	numChunks := lenVer1
	if lenVer2 < numChunks {
		numChunks = lenVer2
	}

	for i := 0; i < numChunks; i++ {
		verNum1, err := strconv.Atoi(ver1[i])
		if err != nil {
			return 0, err
		}

		verNum2, err := strconv.Atoi(ver2[i])
		if err != nil {
			return 0, err
		}

		if verNum1 < verNum2 {
			return -1, nil
		}

		if verNum1 > verNum2 {
			return 1, nil
		}
	}

	if lenVer2 > numChunks {
		return -1, nil
	}

	if lenVer1 > numChunks {
		return 1, nil
	}

	return 0, nil
}

// True if two maps are equivalent.
// Nil is considered to be the same as an empty map.
func CompareStringMaps[T comparable](a, b map[string]T) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if v != b[k] {
			return false
		}
	}

	return true
}
