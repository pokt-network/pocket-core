package types

import (
	lru2 "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/golang-lru/simplelru"
)

// Cache is a thread-safe fixed cap LRU cache.
type Cache struct {
	lru simplelru.LRUCache
	cap int
}

// New creates an LRU of the given cap.
func NewCache(size int) *Cache {
	lru, err := lru2.New(size)
	if err != nil {
		panic(err)
	}
	c := &Cache{
		lru: lru,
		cap: size,
	}
	return c
}

// Purge is used to completely clear the cache.
func (c *Cache) Purge() {
	c.lru.Purge()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) Add(key string, value interface{}) (evicted bool) {
	evicted = c.lru.Add(key, value)
	return evicted
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) AddWithCtx(ctx Ctx, key string, value interface{}) (evicted bool) {
	if ctx.IsPrevCtx() {
		return
	}
	evicted = c.lru.Add(key, value)
	return evicted
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	return c.lru.Get(key)
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) GetWithCtx(ctx Ctx, key string) (value interface{}, ok bool) {
	if ctx.IsPrevCtx() {
		return
	}
	return c.Get(key)
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *Cache) Contains(key string) bool {
	return c.lru.Contains(key)
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache) Peek(key string) (value interface{}, ok bool) {
	return c.lru.Peek(key)
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) PeekOrAdd(key, value interface{}) (previous interface{}, ok, evicted bool) {
	previous, ok = c.lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lru.Add(key, value)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
func (c *Cache) RemoveWithCtx(ctx Ctx, key string) (present bool) {
	if ctx.IsPrevCtx() {
		return
	}
	return c.Remove(key)
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) (present bool) {
	present = c.lru.Remove(key)
	return
}

// Resize changes the cache capacity.
func (c *Cache) Resize(size int) (evicted int) {
	evicted = c.lru.Resize(size)
	return evicted
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() (key string, value interface{}, ok bool) {
	k, v, ok := c.lru.RemoveOldest()
	if !ok {
		return
	}
	key, ok = k.(string)
	return key, v, ok
}

// GetOldest returns the oldest entry
func (c *Cache) GetOldest() (key string, value interface{}, ok bool) {
	k, v, ok := c.lru.GetOldest()
	if !ok {
		return
	}
	key, ok = k.(string)
	return key, v, ok
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *Cache) Keys() []interface{} {
	keys := c.lru.Keys()
	return keys
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	length := c.lru.Len()
	return length
}

func (c *Cache) Cap() int {
	return c.cap
}
