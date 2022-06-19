package codec

import (
	"encoding/hex"
	lru "github.com/hashicorp/golang-lru"
	"reflect"
)

// Codec cache speeds up the app by saving previously decoded byte objects in memory.
// This will significantly improve subsequent recent Unmarshal calls for data.
// The implementation uses LRU methodology with a high capacity

var GlobalCodecCache = NewCodecCache(5000000)

type CodecCache struct {
	lru *lru.Cache
}

func NewCodecCache(capacity int) *CodecCache {
	cache, err := lru.New(capacity)
	if err != nil {
		panic("an error occurred initializing the codec cache: " + err.Error())
	}
	return &CodecCache{
		lru: cache,
	}
}

func (cc *CodecCache) GetAndAssign(bz []byte, ptr interface{}) bool {
	val, _ := cc.lru.Get(hex.EncodeToString(bz))
	if val == nil {
		return false
	}
	v := reflect.ValueOf(ptr).Elem()
	oVal := reflect.ValueOf(val)
	if v.String() == oVal.String() {
		// set pointer value as cached object
		v.Set(oVal)
		return true
	}
	// not assignable
	return false
}

func (cc *CodecCache) Add(bz []byte, val interface{}) {
	cc.lru.Add(hex.EncodeToString(bz), val)
}
func (cc *CodecCache) AddPtr(bz []byte, ptr interface{}) {
	v := reflect.ValueOf(ptr)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	cc.Add(bz, v.Interface())
}
