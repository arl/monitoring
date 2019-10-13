package main

import (
	"container/list"
	"sync"
)

// A Cache stores already data for faster retrieval.
type Cache interface {
	// Add adds a (key, value) pair to the cache.
	Add(key string, value interface{})
	// Get retrieves the value corresponding to key.
	Get(key string) (value interface{}, ok bool)
}

// A LRUCache implements implements a cache with LRU policy.
//
// An LRUCache has a fixed size, when full the Least Recently Used element
// is removed.
type LRUCache struct {
	mu sync.Mutex               // protects concurrent access on m and l
	m  map[string]*list.Element // maps cached keys to list elements in l
	l  *list.List               // list of cached values

	maxcap int // maximum cache capacity
}

type lruNode struct {
	key   string
	value interface{}
}

// NewLRUCache creates a new LRUCache of maximum capacity maxcap.
func NewLRUCache(maxcap int) *LRUCache {
	if maxcap < 0 {
		panic("LRUCache maximum capacity must be positive!")
	}
	return &LRUCache{
		maxcap: maxcap,
		m:      make(map[string]*list.Element),
		l:      list.New(),
	}
}

// Add adds a (key, value) pair to the cache.
func (c *LRUCache) Add(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.m[k]; ok {
		// We already have that element, but let's move it up front
		c.l.MoveToFront(elem)
		// Update value
		elem.Value.(*lruNode).value = v
	} else {
		// Add a new element at the front of the list
		elem := c.l.PushFront(&lruNode{key: k, value: v})
		// Saves that element in the map for fast 0[1] lookup
		c.m[k] = elem
	}

	if c.l.Len() > c.maxcap {
		// We got too big, remove the least recently used element (back of the list)
		elem := c.l.Back()
		c.l.Remove(elem)
		// Removes it from the hashmap as well
		delete(c.m, elem.Value.(*lruNode).key)
	}
}

// Get retrieves the value corresponding to key.
func (c *LRUCache) Get(k string) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.m[k]
	if !ok {
		// Cache miss
		return nil, false
	}
	// Cache hit: move key up front
	c.l.MoveToFront(elem)
	return elem.Value.(*lruNode).value, true
}
