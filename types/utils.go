package types

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"time"

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
