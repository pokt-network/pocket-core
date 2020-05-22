package types

import (
	"github.com/hashicorp/golang-lru/simplelru"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache struct {
	lru simplelru.LRUCache
}

// New creates an LRU of the given size.
func New(size int) (*Cache, error) {
	lru, err := simplelru.NewLRU(size, nil)
	if err != nil {
		return nil, err
	}
	c := &Cache{
		lru: lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (c *Cache) Purge() {
	c.lru.Purge()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) Add(key string, value CacheObject) (evicted bool) {
	evicted = c.lru.Add(key, value)
	return evicted
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value CacheObject, ok bool) {
	v, ok := c.lru.Get(key)
	if !ok {
		return
	}
	value, ok = v.(CacheObject)
	return
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *Cache) Contains(key string) bool {
	containKey := c.lru.Contains(key)
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache) Peek(key string) (value CacheObject, ok bool) {
	v, ok := c.lru.Peek(key)
	if !ok {
		return
	}
	value, ok = v.(CacheObject)
	return value, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) ContainsOrAdd(key, value CacheObject) (ok, evicted bool) {
	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) PeekOrAdd(key, value CacheObject) (previous interface{}, ok, evicted bool) {
	previous, ok = c.lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lru.Add(key, value)
	return nil, false, evicted
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) (present bool) {
	present = c.lru.Remove(key)
	return
}

// Resize changes the cache size.
func (c *Cache) Resize(size int) (evicted int) {
	evicted = c.lru.Resize(size)
	return evicted
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() (key string, value CacheObject, ok bool) {
	k, v, ok := c.lru.RemoveOldest()
	if !ok {
		return
	}
	value, ok = v.(CacheObject)
	if !ok {
		return
	}
	key, ok = k.(string)
	return
}

// GetOldest returns the oldest entry
func (c *Cache) GetOldest() (key string, value CacheObject, ok bool) {
	k, v, ok := c.lru.GetOldest()
	if !ok {
		return
	}
	value, ok = v.(CacheObject)
	if !ok {
		return
	}
	key, ok = k.(string)
	return
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
