// Package cache provides generic FIFO and LRU caches for the go-textual
// framework. These caches are used throughout the strip, widget, and layout
// packages to memoize expensive render, crop, and style computations.
package cache

import "container/list"

// FIFOCache is a fixed-capacity, first-in-first-out generic cache.
// When the cache is full, the oldest entry is evicted.
type FIFOCache[K comparable, V any] struct {
	cap   int
	items map[K]V
	keys  []K // insertion order
}

// NewFIFO creates a new FIFOCache with the given capacity.
// A capacity of zero disables caching (Get always misses, Set is a no-op).
func NewFIFO[K comparable, V any](capacity int) *FIFOCache[K, V] {
	return &FIFOCache[K, V]{
		cap:   capacity,
		items: make(map[K]V, capacity),
	}
}

// Get returns the value for key and whether it was found.
func (c *FIFOCache[K, V]) Get(key K) (V, bool) {
	v, ok := c.items[key]
	return v, ok
}

// Set stores key→value. If the cache is at capacity the oldest entry is
// evicted. If cap is zero the call is a no-op.
func (c *FIFOCache[K, V]) Set(key K, value V) {
	if c.cap == 0 {
		return
	}
	if _, exists := c.items[key]; !exists {
		if len(c.keys) >= c.cap {
			oldest := c.keys[0]
			c.keys = c.keys[1:]
			delete(c.items, oldest)
		}
		c.keys = append(c.keys, key)
	}
	c.items[key] = value
}

// Len returns the number of entries currently in the cache.
func (c *FIFOCache[K, V]) Len() int { return len(c.items) }

// Clear removes all entries from the cache.
func (c *FIFOCache[K, V]) Clear() {
	c.items = make(map[K]V, c.cap)
	c.keys = c.keys[:0]
}

// -----------------------------------------------------------------------

// lruEntry is stored in the LRU list.
type lruEntry[K comparable, V any] struct {
	key   K
	value V
}

// LRUCache is a fixed-capacity least-recently-used generic cache.
// When the cache is full, the least-recently-used entry is evicted.
type LRUCache[K comparable, V any] struct {
	cap   int
	items map[K]*list.Element
	list  *list.List // front = most recently used
}

// NewLRU creates a new LRUCache with the given capacity.
// A capacity of zero disables caching.
func NewLRU[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		cap:   capacity,
		items: make(map[K]*list.Element, capacity),
		list:  list.New(),
	}
}

// Get returns the value for key and whether it was found.
// A hit moves the entry to the front (most recently used).
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	if el, ok := c.items[key]; ok {
		c.list.MoveToFront(el)
		return el.Value.(*lruEntry[K, V]).value, true
	}
	var zero V
	return zero, false
}

// Set stores key→value. If the cache is at capacity the LRU entry is evicted.
// If cap is zero the call is a no-op.
func (c *LRUCache[K, V]) Set(key K, value V) {
	if c.cap == 0 {
		return
	}
	if el, ok := c.items[key]; ok {
		el.Value.(*lruEntry[K, V]).value = value
		c.list.MoveToFront(el)
		return
	}
	if c.list.Len() >= c.cap {
		// Evict least recently used (back of list).
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			delete(c.items, oldest.Value.(*lruEntry[K, V]).key)
		}
	}
	entry := &lruEntry[K, V]{key: key, value: value}
	el := c.list.PushFront(entry)
	c.items[key] = el
}

// Len returns the number of entries currently in the cache.
func (c *LRUCache[K, V]) Len() int { return c.list.Len() }

// Clear removes all entries from the cache.
func (c *LRUCache[K, V]) Clear() {
	c.items = make(map[K]*list.Element, c.cap)
	c.list.Init()
}
